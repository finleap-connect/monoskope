package messaging

import (
	"context"
	"time"

	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

type remoteEventBus struct {
	eventStoreAddr   string
	eventStoreClient esApi.EventStoreClient
	conn             *grpc.ClientConn
}

func NewRemoteEventBus(eventStoreAddr string) es.EventBusPublisher {
	return &remoteEventBus{
		eventStoreAddr: eventStoreAddr,
	}
}

func (b *remoteEventBus) Connect(ctx context.Context) error {
	// Create EventStore client
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpcUtil.
		NewGrpcConnectionFactory(b.eventStoreAddr).
		WithInsecure().
		WithRetry().
		WithBlock().
		Build(ctx)
	if err != nil {
		return err
	}
	b.conn = conn

	b.eventStoreClient = esApi.NewEventStoreClient(conn)

	return nil
}

func (b *remoteEventBus) PublishEvent(ctx context.Context, event es.Event) error {
	protoEvent, err := es.NewProtoFromEvent(event)
	if err != nil {
		return err
	}

	storeClient, err := b.eventStoreClient.Store(ctx)
	if err != nil {
		return err
	}

	err = storeClient.Send(protoEvent)
	if err != nil {
		return err
	}

	return storeClient.CloseSend()
}

func (b *remoteEventBus) Close() error {
	return b.conn.Close()
}
