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

package formatters

import (
	"context"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"time"

	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

// UserSnapshotter implements basic snapshot creation for the user aggregate
type UserSnapshotter struct {
	*Snapshotter[*projections.User]
}

func NewUserSnapshotter(esClient esApi.EventStoreClient, projector es.Projector[*projections.User]) *UserSnapshotter {
	return &UserSnapshotter{Snapshotter: &Snapshotter[*projections.User]{esClient, projector}}
}

// CreateRoleBindingSnapshots returns a list of userRoleBinding snapshots for the user specified by its id
// This is a temporary implementation until snapshots are fully implemented,
// and it is not meant to be used extensively.
func (s *UserSnapshotter) CreateRoleBindingSnapshots(ctx context.Context, userId uuid.UUID, timestamp time.Time) []*projections.UserRoleBinding {
	var userRoleBindings []*projections.UserRoleBinding
	roleBindingEvents, err := s.esClient.Retrieve(ctx, &esApi.EventFilter{
		MaxTimestamp:  timestamppb.New(timestamp),
		AggregateType: wrapperspb.String(aggregates.UserRoleBinding.String()),
	})
	if err != nil {
		return userRoleBindings
	}

	roleBindingProjector := projectors.NewUserRoleBindingProjector()
	roleBindings := make(map[string]*projections.UserRoleBinding)
	for {
		eventPorto, err := roleBindingEvents.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		e, err := es.NewEventFromProto(eventPorto)
		if err != nil {
			continue
		}

		userRoleBinding, exist := roleBindings[e.AggregateID().String()]
		if !exist {
			userRoleBinding = roleBindingProjector.NewProjection(e.AggregateID())
			roleBindings[e.AggregateID().String()] = userRoleBinding
		}

		if userRoleBinding.UserId != "" && userRoleBinding.UserId != userId.String() {
			continue // after projecting once to ensure the roleBinding is irrelevant
		}

		userRoleBinding, err = roleBindingProjector.Project(ctx, e, userRoleBinding)
		if err != nil {
			continue
		}

		if !exist && userRoleBinding.UserId == userId.String() {
			userRoleBindings = append(userRoleBindings, userRoleBinding)
		}
	}

	return userRoleBindings
}
