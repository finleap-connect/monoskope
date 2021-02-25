package usecases

import (
	"io"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/metrics"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StoreEventsUseCase struct {
	UseCaseBase

	store   es.Store
	bus     es.EventBusPublisher
	stream  esApi.EventStore_StoreServer
	metrics *metrics.EventStoreMetrics
}

// NewStoreEventsUseCase creates a new usecase which stores all events in the store
// and broadcasts these events via the message bus
func NewStoreEventsUseCase(stream esApi.EventStore_StoreServer, store es.Store, bus es.EventBusPublisher, metrics *metrics.EventStoreMetrics) UseCase {
	useCase := &StoreEventsUseCase{
		UseCaseBase: UseCaseBase{
			log: logger.WithName("store-events-use-case"),
			ctx: stream.Context(),
		},
		store:   store,
		bus:     bus,
		stream:  stream,
		metrics: metrics,
	}
	return useCase
}

func (u *StoreEventsUseCase) Run() error {
	for {
		startTime := time.Now()

		// Read next event
		event, err := u.stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}

		if err != nil { // Some other error
			return errors.TranslateToGrpcError(err)
		}

		// Count transmitted event
		u.metrics.TransmittedTotalCounter.WithLabelValues(event.Type, event.AggregateType).Inc()

		// Convert from proto events to storage events
		ev, err := es.NewEventFromProto(event)
		if err != nil {
			return err
		}

		// Store events in database
		u.log.Info("Saving events in the store...")
		if err := u.store.Save(u.ctx, []es.Event{ev}); err != nil {
			return err
		}

		// Count sucessfully stored event
		u.metrics.StoredTotalCounter.WithLabelValues(event.Type, event.AggregateType).Inc()

		// Send events to message bus
		u.log.Info("Sending events to the message bus...")
		if err := u.bus.PublishEvent(u.ctx, ev); err != nil {
			return err
		}
		u.metrics.StoredHistogram.WithLabelValues(event.Type, event.AggregateType).Observe(time.Since(startTime).Seconds())
	}

	return u.stream.SendAndClose(&emptypb.Empty{})
}
