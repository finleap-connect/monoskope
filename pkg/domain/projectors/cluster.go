// Copyright 2022 Monoskope Authors
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

package projectors

import (
	"bytes"
	"context"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
)

type clusterProjector struct {
	*domainProjector
}

func NewClusterProjector() es.Projector[*projections.Cluster] {
	return &clusterProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (u *clusterProjector) NewProjection(id uuid.UUID) *projections.Cluster {
	return projections.NewClusterProjection(id)
}

// Project updates the state of the projection according to the given event.
func (c *clusterProjector) Project(ctx context.Context, event es.Event, cluster *projections.Cluster) (*projections.Cluster, error) {
	// Apply the changes for the event.
	switch event.EventType() {
	case events.ClusterCreated:
		data := new(eventdata.ClusterCreated)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}

		cluster.DisplayName = data.GetName()
		cluster.Name = data.GetLabel()
		cluster.ApiServerAddress = data.GetApiServerAddress()
		cluster.CaCertBundle = data.GetCaCertificateBundle()

		if err := c.projectCreated(event, cluster.DomainProjection); err != nil {
			return nil, err
		}
	case events.ClusterCreatedV2:
		data := new(eventdata.ClusterCreatedV2)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}

		cluster.DisplayName = data.GetDisplayName()
		cluster.Name = data.GetName()
		cluster.ApiServerAddress = data.GetApiServerAddress()
		cluster.CaCertBundle = data.GetCaCertificateBundle()

		if err := c.projectCreated(event, cluster.DomainProjection); err != nil {
			return nil, err
		}
	case events.ClusterUpdated:
		data := new(eventdata.ClusterUpdated)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}
		if len(data.GetDisplayName()) > 0 && cluster.DisplayName != data.GetDisplayName() {
			cluster.DisplayName = data.GetDisplayName()
		}
		if len(data.GetApiServerAddress()) > 0 && cluster.ApiServerAddress != data.GetApiServerAddress() {
			cluster.ApiServerAddress = data.GetApiServerAddress()
		}
		if len(data.GetCaCertificateBundle()) > 0 && !bytes.Equal(cluster.CaCertBundle, data.GetCaCertificateBundle()) {
			cluster.CaCertBundle = data.GetCaCertificateBundle()
		}
	case events.ClusterBootstrapTokenCreated:
		data := new(eventdata.ClusterBootstrapTokenCreated)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}
		cluster.BootstrapToken = data.GetJwt()
	case events.ClusterDeleted:
		if err := c.projectDeleted(event, cluster.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	if err := c.projectModified(event, cluster.DomainProjection); err != nil {
		return nil, err
	}
	cluster.IncrementVersion()

	return cluster, nil
}
