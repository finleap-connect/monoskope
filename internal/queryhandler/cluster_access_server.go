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
	"time"

	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

// clusterAccessServer is the implementation of the ClusterAccessService API
type clusterAccessServer struct {
	api.UnimplementedClusterAccessServer

	bindingRepo repositories.ReadOnlyTenantClusterBindingRepository
	userRepo    repositories.ReadOnlyUserRepository
	tenantRepo  repositories.ReadOnlyTenantUserRepository
	clusterRepo repositories.ReadOnlyClusterRepository
}

// NewClusterServiceServer returns a new configured instance of clusterServiceServer
func NewClusterAccessServer(clusterRepo repositories.ReadOnlyClusterRepository) *clusterAccessServer {
	return &clusterAccessServer{
		clusterRepo: clusterRepo,
	}
}

func NewClusterAccessClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.ClusterAccessClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewClusterAccessClient(conn), nil
}
