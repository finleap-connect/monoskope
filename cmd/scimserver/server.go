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
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/finleap-connect/monoskope/internal/scimserver"
	"github.com/finleap-connect/monoskope/internal/telemetry"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	commandHandlerApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/sync/errgroup"
)

var (
	httpApiAddr        string
	healthApiAddr      string
	commandHandlerAddr string
	queryHandlerAddr   string
)

var serveCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the SCIM server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("serve-cmd")
		ctx := cmd.Context()

		// Enable OpenTelemetry optionally
		if telemetry.GetIsOpenTelemetryEnabled() {
			log.Info("Initializing open telemetry...")
			shutdownTelemetry, err := telemetry.InitOpenTelemetry(ctx)
			if err != nil {
				return err
			}
			defer util.PanicOnError(shutdownTelemetry())
		}

		// Create CommandHandler client
		log.Info("Connecting command handler...", "commandHandlerAddr", commandHandlerAddr)
		conn, commandHandlerClient, err := grpcUtil.NewClientWithAuthForward(ctx, commandHandlerAddr, false, commandHandlerApi.NewCommandHandlerClient)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Create User client
		log.Info("Connecting queryhandler...", "queryHandlerAddr", queryHandlerAddr)
		conn, userClient, err := grpcUtil.NewClientWithAuthForward(ctx, queryHandlerAddr, false, domainApi.NewUserClient)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Add readiness check
		health := healthcheck.NewHandler()
		health.AddReadinessCheck("ready", func() error { return nil })

		healthListener, err := net.Listen("tcp", healthApiAddr)
		if err != nil {
			return err
		}
		defer healthListener.Close()

		// Set up SCIM server
		log.Info("Setting up SCIM server...")

		scimListener, err := net.Listen("tcp", httpApiAddr)
		if err != nil {
			return err
		}
		defer scimListener.Close()

		shutdown := util.NewShutdownWaitGroup()

		providerConfig := scimserver.NewProvierConfig()
		userHandler := scimserver.NewUserHandler(commandHandlerClient, userClient)
		groupHandler := scimserver.NewGroupHandler(commandHandlerClient, userClient)
		scimServer := scimserver.NewServer(providerConfig, userHandler, groupHandler)

		// Start routine waiting for signals
		shutdown.RegisterSignalHandler(func() {
			util.PanicOnError(healthListener.Close())
			util.PanicOnError(scimListener.Close())
		})

		// Finally start the servers
		eg, _ := errgroup.WithContext(cmd.Context())
		eg.Go(func() error {
			return http.Serve(healthListener, health)
		})
		eg.Go(func() error {
			return http.Serve(scimListener, scimServer)
		})
		log.Info("Ready!")
		err = eg.Wait()

		if !shutdown.IsExpected() && err != nil {
			panic(fmt.Sprintf("shutdown unexpected: %v", err))
		}

		// Check if we are expecting shutdown
		// Wait for both shutdown signals and close the channel
		if ok := shutdown.WaitOrTimeout(30 * time.Second); !ok {
			panic("shutting down gracefully exceeded 30 seconds")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Local flags
	flags := serveCmd.Flags()
	flags.StringVar(&httpApiAddr, "http-api-addr", ":8081", "Address the HTTP service will listen on")
	flags.StringVar(&healthApiAddr, "health-api-addr", ":8082", "Address the health check HTTP service will listen on")
	flags.StringVar(&commandHandlerAddr, "command-handler-api-addr", ":8081", "Address the command handler gRPC service is listening on")
	flags.StringVar(&queryHandlerAddr, "query-handler-api-addr", ":8082", "Address the query handler gRPC service is listening on")
}
