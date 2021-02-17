package main

import (
	"context"

	qhApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	commonApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
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
		log := logger.WithName("server-cmd")
		ctx := context.Background()

		// Create EventStore client
		log.Info("Connectin event store...", "eventStoreAddr", eventStoreAddr)
		esConnection, esClient, err := util.NewEventStoreClient(eventStoreAddr)
		if err != nil {
			return err
		}
		defer esConnection.Close()

		// init message bus consumer
		log.Info("Setting up message bus consumer...")
		ebConsumer, err := util.NewEventBusConsumer("queryhandler", msgbusPrefix)
		if err != nil {
			return err
		}
		defer ebConsumer.Close()

		// Setup domain
		log.Info("Seting up es/cqrs...")
		userRepo, err := util.SetupQueryHandlerDomain(ctx, ebConsumer, esClient)
		if err != nil {
			return err
		}

		// Create gRPC server and register implementation+
		log.Info("Creating gRPC server...")
		grpcServer := grpcUtil.NewServer("queryhandler-grpc", keepAlive)
		grpcServer.RegisterService(func(s grpc.ServiceRegistrar) {
			qhApi.RegisterTenantServiceServer(s, queryhandler.NewTenantServiceServer())
			qhApi.RegisterUserServiceServer(s, queryhandler.NewUserServiceServer(userRepo))
			commonApi.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
		})

		// Finally start the server
		log.Info("gRPC server start serving...")
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
