// Copyright 2021 Monoskope Authors
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
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/finleap-connect/monoskope/internal/commandhandler"
	"github.com/finleap-connect/monoskope/internal/scimserver"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/sync/errgroup"
)

var (
	commandHandlerAddr string
)

var serveCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the SCIM server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("serve-cmd")
		ctx := context.Background()

		// Create CommandHandler client
		log.Info("Connecting command handler...", "commandHandlerAddr", commandHandlerAddr)
		conn, _, err := commandhandler.NewServiceClient(ctx, commandHandlerAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Add readiness check
		health := healthcheck.NewHandler()
		health.AddReadinessCheck("ready", func() error { return nil })

		healthListener, err := net.Listen("tcp", "0.0.0.0:8086")
		if err != nil {
			return err
		}
		defer healthListener.Close()

		scimListener, err := net.Listen("tcp", "0.0.0.0:8080")
		if err != nil {
			return err
		}
		defer scimListener.Close()

		shutdown := util.NewShutdownWaitGroup()

		providerConfig := scimserver.NewProvierConfig()
		userHandler := scimserver.NewUserHandler()
		groupHandler := scimserver.NewGroupHandler()
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
	flags.StringVar(&commandHandlerAddr, "command-handler-api-addr", ":8081", "Address the command handler gRPC service is listening on")
}