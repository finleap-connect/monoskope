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

package queryhandler

import (
	"context"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	auditApi "github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters/audit"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"strings"
	"time"

	doApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

// auditLogServer is the implementation of the auditLogService API
type auditLogServer struct {
	doApi.UnimplementedAuditLogServer

	esClient       esApi.EventStoreClient
	auditFormatter audit.AuditFormatter
	userRepo       repositories.ReadOnlyUserRepository
}

// NewAuditLogServer returns a new configured instance of auditLogServer
func NewAuditLogServer(esClient esApi.EventStoreClient, efRegistry event.EventFormatterRegistry, userRepo repositories.ReadOnlyUserRepository) *auditLogServer {
	return &auditLogServer{
		esClient:       esClient,
		auditFormatter: audit.NewAuditFormatter(esClient, efRegistry),
		userRepo:       userRepo,
	}
}

// NewAuditLogClient returns a new configured instance of AuditLogClient along with the connection used
func NewAuditLogClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, doApi.AuditLogClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, doApi.NewAuditLogClient(conn), nil
}

// GetByDateRange returns human-readable events within the specified data range
func (s *auditLogServer) GetByDateRange(request *doApi.GetAuditLogByDateRangeRequest, stream doApi.AuditLog_GetByDateRangeServer) error {
	eventFilter := &esApi.EventFilter{MinTimestamp: request.MinTimestamp, MaxTimestamp: request.MaxTimestamp}
	eventsStream, err := s.esClient.Retrieve(stream.Context(), eventFilter)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for {
		e, err := eventsStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}

		hre := s.auditFormatter.NewHumanReadableEvent(stream.Context(), e)
		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}

// GetByUser returns human-readable events caused by others actions on the given user
func (s *auditLogServer) GetByUser(request *doApi.GetByUserRequest, stream doApi.AuditLog_GetByUserServer) error {
	user, err := s.userRepo.ByEmail(stream.Context(), request.Email.GetValue())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	userEventsStream, err := s.esClient.Retrieve(stream.Context(), &esApi.EventFilter{AggregateId: wrapperspb.String(user.Id)})
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	
	for {
		e, err := userEventsStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}

		hre := s.auditFormatter.NewHumanReadableEvent(stream.Context(), e)
		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	rolesEventsStream, err := s.esClient.Retrieve(stream.Context(), &esApi.EventFilter{AggregateType: wrapperspb.String(aggregates.UserRoleBinding.String())})
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for {
		e, err := rolesEventsStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}

		hre := s.auditFormatter.NewHumanReadableEvent(stream.Context(), e)
		if !strings.Contains(hre.Details, user.Email) {
			continue
		}
		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}

// GetUserActions returns human-readable events caused by the given user actions
func (s *auditLogServer) GetUserActions(request *doApi.GetUserActionsRequest, stream doApi.AuditLog_GetUserActionsServer) error {
	if request.DateRange.MaxTimestamp.AsTime().Sub(request.DateRange.MinTimestamp.AsTime()) > time.Hour*24*365 {
		return errors.TranslateToGrpcError(errors.ErrInvalidArgument("date range cannot exceed one year")) // see PR #90
	}

	user, err := s.userRepo.ByEmail(stream.Context(), request.Email.GetValue())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	eventFilter := &esApi.EventFilter{MinTimestamp: request.DateRange.MinTimestamp, MaxTimestamp: request.DateRange.MaxTimestamp}
	eventsStream, err := s.esClient.Retrieve(stream.Context(), eventFilter)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for {
		e, err := eventsStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
		if e.Metadata[auth.HeaderAuthId] != user.Id {
			continue
		}

		hre := s.auditFormatter.NewHumanReadableEvent(stream.Context(), e)
		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}

// GetUsersOverview returns users overview at the specified timestamp, tenants/clusters they belong to, and their roles
func (s *auditLogServer) GetUsersOverview(request *doApi.GetUsersOverviewRequest, stream doApi.AuditLog_GetUsersOverviewServer) error {
	eventsStream, err := s.esClient.Retrieve(stream.Context(), &esApi.EventFilter{MaxTimestamp: request.Timestamp})
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for {
		e, err := eventsStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
		if e.Type != events.UserCreated.String() {
			continue
		}

		uo := s.auditFormatter.NewUserOverview(stream.Context(), uuid.MustParse(e.AggregateId), request.Timestamp.AsTime())
		err = stream.Send(uo)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}
