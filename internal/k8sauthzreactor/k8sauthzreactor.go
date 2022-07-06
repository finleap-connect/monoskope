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

package k8sauthzreactor

import (
	"context"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	// v1 "k8s.io/api/rbac/v1"
)

const reactor_name = "k8s-authz-reactor"

type k8sAuthZReactor struct {
	log logger.Logger
}

// NewK8sAuthZReactor creates a new k8sAuthZReactor.
func NewK8sAuthZReactor() es.Reactor {
	return &k8sAuthZReactor{
		log: logger.WithName(reactor_name),
	}
}

// HandleEvent handles a given event returns 0..* Events in reaction or an error
func (r *k8sAuthZReactor) HandleEvent(ctx context.Context, event es.Event, eventsChannel chan<- es.Event) error {
	_, err := users.CreateUserContext(ctx, users.NewSystemUser(reactor_name))
	if err != nil {
		return err
	}

	switch event.EventType() {
	case events.ClusterCreated:
		data := &eventdata.ClusterCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		return nil
	case events.ClusterCreatedV2:
		data := &eventdata.ClusterCreatedV2{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		return nil
	}

	return nil
}
