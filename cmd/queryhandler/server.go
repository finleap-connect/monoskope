package main

import (
	"context"

	commonApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	qhApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/util"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr        string
	metricsAddr    string
	keepAlive      bool
	eventStoreAddr string
	msgbusPrefix   string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		ctx := context.Background()

		// Create EventStore client
		esConnection, esClient, err := util.NewEventStoreClient(eventStoreAddr)
		if err != nil {
			return err
		}
		defer esConnection.Close()

		// init message bus consumer
		ebConsumer, err := util.NewEventBusConsumer("queryhandler", msgbusPrefix)
		if err != nil {
			return err
		}
		defer ebConsumer.Close()

		// Setup domain
		userRepo, err := util.SetupQueryHandlerDomain(ctx, ebConsumer, esClient)
		if err != nil {
			return err
		}

		// Create gRPC server and register implementation+
		grpcServer := grpcUtil.NewServer("queryhandler-grpc", keepAlive)
		grpcServer.RegisterService(func(s grpc.ServiceRegistrar) {
			qhApi.RegisterTenantServiceServer(s, queryhandler.NewTenantServiceServer())
			qhApi.RegisterUserServiceServer(s, queryhandler.NewUserServiceServer(userRepo))
			commonApi.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
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
	flags.StringVar(&msgbusPrefix, "msgbus-routing-key-prefix", "m8", "Prefix for all messages emitted to the msg bus")
}
