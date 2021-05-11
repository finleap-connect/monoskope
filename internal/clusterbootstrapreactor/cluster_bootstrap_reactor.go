package clusterboostrap

import (
	"context"

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
	switch event.EventType() {
	case events.ClusterCreated:
		//TODO: Create new shiny bootstrap token
	}
	panic("not implemented")
}
