package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	ggrpc "google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr          string
	metricsAddr      string
	keepAlive        bool
	eventStoreAddr   string
	queryHandlerAddr string
	enableSuperusers bool
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		ctx := context.Background()
		log := logger.WithName("server-cmd")

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		conn, esClient, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Create UserService client
		log.Info("Connecting query handler...", "queryHandlerAddr", queryHandlerAddr)
		conn, userSvcClient, err := queryhandler.NewUserServiceClient(ctx, queryHandlerAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Create TenantService client
		conn, tenantSvcClient, err := queryhandler.NewTenantServiceClient(ctx, queryHandlerAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Setup domain
		var cmdRegistry es.CommandRegistry
		log.Info("Seting up es/cqrs...")
		if enableSuperusers {
			if u := strings.Split(os.Getenv("SUPERUSERS"), ","); len(u) != 0 {
				cmdRegistry, err = domain.SetupCommandHandlerDomain(ctx, userSvcClient, tenantSvcClient, esClient, u...)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("no valid list of superusers provided")
			}
		} else {
			cmdRegistry, err = domain.SetupCommandHandlerDomain(ctx, userSvcClient, tenantSvcClient, esClient)
			if err != nil {
				return err
			}
		}

		// Create gRPC server and register implementation
		log.Info("Creating gRPC server...")
		grpcServer := grpc.NewServer("commandhandler-grpc", keepAlive)

		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterCommandHandlerServer(s, commandhandler.NewApiServer(cmdRegistry))
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
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
	flags.StringVar(&queryHandlerAddr, "query-handler-api-addr", ":8081", "Address the queryhandler gRPC service is listening on")
	flags.BoolVar(&enableSuperusers, "enable-superusers", false, "Enable superuser functionality for initial system admin creation")
}
