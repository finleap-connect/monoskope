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

package snapshots

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

type UserRoleBindingSnapshotter struct {
	*Snapshotter[*projections.UserRoleBinding]
}

func NewUserRoleBindingSnapshotter(esClient esApi.EventStoreClient) *UserRoleBindingSnapshotter {
	return &UserRoleBindingSnapshotter{Snapshotter: &Snapshotter[*projections.UserRoleBinding]{esClient, projectors.NewUserRoleBindingProjector()}}
}

// CreateAllSnapshots returns all userRoleBinding snapshots of the user specified by its id
func (s *UserRoleBindingSnapshotter) CreateAllSnapshots(ctx context.Context, userId uuid.UUID, timestamp time.Time) []*projections.UserRoleBinding {
	var userRoleBindings []*projections.UserRoleBinding
	roleBindingEvents, err := s.esClient.Retrieve(ctx, &esApi.EventFilter{
		MaxTimestamp:  timestamppb.New(timestamp),
		AggregateType: wrapperspb.String(aggregates.UserRoleBinding.String()),
	})
	if err != nil {
		return userRoleBindings
	}

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
			userRoleBinding = s.projector.NewProjection(e.AggregateID())
			roleBindings[e.AggregateID().String()] = userRoleBinding
		}

		if userRoleBinding.UserId != "" && userRoleBinding.UserId != userId.String() {
			continue // after projecting once to ensure the roleBinding is irrelevant
		}

		userRoleBinding, err = s.projector.Project(ctx, e, userRoleBinding)
		if err != nil {
			continue
		}

		if !exist && userRoleBinding.UserId == userId.String() {
			userRoleBindings = append(userRoleBindings, userRoleBinding)
		}
	}

	return userRoleBindings
}
