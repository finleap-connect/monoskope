package reactors

import (
	"context"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
)

type clusterBootstrapReactor struct {
	signer jwt.JWTSigner
}

// NewClusterBootstrapReactor creates a new Reactor.
func NewAggregateHandler(signer jwt.JWTSigner) es.Reactor {
	return &clusterBootstrapReactor{
		signer: signer,
	}
}

// HandleEvent handles a given event returns 0..* Events in reaction or an error
func (r *clusterBootstrapReactor) HandleEvent(ctx context.Context, event es.Event) ([]es.Event, error) {
	eventsToEmit := make([]es.Event, 0)

	switch event.EventType() {
	case events.ClusterCreated:
		data := &eventdata.ClusterCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}

		rawJWT, err := r.signer.GenerateSignedToken(jwt.NewClusterBootstrapToken(data.Name))
		if err != nil {
			return nil, err
		}

		eventData := &eventdata.ClusterBootstrapTokenCreated{
			JWT: rawJWT,
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
