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
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/eventformatter"
	"github.com/finleap-connect/monoskope/pkg/domain"
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

	// TODO: how to deal with this (need access to all repos/stores)
	// 	- AuditLogRepo that have deps on all repos
	// 	- or
	store esApi.EventStoreClient
	//qhDomain domain.QueryHandlerDomain

	eventFormatter *eventformatter.EventFormatter
}

// NewAuditLogServer returns a new configured instance of auditLogServer
func NewAuditLogServer(store esApi.EventStoreClient, qhDomain domain.QueryHandlerDomain) *auditLogServer {
	return &auditLogServer{
		store:          store,
		eventFormatter: eventformatter.NewEventFormatter(qhDomain),
	}
}

func NewAuditLogClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, doApi.AuditLogClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, doApi.NewAuditLogClient(conn), nil
}

func (s *auditLogServer) GetByDateRange(request *doApi.GetAuditLogByDateRangeRequest, stream doApi.AuditLog_GetByDateRangeServer) error {
	ctx := context.Background()

	eventFilter := &esApi.EventFilter{MinTimestamp: request.MinTimestamp, MaxTimestamp: request.MaxTimestamp}
	events, err := s.store.Retrieve(ctx, eventFilter)
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

		hre := s.eventFormatter.NewHumanReadableEvent(ctx, e)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}

		err = stream.Send(hre)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}

	return nil
}