package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// postgresEventStore implements an EventStore for PostgreSQL.
type postgresEventStore struct {
	log         logger.Logger
	db          *pg.DB
	conf        *postgresStoreConfig
	isConnected bool
	ctx         context.Context
	cancel      context.CancelFunc
}

// eventRecord is the model for entries in the events table in the database.
type eventRecord struct {
	tableName struct{} `sql:"events"`

	EventID          uuid.UUID         `pg:"event_id,type:uuid,pk"`
	EventType        evs.EventType     `pg:"event_type,type:varchar(250)"`
	AggregateID      uuid.UUID         `pg:"aggregate_id,type:uuid,unique:aggregate"`
	AggregateType    evs.AggregateType `pg:"aggregate_type,type:varchar(250),unique:aggregate"`
	AggregateVersion uint64            `pg:"aggregate_version,unique:aggregate"`
	Timestamp        time.Time         `pg:""`
	Metadata         map[string][]byte `pg:"metadata,type:jsonb"`
	RawData          json.RawMessage   `pg:"data,type:jsonb"`
}

var models []interface{}

func init() {
	// Just to silence linter
	eventsTbl := &eventRecord{}
	_ = eventsTbl.tableName

	models = []interface{}{
		(*eventRecord)(nil),
	}
}

// createTables creates the event table in the database.
func (s *postgresEventStore) createTables(ctx context.Context, db *pg.DB) error {
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// newEventRecord returns a new EventRecord for an event.
func (s *postgresEventStore) newEventRecord(ctx context.Context, event evs.Event) (*eventRecord, error) {
	return &eventRecord{
		EventID:          uuid.New(),
		AggregateID:      event.AggregateID(),
		AggregateType:    event.AggregateType(),
		EventType:        event.EventType(),
		RawData:          json.RawMessage(event.Data()),
		Timestamp:        event.Timestamp(),
		AggregateVersion: event.AggregateVersion(),
		Metadata:         event.Metadata(),
	}, nil
}

// NewPostgresEventStore creates a new EventStore.
func NewPostgresEventStore(config *postgresStoreConfig) (evs.Store, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := &postgresEventStore{
		log:    logger.WithName("postgres-store"),
		conf:   config,
		ctx:    ctx,
		cancel: cancel,
	}
	return s, nil
}

// Connect starts automatic reconnect with postgres db
func (s *postgresEventStore) Connect(ctx context.Context) error {
	go s.handleReconnect()
	for !s.isConnected {
		select {
		case <-s.ctx.Done():
			s.log.Info("Connection aborted because of shutdown.")
			return errors.ErrCouldNotConnect
		case <-ctx.Done():
			s.log.Info("Connection aborted because context deadline exceeded.")
			return errors.ErrCouldNotConnect
		case <-time.After(300 * time.Millisecond):
		}
	}
	return nil
}

// Save implements the Save method of the EventStore interface.
func (s *postgresEventStore) Save(ctx context.Context, events []evs.Event) error {
	if len(events) == 0 {
		return errors.ErrNoEventsToAppend
	}

	// Validate incoming events and create all event records.
	eventRecords := make([]eventRecord, len(events))
	aggregateID := events[0].AggregateID()
	aggregateType := events[0].AggregateType()
	nextVersion := events[0].AggregateVersion()
	for i, event := range events {
		// Only accept events belonging to the same aggregate.
		if event.AggregateID() != aggregateID || event.AggregateType() != aggregateType {
			return errors.ErrInvalidAggregateType
		}

		// Only accept events that apply to the correct aggregate version.
		if event.AggregateVersion() != nextVersion {
			return errors.ErrIncorrectAggregateVersion
		}

		// Create the event record for the DB.
		e, err := s.newEventRecord(ctx, event)
		if err != nil {
			return err
		}
		eventRecords[i] = *e

		// Increment to checking order of following events.
		nextVersion++
	}

	// Append events to the store.
	err := retryWithExponentialBackoff(5, 500*time.Millisecond, func() (e error) {
		if !s.isConnected {
			return errors.ErrConnectionClosed
		}
		_, e = s.db.Model(&eventRecords).Insert()
		return e
	}, func(e error) bool {
		if pgErr, ok := e.(pg.Error); ok {
			if pgErr.Field(byte('C')) == "40001" { // serialization_failure, see https://www.postgresql.org/docs/10/errcodes-appendix.html, see https://www.cockroachlabs.com/docs/stable/common-errors.html#result-is-ambiguous
				return true
			}
		}
		return false
	})
	if pgErr, ok := err.(pg.Error); ok {
		if pgErr.IntegrityViolation() {
			s.log.Info(errors.ErrAggregateVersionAlreadyExists.Error(), "error", pgErr)
			return errors.ErrAggregateVersionAlreadyExists
		}
	}
	if err != nil {
		s.log.Error(err, errors.ErrCouldNotSaveEvents.Error())
		return errors.ErrCouldNotSaveEvents
	}

	s.log.Info("Saved event(s) successullfy", "eventCount", len(eventRecords))
	return nil
}

// retryWithExponentialBackoff retries a given function on error if either the recoverable function returns true or still attempts left
func retryWithExponentialBackoff(attempts int, initialBackoff time.Duration, f func() error, recoverable func(error) bool) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if !recoverable(err) {
			return err
		}
		if i >= (attempts - 1) {
			break
		}
		milliseconds := math.Pow(float64(initialBackoff.Milliseconds()), float64(2*(i+1)))
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// Load implements the Load method of the EventStore interface.
func (s *postgresEventStore) Load(ctx context.Context, storeQuery *evs.StoreQuery) ([]evs.Event, error) {
	if !s.isConnected {
		return nil, errors.ErrConnectionClosed
	}

	var events []evs.Event

	// Basic query to query all events
	dbQuery := s.db.
		WithContext(ctx).
		Model((*eventRecord)(nil)).
		Order("timestamp ASC")

	// Translate the abstrace query to a postgres query
	mapStoreQuery(storeQuery, dbQuery)

	err := dbQuery.ForEach(func(e *eventRecord) (err error) {
		events = append(events, pgEvent{
			eventRecord: *e,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *postgresEventStore) Close() error {
	s.log.Info("Shutting down...")

	s.cancel()
	err := s.db.Close()
	if err != nil {
		return err
	}

	s.log.Info("Shutdown complete.")

	return nil
}

// mapStoreQuery maps the generic query struct to a postgress orm query
func mapStoreQuery(storeQuery *evs.StoreQuery, dbQuery *orm.Query) {
	if storeQuery == nil {
		return
	}

	if storeQuery.AggregateId != nil {
		_ = dbQuery.Where("aggregate_id = ?", storeQuery.AggregateId)
	}
	if storeQuery.AggregateType != nil {
		_ = dbQuery.Where("aggregate_type = ?", storeQuery.AggregateType)
	}

	if storeQuery.MinVersion != nil {
		_ = dbQuery.Where("aggregate_version >= ?", storeQuery.MinVersion)
	}
	if storeQuery.MaxVersion != nil {
		_ = dbQuery.Where("aggregate_version <= ?", storeQuery.MaxVersion)
	}

	if storeQuery.MinTimestamp != nil {
		_ = dbQuery.Where("timestamp >= ?", storeQuery.MinTimestamp)
	}
	if storeQuery.MaxTimestamp != nil {
		_ = dbQuery.Where("timestamp <= ?", storeQuery.MaxTimestamp)
	}
}

// Clear clears the event storage. This is only for testing purposes.
func (s *postgresEventStore) clear(ctx context.Context) error {
	return s.db.
		RunInTransaction(ctx, func(tx *pg.Tx) (err error) {
			_, err = tx.Model((*eventRecord)(nil)).Where("1=1").Delete()
			return err
		})
}

// init will initialize db
func (s *postgresEventStore) init(db *pg.DB) error {
	err := s.createTables(s.ctx, db)
	s.isConnected = err == nil
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresEventStore) handleReInit(db *pg.DB) error {
	for {
		err := s.init(db)
		if err != nil {
			s.log.Info("Failed to initialize channel. Retrying...", "error", err.Error())

			select {
			case <-s.ctx.Done():
				s.log.Info("Aborting init. Shutting down...")
				return errors.ErrCouldNotConnect
			case <-time.After(s.conf.ReInitDelay):
				continue
			}
		}
		return nil
	}
}

// connect will create a new db connection
func (s *postgresEventStore) connect() (*pg.DB, error) {
	s.log.Info("Attempting to connect...")

	db := pg.Connect(s.conf.pgOptions)
	if err := db.Ping(s.ctx); err != nil {
		return nil, err
	}

	s.db = db
	s.log.Info("Connection established.")

	return db, nil
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (s *postgresEventStore) handleReconnect() {
	for {
		s.isConnected = false
		db, err := s.connect()
		if err != nil {
			s.log.Info("Failed to connect. Retrying...", "error", err.Error())

			select {
			case <-s.ctx.Done():
				s.log.Info("Automatic reconnect stopped. Shutdown.")
				return
			case <-time.After(s.conf.ReconnectDelay):
			}
			continue
		}

		err = s.handleReInit(db)
		if err != nil {
			s.log.Info("Failed to init. Retrying...", "error", err.Error())
		} else {
			for {
				if err := db.Ping(s.ctx); err != nil {
					s.isConnected = false
					s.log.Info("Connection closed. Reconnecting...")
					break
				}
				time.Sleep(s.conf.ReconnectDelay)
			}
		}
	}
}

// pgEvent is the private implementation of the Event interface for a postgres event store.
type pgEvent struct {
	eventRecord
}

// EventType implements the EventType method of the Event interface.
func (e pgEvent) EventType() evs.EventType {
	return e.eventRecord.EventType
}

// Data implements the Data method of the Event interface.
func (e pgEvent) Data() evs.EventData {
	return evs.EventData(e.RawData)
}

// Timestamp implements the Timestamp method of the Event interface.
func (e pgEvent) Timestamp() time.Time {
	return e.eventRecord.Timestamp
}

// AggregateType implements the AggregateType method of the Event interface.
func (e pgEvent) AggregateType() evs.AggregateType {
	return e.eventRecord.AggregateType
}

// AggrgateID implements the AggrgateID method of the Event interface.
func (e pgEvent) AggregateID() uuid.UUID {
	return e.eventRecord.AggregateID
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e pgEvent) AggregateVersion() uint64 {
	return e.eventRecord.AggregateVersion
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e pgEvent) Metadata() map[string][]byte {
	return e.eventRecord.Metadata
}

// String implements the String method of the Event interface.
func (e pgEvent) String() string {
	return fmt.Sprintf("%s@%d", e.eventRecord.EventType, e.eventRecord.AggregateVersion)
}
