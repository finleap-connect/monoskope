package reactors

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/certificatemanagement"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const (
	ISSUER                     = "cluster-bootstrap-reactor"
	DOMAIN                     = "@monoskope.local"
	RECONCILIATION_MAX_BACKOFF = 1 * time.Minute
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
		var name = data.Name
		var email = data.Label + DOMAIN

		r.log.Info("Generating bootstrap token...", "AggregateID", event.AggregateID(), "Name", data.Name, "Label", data.Label)
		rawJWT, err := r.signer.GenerateSignedToken(jwt.NewClusterBootstrapToken(&jwt.StandardClaims{
			Name:  name,
			Email: email,
		}, uuid.New().String(), ISSUER))
		if err != nil {
			r.log.Error(err, "Generating bootstrap token failed.", "AggregateID", event.AggregateID(), "Name", data.Name, "Label", data.Label)
			return err
		}
		r.log.Info("Generating bootstrap token succeeded.", "AggregateID", event.AggregateID(), "Name", data.Name, "Label", data.Label)

		eventsChannel <- es.NewEvent(
			ctx,
			events.ClusterBootstrapTokenCreated,
			es.ToEventDataFromProto(&eventdata.ClusterBootstrapTokenCreated{
				Jwt: rawJWT,
			}),
			time.Now().UTC(),
			event.AggregateType(),
			event.AggregateID(),
			event.AggregateVersion()+1)

		userId := uuid.New()
		r.log.Info("Creating user and rolebinding.", "AggregateID", userId, "Name", data.Name, "Email", email)
		eventsChannel <- es.NewEvent(
			ctx,
			events.UserCreated,
			es.ToEventDataFromProto(&eventdata.UserCreated{
				Name:  name,
				Email: email,
			}),
			time.Now().UTC(),
			aggregates.User,
			userId,
			1)

		eventsChannel <- es.NewEvent(
			ctx,
			events.UserRoleBindingCreated,
			es.ToEventDataFromProto(&eventdata.UserRoleAdded{
				UserId: userId.String(),
				Role:   roles.K8sOperator.String(),
				Scope:  scopes.System.String(),
			}),
			time.Now().UTC(),
			aggregates.UserRoleBinding,
			uuid.New(),
			1)
		r.log.Info("Creating user and rolebinding succeeded.", "AggregateID", userId, "Name", data.Name, "Email", email)
	case events.CertificateRequested:
		data := &eventdata.CertificateRequested{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		if data.ReferencedAggregateType != aggregates.Cluster.String() {
			return errors.New("event CertificateRequested only supported for aggregate Cluster")
		}

		r.log.Info("Generating certificate signing request...", "AggregateID", event.AggregateID())
		if err := r.certManager.RequestCertificate(ctx, event.AggregateID(), data.GetSigningRequest()); err != nil {
			r.log.Error(err, "Generating certificate signing request failed", "AggregateID", event.AggregateID())
			return err
		}
		r.log.Info("Generating certificate signing request succeeded", "AggregateID", event.AggregateID())

		eventsChannel <- es.NewEvent(
			ctx,
			events.CertificateRequestIssued,
			nil,
			time.Now().UTC(),
			event.AggregateType(),
			event.AggregateID(),
			event.AggregateVersion()+1)

		go r.reconcile(ctx, event, eventsChannel)
	}

	return nil
}

func (r *clusterBootstrapReactor) reconcile(ctx context.Context, event es.Event, eventsChannel chan<- es.Event) {
	defer close(eventsChannel)

	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = RECONCILIATION_MAX_BACKOFF

	err := backoff.Retry(func() error {
		r.log.Info("Certificate reconciliation started...", "AggregateID", event.AggregateID())
		ca, cert, err := r.certManager.GetCertificate(ctx, event.AggregateID())
		if err != nil {
			r.log.Info("Certificate reconciliation finished.", "AggregateID", event.AggregateID(), "State", err)
			return err
		}

		r.log.Info("Certificate reconciliation finished.", "AggregateID", event.AggregateID(), "State", "certificate issued successfully")
		eventsChannel <- es.NewEvent(
			ctx,
			events.CertificateIssued,
			es.ToEventDataFromProto(&eventdata.CertificateIssued{
				Certificate: &common.Certificate{
					Ca:          ca,
					Certificate: cert,
				},
			}),
			time.Now().UTC(),
			event.AggregateType(),
			event.AggregateID(),
			event.AggregateVersion()+1)

		return nil
	}, params)

	if err != nil {
		r.log.Error(err, "Certificate reconciliation failed.")
		eventsChannel <- es.NewEvent(
			ctx,
			events.CertificateIssueingFailed,
			nil,
			time.Now().UTC(),
			event.AggregateType(),
			event.AggregateID(),
			event.AggregateVersion()+1)
	}
}
