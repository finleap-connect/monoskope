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
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	ggrpc "google.golang.org/grpc"
)

var (
	apiAddr      string
	metricsAddr  string
	keepAlive    bool
	msgbusPrefix string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("server-cmd")

		// init message bus publisher
		log.Info("Setting up message bus publisher...")
		publisher, err := messagebus.NewEventBusPublisher("eventstore", msgbusPrefix)
		if err != nil {
			log.Error(err, "Failed to configure message bus publisher.")
			return err
		}
		defer publisher.Close()

		// init event store
		log.Info("Setting up event store...")
		store, err := eventstore.NewEventStore()
		if err != nil {
			log.Error(err, "Failed to configure event store.")
			return err
		}
		defer store.Close()

		// Create the server
		log.Info("Creating gRPC server...")
		grpcServer := grpc.NewServer("event-store-grpc", keepAlive)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api_es.RegisterEventStoreServer(s, eventstore.NewApiServer(store, publisher))
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
	flags.StringVar(&msgbusPrefix, "msgbus-routing-key-prefix", "m8", "Prefix for all messages emitted to the msg bus")
}
