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
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/k8sauthzreactor"
	"github.com/finleap-connect/monoskope/internal/messagebus"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/eventhandler"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

var (
	eventStoreAddr string
	msgbusPrefix   string
)

var serveCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("serve-cmd")
		ctx := context.Background()

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		esConnection, esClient, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		if err != nil {
			return err
		}
		defer esConnection.Close()

		// Init message bus consumer
		log.Info("Setting up message bus consumer...")
		msgBus, err := messagebus.NewEventBusConsumer("k8s-authz-reactor", msgbusPrefix)
		if err != nil {
			return err
		}
		defer msgBus.Close()

		// Set up reactor
		reactorEventHandler := eventhandler.NewReactorEventHandler(esClient, k8sauthzreactor.NewK8sAuthZReactor())
		defer reactorEventHandler.Stop()

		// Register event handler with event bus
		if err := msgBus.AddWorker(ctx,
			reactorEventHandler,
			"k8s-authz-reactor",
			msgBus.Matcher().MatchAggregateType(aggregates.User),
			msgBus.Matcher().MatchAggregateType(aggregates.Tenant),
			msgBus.Matcher().MatchAggregateType(aggregates.UserRoleBinding),
			msgBus.Matcher().MatchAggregateType(aggregates.Cluster),
			msgBus.Matcher().MatchAggregateType(aggregates.TenantClusterBinding),
		); err != nil {
			return err
		}

		// Add readiness check
		health := healthcheck.NewHandler()
		health.AddReadinessCheck("ready", func() error { return nil })

		listener, err := net.Listen("tcp", "0.0.0.0:8086")
		if err != nil {
			return err
		}
		defer listener.Close()

		shutdown := util.NewShutdownWaitGroup()

		// Start routine waiting for signals
		shutdown.RegisterSignalHandler(func() {
			// Stop the HTTP servers
			reactorEventHandler.Stop()
			util.PanicOnError(listener.Close())
			util.PanicOnError(msgBus.Close())
		})

		err = http.Serve(listener, health)
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
	flags.StringVar(&eventStoreAddr, "event-store-api-addr", ":8081", "Address the eventstore gRPC service is listening on")
	flags.StringVar(&msgbusPrefix, "msgbus-routing-key-prefix", "m8", "Prefix for all messages emitted to the msg bus")
}
