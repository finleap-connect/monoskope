package eventsourcing

import (
	"sync"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type AggregateRegistry interface {
	RegisterAggregate(func(id uuid.UUID) Aggregate)
	CreateAggregate(AggregateType, uuid.UUID) (Aggregate, error)
}

type aggregateRegistry struct {
	log        logger.Logger
	mutex      sync.RWMutex
	aggregates map[AggregateType]func(id uuid.UUID) Aggregate
}

// NewAggregateRegistry creates a new aggregate registry
func NewAggregateRegistry() AggregateRegistry {
	return &aggregateRegistry{
		log:        logger.WithName("aggregate-registry"),
		aggregates: make(map[AggregateType]func(id uuid.UUID) Aggregate),
	}
}

// RegisterAggregate registers an aggregate factory for a type. The factory is
// used to create concrete aggregate types.
//
// An example would be:
//     RegisterAggregate(func() Aggregate { return &MyAggregate{} })
func (r *aggregateRegistry) RegisterAggregate(factory func(id uuid.UUID) Aggregate) {
	cmd := factory(uuid.Nil)
	if cmd == nil {
		r.log.Info("factory does not create aggregates")
		panic(errors.ErrFactoryInvalid)
	}

	aggregateType := cmd.Type()
	if aggregateType.String() == "" {
		r.log.Info("attempt to register empty aggregate type")
		panic(errors.ErrEmptyAggregateType)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.aggregates[aggregateType]; ok {
		r.log.Info("attempt to register aggregate already registered", "aggregateType", aggregateType)
		panic(errors.ErrAggregateTypeAlreadyRegistered)
	}
	r.aggregates[aggregateType] = factory

	r.log.Info("aggregate has been registered.", "aggregateType", aggregateType)
}

// CreateAggregate creates an aggregate of a type with an ID using the factory
// registered with RegisterAggregate.
func (r *aggregateRegistry) CreateAggregate(aggregateType AggregateType, id uuid.UUID) (Aggregate, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if factory, ok := r.aggregates[aggregateType]; ok {
		return factory(id), nil
	}
	r.log.Info("trying to create a aggregate of non-registered type", "aggregateType", aggregateType)
	return nil, errors.ErrAggregateNotRegistered
}
