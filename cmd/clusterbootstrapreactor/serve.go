package main

import (
	"context"
	"net/http"
	"os"

	"github.com/heptiolabs/healthcheck"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	_ "go.uber.org/automaxprocs"
)

var (
	healthAddr     string
	metricsAddr    string
	eventStoreAddr string
	msgbusPrefix   string
)

var serveCmd = &cobra.Command{
	Use:   "serve [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("serve-cmd")
		ctx := context.Background()

		// Add health check
		promRegistry := prom.NewRegistry()
		healthCheckHandler := healthcheck.NewMetricsHandler(promRegistry, os.Getenv("K8S_NAMESPACE"))
		go func() {
			log.Info(http.ListenAndServe(healthAddr, healthCheckHandler).Error())
		}()

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		esConnection, _, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		if err != nil {
			return err
		}
		defer esConnection.Close()

		// Init message bus consumer
		log.Info("Setting up message bus consumer...")
		ebConsumer, err := messagebus.NewEventBusConsumer("cluster-bootstrap-reactor", msgbusPrefix)
		if err != nil {
			return err
		}
		defer ebConsumer.Close()

		// Finally start the servers
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Local flags
	flags := serveCmd.Flags()
	flags.StringVarP(&healthAddr, "addr", "a", ":8086", "Address the health and readiness check http service will listen on")
	flags.StringVar(&metricsAddr, "metrics-addr", ":9102", "Address the metrics http service will listen on")
	flags.StringVar(&eventStoreAddr, "event-store-api-addr", ":8081", "Address the eventstore gRPC service is listening on")
	flags.StringVar(&msgbusPrefix, "msgbus-routing-key-prefix", "m8", "Prefix for all messages emitted to the msg bus")
}
