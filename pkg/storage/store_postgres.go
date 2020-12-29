package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/google/uuid"
)

// EventStore implements an EventStore for PostgreSQL.
type EventStore struct {
	db *pg.DB
}

// EventRecord is the model for entries in the events table in the database.
type EventRecord struct {
	tableName struct{} `sql:"events"`

	EventID          uuid.UUID       `sql:"event_id,type:uuid,pk"`
	EventType        EventType       `sql:"event_type,type:varchar(250)"`
	AggregateID      uuid.UUID       `sql:"aggregate_id,type:uuid,unique:aggregate"`
	AggregateType    AggregateType   `sql:"aggregate_type,type:varchar(250),unique:aggregate"`
	AggregateVersion uint64          `sql:"aggregate_version,unique:aggregate"`
	Timestamp        time.Time       `sql:"timestamp"`
	Context          json.RawMessage `sql:"context,type:jsonb"`
	RawData          json.RawMessage `sql:"data,type:jsonb"`
}

var tables []interface{}

func init() {
	// Just to silence linter
	eventsTbl := &EventRecord{}
	_ = eventsTbl.tableName

	tables = []interface{}{
		(*EventRecord)(nil),
	}
}

// createTables creates the event table in the database.
func (s *EventStore) createTables(opts *orm.CreateTableOptions) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		for _, table := range tables {
			if err := s.db.CreateTable(table, opts); err != nil {
				return err
			}
		}
		return nil
	})
}

// newEventRecord returns a new EventRecord for an event.
func (s *EventStore) newEventRecord(ctx context.Context, event Event) (*EventRecord, error) {
	// Marshal event data if there is any.
	eventData, err := json.Marshal(event.Data())
	if err != nil {
		return nil, EventStoreError{
			BaseErr: err,
			Err:     ErrCouldNotMarshalEvent,
		}
	}

	// Marshal event context if there is any.
	context, err := json.Marshal(ctx)
	if err != nil {
		return nil, EventStoreError{
			BaseErr: err,
			Err:     ErrCouldNotMarshalEventContext,
		}
	}

	return &EventRecord{
		EventID:          uuid.New(),
		AggregateID:      event.AggregateID(),
		AggregateType:    event.AggregateType(),
		EventType:        event.EventType(),
		RawData:          eventData,
		Timestamp:        event.Timestamp(),
		AggregateVersion: event.AggregateVersion(),
		Context:          context,
	}, nil
}

// NewPostgresEventStore creates a new EventStore.
func NewPostgresEventStore(db *pg.DB) (Store, error) {
	s := &EventStore{
		db: db,
	}
	err := s.createTables(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Save implements the Save method of the EventStore interface.
func (s *EventStore) Save(ctx context.Context, events []Event) error {
	if len(events) == 0 {
		return EventStoreError{
			Err: ErrNoEventsToAppend,
		}
	}

	// Validate incoming events and create all event records.
	eventRecords := make([]EventRecord, len(events))
	aggregateID := events[0].AggregateID()
	aggregateType := events[0].AggregateType()
	nextVersion := events[0].AggregateVersion()
	for i, event := range events {
		// Only accept events belonging to the same aggregate.
		if event.AggregateID() != aggregateID || event.AggregateType() != aggregateType {
			return EventStoreError{
				Err: ErrInvalidAggregateType,
			}
		}

		// Only accept events that apply to the correct aggregate version.
		if event.AggregateVersion() != nextVersion {
			return EventStoreError{
				Err: ErrIncorrectAggregateVersion,
			}
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
	err := s.db.WithContext(ctx).Insert(&eventRecords)
	if pgErr, ok := err.(pg.Error); ok {
		if pgErr.IntegrityViolation() {
			return EventStoreError{
				BaseErr: err,
				Err:     ErrAggregateVersionAlreadyExists,
			}
		}
	}
	if err != nil {
		return EventStoreError{
			BaseErr: err,
			Err:     ErrCouldNotSaveEvents,
		}
	}

	return nil
}

// Load implements the Load method of the EventStore interface.
func (s *EventStore) Load(ctx context.Context, storeQuery *StoreQuery) ([]Event, error) {
	var events []Event

	// Basic query to query all events
	dbQuery := s.db.
		WithContext(ctx).
		Model((*EventRecord)(nil)).
		Order("timestamp ASC").
		Order("aggregate_version ASC")

	// Translate the abstrace query to a postgres query
	mapStoreQuery(storeQuery, dbQuery)

	err := dbQuery.ForEach(func(e *EventRecord) (err error) {
		events = append(events, pgEvent{
			EventRecord: *e,
		})
		return nil
	})
	if err != nil {
		return nil, EventStoreError{
			BaseErr: err,
			Err:     err,
		}
	}

	return events, nil
}

// mapStoreQuery maps the generic query struct to a postgress orm query
func mapStoreQuery(storeQuery *StoreQuery, dbQuery *orm.Query) {
	if storeQuery == nil {
		return
	}

	if storeQuery.AggregateId != nil {
		_ = dbQuery.Where("aggregate_id = ?", storeQuery.AggregateId)
	} else if storeQuery.AggregateType != nil {
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
func (s *EventStore) clear(ctx context.Context) error {
	return s.db.
		WithContext(ctx).
		RunInTransaction(func(tx *pg.Tx) (err error) {
			_, err = tx.Model((*EventRecord)(nil)).Where("1=1").Delete()
			return err
		})
}

// pgEvent is the private implementation of the Event interface for a postgres event store.
type pgEvent struct {
	EventRecord
}

// EventType implements the EventType method of the Event interface.
func (e pgEvent) EventType() EventType {
	return e.EventRecord.EventType
}

// Data implements the Data method of the Event interface.
func (e pgEvent) Data() EventData {
	return EventData(e.RawData)
}

// Timestamp implements the Timestamp method of the Event interface.
func (e pgEvent) Timestamp() time.Time {
	return e.EventRecord.Timestamp
}

// AggregateType implements the AggregateType method of the Event interface.
func (e pgEvent) AggregateType() AggregateType {
	return e.EventRecord.AggregateType
}

// AggrgateID implements the AggrgateID method of the Event interface.
func (e pgEvent) AggregateID() uuid.UUID {
	return e.EventRecord.AggregateID
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e pgEvent) AggregateVersion() uint64 {
	return e.EventRecord.AggregateVersion
}

// String implements the String method of the Event interface.
func (e pgEvent) String() string {
	return fmt.Sprintf("%s@%d", e.EventRecord.EventType, e.EventRecord.AggregateVersion)
}
