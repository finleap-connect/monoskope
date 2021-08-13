package eventsourcing

import (
	"sync"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type AggregateRegistry interface {
	RegisterAggregate(func() Aggregate)
	CreateAggregate(AggregateType) (Aggregate, error)
}

type aggregateRegistry struct {
	log        logger.Logger
	mutex      sync.RWMutex
	aggregates map[AggregateType]func() Aggregate
}

var DefaultAggregateRegistry AggregateRegistry

func init() {
	DefaultAggregateRegistry = NewAggregateRegistry()
}

// NewAggregateRegistry creates a new aggregate registry
func NewAggregateRegistry() AggregateRegistry {
	return &aggregateRegistry{
		log:        logger.WithName("aggregate-registry"),
		aggregates: make(map[AggregateType]func() Aggregate),
	}
}

// RegisterAggregate registers an aggregate factory for a type. The factory is
// used to create concrete aggregate types.
//
// An example would be:
//     RegisterAggregate(func() Aggregate { return &MyAggregate{} })
func (r *aggregateRegistry) RegisterAggregate(factory func() Aggregate) {
	aggregate := factory()
	if aggregate == nil {
		r.log.Info("factory does not create aggregates")
		panic(errors.ErrFactoryInvalid)
	}

	aggregateType := aggregate.Type()
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
func (r *aggregateRegistry) CreateAggregate(aggregateType AggregateType) (Aggregate, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if factory, ok := r.aggregates[aggregateType]; ok {
		return factory(), nil
	}
	r.log.Info("trying to create a aggregate of non-registered type", "aggregateType", aggregateType)
	return nil, errors.ErrAggregateNotRegistered
}
