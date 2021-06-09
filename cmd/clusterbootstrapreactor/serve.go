package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/clusterbootstrapreactor"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/certificatemanagement"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	_ "go.uber.org/automaxprocs"
)

var (
	healthAddr     string
	metricsAddr    string
	eventStoreAddr string
	msgbusPrefix   string
	certIssuer     string
	certIssuerKind string
	certDuration   string
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
		k8sNamespace := os.Getenv("K8S_NAMESPACE")
		if k8sNamespace == "" {
			return errors.New("K8S_NAMESPACE env variable not set")
		}

		// Add health check handling
		ready := false
		promRegistry := prom.NewRegistry()
		healthCheckHandler := healthcheck.NewMetricsHandler(promRegistry, k8sNamespace)
		healthCheckHandler.AddReadinessCheck("setup complete", func() error {
			if !ready {
				return errors.New("starting up...")
			}
			return nil
		})
		go func() {
			log.Info(http.ListenAndServe(healthAddr, healthCheckHandler).Error())
		}()

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		esConnection, esClient, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
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

		// Set up K8s client
		k8sClient, err := k8s.NewClient()
		if err != nil {
			return err
		}

		// Set up CertificateManager
		duration, err := time.ParseDuration(certDuration)
		if err != nil {
			return err
		}
		certManager := certificatemanagement.NewCertManagerClient(k8sClient, k8sNamespace, certIssuerKind, certIssuer, duration)

		// Set up
		err = clusterbootstrapreactor.SetupClusterBootstrapReactor(ctx, ebConsumer, esClient, certManager)
		if err != nil {
			return err
		}

		ready = true

		// Wait for interrupt signal sent from terminal or on sigterm
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGQUIT)
		<-sigint

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

	flags.StringVarP(&certDuration, "certificate-duration", "d", "48h", "Certificate validity to request certificates for")
	flags.StringVarP(&certIssuerKind, "certificate-issuer-kind", "k", "Issuer", "Certificate issuer kind to request certificates from")

	flags.StringVarP(&certIssuer, "certificate-issuer", "i", "", "Certificate issuer name to request certificates from")
	util.PanicOnError(cobra.MarkFlagRequired(flags, "certificate-issuer"))
}
