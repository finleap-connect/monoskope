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
	db      *pg.DB
	encoder Encoder
}

type EventLog struct {
	tableName struct{} `sql:"events"`

	EventID        uuid.UUID              `sql:"event_id,type:uuid,pk"`
	SequenceNumber uint64                 `sql:"sequence_number,type:serial,pk"`
	EventType      EventType              `sql:"event_type,type:varchar(250)"`
	AggregateID    uuid.UUID              `sql:"aggregate_id,type:uuid,unique:aggregate"`
	AggregateType  AggregateType          `sql:"aggregate_type,type:varchar(250),unique:aggregate"`
	Version        uint64                 `sql:"version,unique:aggregate"`
	Timestamp      time.Time              `sql:"timestamp"`
	Context        map[string]interface{} `sql:"context"`
	RawData        json.RawMessage        `sql:"data,type:jsonb"`
	data           EventData              `sql:"-"`
}

var tables = []interface{}{
	(*EventLog)(nil),
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

// NewEventStore creates a new EventStore.
func NewEventStore(db *pg.DB) (*EventStore, error) {
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
