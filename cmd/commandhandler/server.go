package main

import (
	"context"
	"net"
	"time"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr        string
	metricsAddr    string
	keepAlive      bool
	eventStoreAddr string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Create EventStore client
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		conn, err := grpc.
			NewGrpcConnectionFactory(eventStoreAddr).
			WithInsecure().
			WithRetry().
			WithBlock().
			Build(ctx)
		if err != nil {
			return err
		}
		esClient := api_es.NewEventStoreClient(conn)

		// Gateway API server
		commandHandlerApiServer := commandhandler.NewApiServer(esClient)

		// Create gRPC server and register implementation
		grpcServer := grpc.NewServer("commandhandler-grpc", keepAlive)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterCommandHandlerServer(s, commandHandlerApiServer)
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
		})

		// Setup grpc listener
		apiLis, err := net.Listen("tcp", apiAddr)
		if err != nil {
			return err
		}
		defer apiLis.Close()

		// Setup metrics listener
		metricsLis, err := net.Listen("tcp", metricsAddr)
		if err != nil {
			return err
		}
		defer metricsLis.Close()

		// Finally start the server
		return grpcServer.Serve(apiLis, metricsLis)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// Local flags
	flags := serverCmd.Flags()
	flags.BoolVar(&keepAlive, "keep-alive", false, "If enabled, gRPC will use keepalive and allow long lasting connections")
	flags.StringVarP(&apiAddr, "api-addr", "a", ":8080", "Address the gRPC service will listen on")
	flags.StringVar(&metricsAddr, "metrics-addr", ":9102", "Address the metrics http service will listen on")
	flags.StringVar(&eventStoreAddr, "event-store-api-addr", ":8081", "Address the eventstore gRPC service is listening on")
}
