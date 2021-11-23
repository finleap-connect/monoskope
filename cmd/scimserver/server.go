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

	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

var (
	eventStoreAddr    string
	msgbusPrefix      string
	certIssuer        string
	certIssuerKind    string
	certDuration      string
	jwtPrivateKeyFile string
	issuerURL         string
)

var serveCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the SCIM server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("serve-cmd")
		ctx := context.Background()

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		esConnection, _, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		if err != nil {
			return err
		}
		defer esConnection.Close()

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
			util.PanicOnError(listener.Close())
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
	flags.StringVar(&jwtPrivateKeyFile, "jwt-privatekey", "/etc/clusterbootstrapreactor/signing.key", "Path to the private key for signing JWTs")

	flags.StringVarP(&certDuration, "certificate-duration", "d", "48h", "Certificate validity to request certificates for")
	flags.StringVarP(&certIssuerKind, "certificate-issuer-kind", "k", "Issuer", "Certificate issuer kind to request certificates from")

	flags.StringVarP(&certIssuer, "certificate-issuer", "i", "", "Certificate issuer name to request certificates from")
	util.PanicOnError(cobra.MarkFlagRequired(flags, "certificate-issuer"))

	flags.StringVar(&issuerURL, "issuer-url", "", "The URL of the Monoskope issuer (Gateway)")
	util.PanicOnError(cobra.MarkFlagRequired(flags, "issuer-url"))
}
