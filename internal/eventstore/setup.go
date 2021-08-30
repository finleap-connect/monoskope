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

package eventstore

import (
	"context"
	"os"
	"time"

	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/storage"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

func NewEventStore() (eventsourcing.EventStore, error) {
	var dbUrl string

	if v := os.Getenv("DB_URL"); v != "" {
		dbUrl = v
	}

	conf, err := storage.NewPostgresStoreConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	err = conf.ConfigureTLS()
	if err != nil {
		return nil, err
	}

	store, err := storage.NewPostgresEventStore(conf)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err = store.Open(ctx)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func NewEventStoreClient(ctx context.Context, eventStoreAddr string) (*grpc.ClientConn, esApi.EventStoreClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(eventStoreAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	return conn, esApi.NewEventStoreClient(conn), nil
}
