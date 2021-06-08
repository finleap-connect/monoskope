package reactors

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/certificatemanagement"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type clusterBootstrapReactor struct {
	log         logger.Logger
	signer      jwt.JWTSigner
	certManager certificatemanagement.CertificateManager
}

// NewClusterBootstrapReactor creates a new Reactor.
func NewClusterBootstrapReactor(signer jwt.JWTSigner, certManager certificatemanagement.CertificateManager) es.Reactor {
	return &clusterBootstrapReactor{
		log:         logger.WithName("clusterBootstrapReactor"),
		signer:      signer,
		certManager: certManager,
	}
}

// HandleEvent handles a given event returns 0..* Events in reaction or an error
func (r *clusterBootstrapReactor) HandleEvent(ctx context.Context, event es.Event, eventsChannel chan<- es.Event) error {
	switch event.EventType() {
	case events.ClusterCreated:
		data := &eventdata.ClusterCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}

		rawJWT, err := r.signer.GenerateSignedToken(jwt.NewClusterBootstrapToken(&jwt.StandardClaims{
			Name:  data.Name,
			Email: data.Name + "@monoskope.io",
		}, uuid.New().String(), "cluster-bootstrap-reactor"))
		if err != nil {
			return err
		}

		eventData := &eventdata.ClusterBootstrapTokenCreated{
			JWT: rawJWT,
		}

		eventsChannel <- es.NewEvent(
			ctx,
			events.ClusterBootstrapTokenCreated,
			es.ToEventDataFromProto(eventData),
			time.Now().UTC(),
			event.AggregateType(),
			event.AggregateID(),
			event.AggregateVersion()+1)
	case events.ClusterCertificateRequested:
		data := &eventdata.ClusterCertificateRequested{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}

		if err := r.certManager.RequestCertificate(ctx, event.AggregateID(), data.GetCertificateSigningRequest()); err != nil {
			return err
		}

		eventsChannel <- es.NewEvent(
			ctx,
			events.ClusterOperatorCertificateRequestIssued,
			nil,
			time.Now().UTC(),
			event.AggregateType(),
			event.AggregateID(),
			event.AggregateVersion()+1)
	}

	return nil
}
