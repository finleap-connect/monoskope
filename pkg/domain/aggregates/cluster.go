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

package aggregates

import (
	"bytes"
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)

// ClusterAggregate is an aggregate for K8s Clusters.
type ClusterAggregate struct {
	*DomainAggregateBase
	aggregateManager es.AggregateStore
	displayName      string
	name             string
	apiServerAddr    string
	caCertBundle     []byte
	bootstrapToken   string
}

// ClusterAggregate creates a new ClusterAggregate
func NewClusterAggregate(aggregateManager es.AggregateStore) es.Aggregate {
	return &ClusterAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Cluster),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *ClusterAggregate) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	if err := a.validate(ctx, cmd); err != nil {
		return nil, err
	}

	switch cmd := cmd.(type) {
	case *commands.CreateClusterCommand:
		ed := es.ToEventDataFromProto(&eventdata.ClusterCreatedV2{
			DisplayName:         cmd.GetDisplayName(),
			Name:                cmd.GetName(),
			ApiServerAddress:    cmd.GetApiServerAddress(),
			CaCertificateBundle: cmd.GetCaCertBundle(),
		})
		_ = a.AppendEvent(ctx, events.ClusterCreatedV2, ed)
	case *commands.UpdateClusterCommand:
		ed := new(eventdata.ClusterUpdated)

		displayName := cmd.GetDisplayName()
		apiServerAddr := cmd.GetApiServerAddress()
		caCertBundle := cmd.GetCaCertBundle()

		if displayName != nil && a.displayName != displayName.Value {
			ed.DisplayName = displayName.Value
		}
		if apiServerAddr != nil && a.apiServerAddr != apiServerAddr.Value {
			ed.ApiServerAddress = apiServerAddr.Value
		}
		if caCertBundle != nil && !bytes.Equal(a.caCertBundle, caCertBundle) {
			ed.CaCertificateBundle = caCertBundle
		}
		_ = a.AppendEvent(ctx, events.ClusterUpdated, es.ToEventDataFromProto(ed))
	case *commands.DeleteClusterCommand:
		_ = a.AppendEvent(ctx, events.ClusterDeleted, nil)
	default:
		return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
	return a.DefaultReply(), nil
}

// validate validates the current state of the aggregate and if a specific command is valid in the current state
func (a *ClusterAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateClusterCommand:
		if a.Exists() {
			return domainErrors.ErrClusterAlreadyExists
		}

		// Get all aggregates of same type
		aggregates, err := a.aggregateManager.All(ctx, a.Type())
		if err != nil {
			return err
		}

		if containsCluster(aggregates, cmd.GetName()) {
			return domainErrors.ErrClusterAlreadyExists
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

func containsCluster(values []es.Aggregate, name string) bool {
	for _, value := range values {
		d, ok := value.(*ClusterAggregate)
		if ok {
			if !d.Deleted() && d.name == name {
				return true
			}
		}
	}
	return false
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *ClusterAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.ClusterCreated:
		clusterCreatedV1 := new(eventdata.ClusterCreated)
		err := event.Data().ToProto(clusterCreatedV1)
		if err != nil {
			return err
		}
		a.displayName = clusterCreatedV1.GetName()
		a.name = clusterCreatedV1.GetLabel()
		a.apiServerAddr = clusterCreatedV1.GetApiServerAddress()
		a.caCertBundle = clusterCreatedV1.GetCaCertificateBundle()
	case events.ClusterCreatedV2:
		clusterCreatedV2 := new(eventdata.ClusterCreatedV2)
		err := event.Data().ToProto(clusterCreatedV2)
		if err != nil {
			return err
		}
		a.displayName = clusterCreatedV2.GetDisplayName()
		a.name = clusterCreatedV2.GetName()
		a.apiServerAddr = clusterCreatedV2.GetApiServerAddress()
		a.caCertBundle = clusterCreatedV2.GetCaCertificateBundle()
	case events.ClusterBootstrapTokenCreated:
		data := new(eventdata.ClusterBootstrapTokenCreated)
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}
		a.bootstrapToken = data.GetJwt()
	case events.ClusterUpdated:
		data := new(eventdata.ClusterUpdated)
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}

		if len(data.GetDisplayName()) > 0 && a.displayName != data.GetDisplayName() {
			a.displayName = data.GetDisplayName()
		}
		if len(data.GetApiServerAddress()) > 0 && a.apiServerAddr != data.GetApiServerAddress() {
			a.apiServerAddr = data.GetApiServerAddress()
		}
		if len(data.GetCaCertificateBundle()) > 0 && !bytes.Equal(a.caCertBundle, data.GetCaCertificateBundle()) {
			a.caCertBundle = data.GetCaCertificateBundle()
		}
	case events.ClusterDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
