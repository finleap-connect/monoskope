package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/google/uuid"
)

// ErrCouldNotClearDB is when the database could not be cleared.
var ErrCouldNotClearDB = errors.New("could not clear database")

// ErrConflictVersion is when a version conflict occurs when saving an aggregate.
var ErrVersionConflict = errors.New("Can not create/update aggregate")

// ErrCouldNotMarshalEvent is when an event could not be marshaled into JSON.
var ErrCouldNotMarshalEvent = errors.New("could not marshal event")

// ErrCouldNotMarshalEventContext is when an event could not be marshaled into JSON.
var ErrCouldNotMarshalEventContext = errors.New("could not marshal event context")

// ErrCouldNotUnmarshalEvent is when an event could not be unmarshaled into a concrete type.
var ErrCouldNotUnmarshalEvent = errors.New("could not unmarshal event")

// ErrCouldNotSaveEvents is when events could not be saved.
var ErrCouldNotSaveEvents = errors.New("could not save events")

// EventStore implements an EventStore for PostgreSQL.
type EventStore struct {
	db      *pg.DB
	encoder Encoder
}

type EventRecord struct {
	tableName struct{} `sql:"events"`

	EventID          uuid.UUID       `sql:"event_id,type:uuid,pk"`
	SequenceNumber   uint64          `sql:"sequence_number,type:serial,pk"`
	EventType        EventType       `sql:"event_type,type:varchar(250)"`
	AggregateID      uuid.UUID       `sql:"aggregate_id,type:uuid,unique:aggregate"`
	AggregateType    AggregateType   `sql:"aggregate_type,type:varchar(250),unique:aggregate"`
	AggregateVersion uint64          `sql:"aggregate_version,unique:aggregate"`
	Timestamp        time.Time       `sql:"timestamp"`
	Context          json.RawMessage `sql:"context,type:jsonb"`
	RawData          json.RawMessage `sql:"data,type:jsonb"`
	data             EventData       `sql:"-"`
}

var tables = []interface{}{
	(*EventRecord)(nil),
}

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
	eventData, err := s.encoder.Marshal(event.Data())
	if err != nil {
		return nil, EventStoreError{
			BaseErr: err,
			Err:     ErrCouldNotMarshalEvent,
		}
	}

	// Marshal event context if there is any.
	context, err := s.encoder.Marshal(ctx)
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
func NewPostgresEventStore(db *pg.DB) (*EventStore, error) {
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

	// Build all event records, with incrementing versions starting from the
	// original aggregate version.
	eventRecords := make([]EventRecord, len(events))
	aggregateID := events[0].AggregateID()
	aggregateType := events[0].AggregateType()
	nextVersion := events[0].AggregateVersion()
	for i, event := range events {
		// Only accept events belonging to the same aggregate.
		if event.AggregateID() != aggregateID && event.AggregateType() == aggregateType {
			return EventStoreError{
				Err: ErrInvalidEvent,
			}
		}

		// Only accept events that apply to the correct aggregate version.
		if event.AggregateVersion() != nextVersion {
			return EventStoreError{
				Err: ErrIncorrectEventVersion,
			}
		}

		// Create the event record for the DB.
		e, err := s.newEventRecord(ctx, event)
		if err != nil {
			return err
		}
		eventRecords[i] = *e
		nextVersion++
	}

	// Append events to the store.
	err := s.db.WithContext(ctx).Insert(&eventRecords)
	if err != nil {
		return EventStoreError{
			BaseErr: err,
			Err:     ErrCouldNotSaveEvents,
		}
	}

	return nil
}

// Load implements the Load method of the EventStore interface.
func (s *EventStore) Load(ctx context.Context, query *StoreQuery) ([]Event, error) {
	panic("not implemented")
}
