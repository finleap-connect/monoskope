// Copyright 2022 Monoskope Authors
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

package main

import (
	"github.com/finleap-connect/monoskope/internal/gateway"
	"github.com/finleap-connect/monoskope/internal/telemetry"
	api_domain "github.com/finleap-connect/monoskope/pkg/api/domain"
	api_common "github.com/finleap-connect/monoskope/pkg/api/domain/common"
	api "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/grpc/middleware/auth"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	ggrpc "google.golang.org/grpc"

	"github.com/finleap-connect/monoskope/internal/commandhandler"
	"github.com/finleap-connect/monoskope/internal/common"
	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr        string
	metricsAddr    string
	keepAlive      bool
	eventStoreAddr string
	gatewayAddr    string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		ctx := cmd.Context()
		log := logger.WithName("server-cmd")

		// Enable OpenTelemetry optionally
		log.Info("Initializing open telemetry...")
		shutdownTelemetry, err := telemetry.InitOpenTelemetry(ctx)
		if err != nil && err != telemetry.ErrOpenTelemetryNotEnabled {
			return err
		}
		if shutdownTelemetry != nil {
			defer util.PanicOnError(shutdownTelemetry())
		}

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		conn, esClient, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Setup domain
		log.Info("Seting up es/cqrs...")
		err = domain.SetupCommandHandlerDomain(ctx, esClient)
		if err != nil {
			return err
		}

		// Create Gateway Auth client
		log.Info("Connecting gateway...", "gatewayAddr", gatewayAddr)
		conn, gatewaySvcClient, err := gateway.NewInsecureAuthServerClient(ctx, gatewayAddr)
		if err != nil {
			return err
		}
		defer conn.Close()
		authMiddleware := auth.NewAuthMiddleware(gatewaySvcClient, []string{"/grpc.health.v1.Health/Check"})

		// Create gRPC server and register implementation
		log.Info("Creating gRPC server...")
		grpcServer := grpc.NewServerWithOpts("commandhandler-grpc", keepAlive,
			[]ggrpc.UnaryServerInterceptor{
				authMiddleware.UnaryServerInterceptor(),
			}, []ggrpc.StreamServerInterceptor{
				authMiddleware.StreamServerInterceptor(),
			},
		)

		commandHandlerApiServer := commandhandler.NewApiServer(es.DefaultCommandRegistry)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterCommandHandlerServer(s, commandHandlerApiServer)
			api_domain.RegisterCommandHandlerExtensionsServer(s, commandHandlerApiServer)
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
	flags.StringVar(&gatewayAddr, "gateway-api-addr", ":8081", "Address the gateway gRPC service is listening on")
}
