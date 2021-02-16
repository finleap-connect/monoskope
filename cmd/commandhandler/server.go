package main

import (
	"context"

	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/util"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr          string
	metricsAddr      string
	keepAlive        bool
	eventStoreAddr   string
	queryHandlerAddr string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		ctx := context.Background()

		// Create EventStore client
		conn, esClient, err := util.NewEventStoreClient(eventStoreAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Create UserService client
		conn, userSvcClient, err := util.NewUserServiceClient(queryHandlerAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Setup domain
		err = util.SetupCommandHandlerDomain(ctx, userSvcClient, esClient)
		if err != nil {
			return err
		}

		// Create gRPC server and register implementation
		grpcServer := grpc.NewServer("commandhandler-grpc", keepAlive)

		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterCommandHandlerServer(s, commandhandler.NewApiServer(esClient))
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
		})

		// Finally start the server
		return grpcServer.Serve(apiAddr, metricsAddr)
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
	flags.StringVar(&eventStoreAddr, "query-handler-api-addr", ":8081", "Address the queryhandler gRPC service is listening on")
}
