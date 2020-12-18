package storage

import (
	"encoding/json"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/google/uuid"
)

// EventStore implements an EventStore for PostgreSQL.
type EventStore struct {
	db *pg.DB
}

type Aggregate struct {
	tableName struct{} `sql:"aggregates"`

	AggregateID uuid.UUID `sql:"aggregate_id,type:uuid,pk"`
	Version     uint64    `sql:"version"`
}

type EventLog struct {
	tableName struct{} `sql:"event_store"`

	EventID       uuid.UUID              `sql:"event_id,type:uuid,pk"`
	AggregateID   uuid.UUID              `sql:"aggregate_id,type:uuid"`
	AggregateType AggregateType          `sql:"aggregate_type,type:varchar(250)"`
	EventType     EventType              `sql:"event_type,type:varchar(250)"`
	Timestamp     time.Time              `sql:"timestamp"`
	Version       uint64                 `sql:"version"`
	Context       map[string]interface{} `sql:"context"`
	RawData       json.RawMessage        `sql:"data,type:jsonb"`
	data          EventData              `sql:"-"`
}

var tables = []interface{}{
	(*EventLog)(nil),
	(*Aggregate)(nil),
}

func (s *EventStore) CreateTables(opts *orm.CreateTableOptions) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		for _, table := range tables {
			if err := s.db.CreateTable(table, opts); err != nil {
				return err
			}
		}
		return nil
	})
}

// NewEventStore creates a new EventStore.
func NewEventStore(db *pg.DB) (*EventStore, error) {
	s := &EventStore{
		db: db,
	}
	err := s.CreateTables(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}
