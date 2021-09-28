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

package projectors

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain/user_repo", func() {
	ctx := context.Background()
	userId := uuid.New()
	adminUser := &projections.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}

	It("can handle events", func() {
		userProjector := NewUserProjector()
		userProjection := userProjector.NewProjection(uuid.New())
		protoEventData := &eventdata.UserCreated{
			Name:  adminUser.Name,
			Email: adminUser.Email,
		}
		eventData := eventsourcing.ToEventDataFromProto(protoEventData)
		event := eventsourcing.NewEvent(ctx, events.UserCreated, eventData, time.Now().UTC(), aggregates.User, uuid.MustParse(adminUser.Id), 1)
		event.Metadata()[auth.HeaderAuthId] = userId.String()
		userProjection, err := userProjector.Project(context.Background(), event, userProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(userProjection.Version()).To(Equal(uint64(1)))

		deleteEvent := eventsourcing.NewEvent(ctx, events.UserDeleted, nil, time.Now().UTC(), aggregates.User, uuid.MustParse(adminUser.Id), 2)
		deleteEvent.Metadata()[auth.HeaderAuthId] = userId.String()
		userProjection, err = userProjector.Project(context.Background(), deleteEvent, userProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(userProjection.Version()).To(Equal(uint64(2)))
	})
})
