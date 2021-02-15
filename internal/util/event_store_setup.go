package util

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

func NewEventStore() (eventsourcing.Store, error) {
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
	err = store.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func NewEventStoreClient(eventStoreAddr string) (*grpc.ClientConn, esApi.EventStoreClient, error) {
	// Create EventStore client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpcUtil.
		NewGrpcConnectionFactory(eventStoreAddr).
		WithInsecure().
		WithRetry().
		WithBlock().
		Build(ctx)
	if err != nil {
		return nil, nil, err
	}

	return conn, esApi.NewEventStoreClient(conn), nil
}
