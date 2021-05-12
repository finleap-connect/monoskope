package clusterboostrap

import (
	"context"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type clusterBootstrapReactor struct {
}

// NewClusterBootstrapReactor creates a new Reactor.
func NewAggregateHandler() es.Reactor {
	return &clusterBootstrapReactor{}
}

// HandleEvent handles a given event returns 0..* Events in reaction or an error
func (r *clusterBootstrapReactor) HandleEvent(ctx context.Context, event es.Event) ([]es.Event, error) {
	eventsToEmit := make([]es.Event, 0)

	switch event.EventType() {
	case events.ClusterCreated:
		//TODO: Create new shiny bootstrap token
		eventData := &eventdata.ClusterBootstrapTokenCreated{
			JWT: "JWT",
		}

		eventsToEmit = append(eventsToEmit, es.NewEvent(
			ctx,
			events.ClusterBootstrapTokenCreated,
			es.ToEventDataFromProto(eventData),
			time.Now().UTC(),
			event.AggregateType(),
			event.AggregateID(),
			event.AggregateVersion()+1))
	}

	return eventsToEmit, nil
}
