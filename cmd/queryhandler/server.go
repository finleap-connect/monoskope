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
	ef "github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
	"github.com/finleap-connect/monoskope/pkg/grpc/middleware/auth"
	"github.com/finleap-connect/monoskope/pkg/util"

	qhApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	commonApi "github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/domain"
	grpc "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/logger"
	ggrpc "google.golang.org/grpc"

	"github.com/finleap-connect/monoskope/internal/common"
	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/gateway"
	"github.com/finleap-connect/monoskope/internal/k8sauthz"
	"github.com/finleap-connect/monoskope/internal/messagebus"
	"github.com/finleap-connect/monoskope/internal/queryhandler"
	"github.com/finleap-connect/monoskope/internal/telemetry"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr        string
	metricsAddr    string
	keepAlive      bool
	eventStoreAddr string
	msgbusPrefix   string
	gatewayAddr    string
	k8sAuthZConf   string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("server-cmd")
		ctx := cmd.Context()

		// Enable OpenTelemetry optionally
		log.Info("Initializing open telemetry...")
		shutdownTelemetry, err := telemetry.InitOpenTelemetry(ctx)
		if err != nil && err != telemetry.ErrOpenTelemetryNotEnabled {
			return err
		}
		if shutdownTelemetry != nil {
			defer util.PanicOnErrorFunc(shutdownTelemetry)
		}

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		esConnection, esClient, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		if err != nil {
			return err
		}
		defer util.PanicOnErrorFunc(esConnection.Close)

		// init message bus consumer
		log.Info("Setting up message bus consumer...")
		ebConsumer, err := messagebus.NewEventBusConsumer("queryhandler", msgbusPrefix)
		if err != nil {
			return err
		}
		defer util.PanicOnErrorFunc(ebConsumer.Close)

		// Setup domain
		log.Info("Seting up es/cqrs...")
		qhDomain, err := domain.NewQueryHandlerDomain(ctx, ebConsumer, esClient)
		if err != nil {
			return err
		}

		// Create gRPC server and register implementation
		// Create Gateway Auth client
		log.Info("Connecting gateway...", "gatewayAddr", gatewayAddr)
		conn, gatewaySvcClient, err := gateway.NewInsecureAuthServerClient(ctx, gatewayAddr)
		if err != nil {
			return err
		}
		defer util.PanicOnErrorFunc(conn.Close)

		authMiddleware := auth.NewAuthMiddleware(gatewaySvcClient, []string{"/grpc.health.v1.Health/Check"})

		// Create gRPC server and register implementation
		log.Info("Creating gRPC server...")
		grpcServer := grpc.NewServerWithOpts("queryhandler-grpc", keepAlive,
			[]ggrpc.UnaryServerInterceptor{
				authMiddleware.UnaryServerInterceptor(),
			}, []ggrpc.StreamServerInterceptor{
				authMiddleware.StreamServerInterceptor(),
			},
		)

		// Configure k8s authz reconciliation
		if k8sAuthZConf != "" {
			conf, err := k8sauthz.NewConfigFromFilePath(k8sAuthZConf)
			if err != nil {
				return err
			}
			k8sAuthZManager := k8sauthz.NewManager(qhDomain.UserRepository, qhDomain.ClusterAccessRepo)

			if err := k8sAuthZManager.Run(ctx, conf); err != nil {
				return err
			}
			defer util.PanicOnErrorFunc(k8sAuthZManager.Close)
		}

		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			qhApi.RegisterTenantServer(s, queryhandler.NewTenantServer(qhDomain.TenantRepository, qhDomain.TenantUserRepository))
			qhApi.RegisterUserServer(s, queryhandler.NewUserServer(qhDomain.UserRepository))
			qhApi.RegisterClusterServer(s, queryhandler.NewClusterServer(qhDomain.ClusterRepository))
			qhApi.RegisterClusterAccessServer(s, queryhandler.NewClusterAccessServer(qhDomain.ClusterAccessRepo, qhDomain.TenantClusterBindingRepository))
			qhApi.RegisterAuditLogServer(s, queryhandler.NewAuditLogServer(esClient, ef.DefaultEventFormatterRegistry, qhDomain.UserRepository))
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
	flags.StringVar(&gatewayAddr, "gateway-api-addr", ":8081", "Address the gateway gRPC service is listening on")
	flags.StringVar(&k8sAuthZConf, "k8s-authz-conf-path", "", "Path to load K8sAuthZ config from. If not specified the feature is disabled.")
}
