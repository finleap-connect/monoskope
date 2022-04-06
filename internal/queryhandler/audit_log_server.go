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
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit"
	"github.com/finleap-connect/monoskope/pkg/audit/eventformatter"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
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
func NewAuditLogServer(esClient esApi.EventStoreClient, efRegistry eventformatter.EventFormatterRegistry, userRepo repositories.ReadOnlyUserRepository, tenantRepo repositories.ReadOnlyTenantRepository, clusterRepo repositories.ReadOnlyClusterRepository) *auditLogServer {
	return &auditLogServer{
		esClient:       esClient,
		auditFormatter: audit.NewAuditFormatter(esClient, efRegistry, userRepo, tenantRepo, clusterRepo),
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
	ctx := context.Background()

	eventFilter := &esApi.EventFilter{MinTimestamp: request.MinTimestamp, MaxTimestamp: request.MaxTimestamp}
	events, err := s.esClient.Retrieve(ctx, eventFilter)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for {
		e, err := events.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}

		hre := s.auditFormatter.NewHumanReadableEvent(ctx, e)
		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}

// GetByUser returns human-readable events caused by others actions on the given user
func (s *auditLogServer) GetByUser(request *doApi.GetByUserRequest, stream doApi.AuditLog_GetByUserServer) error {
	ctx := context.Background()

	user, err := s.userRepo.ByEmail(ctx, request.Email.GetValue())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	// TODO: this retrieves create/update/delete user events however
	// 	roles events has the effected aggregate id in their data
	// 	which means retrieving all userRoleBinding events extracting the data and filtering the ones apples to the user
	//	this also mean user events first than roles events second which then need to be sorted by the timestamp
	//	simply awful and cannot be implemented right now without refactoring and extending the way events are aggregated
	events, err := s.esClient.Retrieve(ctx, &esApi.EventFilter{
		AggregateId: wrapperspb.String(user.Id),
	})
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for {
		e, err := events.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}

		hre := s.auditFormatter.NewHumanReadableEvent(ctx, e)
		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}

// GetUserActions returns human-readable events caused by the given user actions
func (s *auditLogServer) GetUserActions(request *doApi.GetUserActionsRequest, stream doApi.AuditLog_GetUserActionsServer) error {
	ctx := context.Background()

	if request.DateRange.MaxTimestamp.AsTime().Sub(request.DateRange.MinTimestamp.AsTime()) > time.Hour*24*365 {
		return errors.TranslateToGrpcError(errors.ErrInvalidArgument("date range cannot exceed one year")) // see PR #90
	}

	user, err := s.userRepo.ByEmail(ctx, request.Email.GetValue())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	events, err := s.esClient.Retrieve(ctx, &esApi.EventFilter{
		MinTimestamp: request.DateRange.MinTimestamp,
		MaxTimestamp: request.DateRange.MaxTimestamp,
	})
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for {
		e, err := events.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
		if e.Metadata[auth.HeaderAuthId] != user.Id {
			continue
		}

		hre := s.auditFormatter.NewHumanReadableEvent(ctx, e)
		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}

// GetUsersOverview returns all users, tenants/clusters they belong to, and their roles
func (s *auditLogServer) GetUsersOverview(request *doApi.GetAllRequest, stream doApi.AuditLog_GetUsersOverviewServer) error {
	ctx := context.Background()

	users, err := s.userRepo.GetAll(ctx, request.IncludeDeleted)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, user := range users {
		uo := s.auditFormatter.NewUserOverview(ctx, user)
		err = stream.Send(uo)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}
