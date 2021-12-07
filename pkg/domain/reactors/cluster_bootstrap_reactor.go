// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reactors

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/certificatemanagement"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
)

const (
	DOMAIN                     = "@monoskope.local"
	RECONCILIATION_MAX_BACKOFF = 1 * time.Minute
)

type clusterBootstrapReactor struct {
	log         logger.Logger
	issuerURL   string
	signer      jwt.JWTSigner
	certManager certificatemanagement.CertificateManager
}

// NewClusterBootstrapReactor creates a new Reactor.
func NewClusterBootstrapReactor(issuerURL string, signer jwt.JWTSigner, certManager certificatemanagement.CertificateManager) es.Reactor {
	return &clusterBootstrapReactor{
		log:         logger.WithName("clusterBootstrapReactor"),
		issuerURL:   issuerURL,
		signer:      signer,
		certManager: certManager,
	}
}

// HandleEvent handles a given event returns 0..* Events in reaction or an error
func (r *clusterBootstrapReactor) HandleEvent(ctx context.Context, event es.Event, eventsChannel chan<- es.Event) error {
	ctx, err := users.CreateUserContext(ctx, users.ReactorUser)
	if err != nil {
		return err
	}

	switch event.EventType() {
	case events.ClusterCreated:
		data := &eventdata.ClusterCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		return r.handleClusterCreated(ctx, data.Name, event, eventsChannel)
	case events.ClusterCreatedV2:
		data := &eventdata.ClusterCreatedV2{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		return r.handleClusterCreated(ctx, data.Name, event, eventsChannel)
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

func (r *clusterBootstrapReactor) handleClusterCreated(ctx context.Context, name string, event es.Event, eventsChannel chan<- es.Event) error {
	var email = name + DOMAIN
	r.log.Info("Generating bootstrap token...", "AggregateID", event.AggregateID(), "Name", name)
	rawJWT, err := r.signer.GenerateSignedToken(jwt.NewClusterBootstrapToken(&jwt.StandardClaims{
		Name:  name,
		Email: email,
	}, r.issuerURL, uuid.New().String()))
	if err != nil {
		r.log.Error(err, "Generating bootstrap token failed.", "AggregateID", event.AggregateID(), "Name", name)
		return err
	}
	r.log.Info("Generating bootstrap token succeeded.", "AggregateID", event.AggregateID(), "Name", name)

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
	r.log.Info("Creating user and rolebinding.", "AggregateID", userId, "Name", name, "Email", email)
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
	r.log.Info("Creating user and rolebinding succeeded.", "AggregateID", userId, "Name", name, "Email", email)

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
				Certificate: &common.CertificateChain{
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
