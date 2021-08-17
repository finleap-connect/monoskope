package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/certificatemanagement"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/reactors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
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

		// Create EventStore client
		log.Info("Connecting event store...", "eventStoreAddr", eventStoreAddr)
		esConnection, esClient, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		if err != nil {
			return err
		}
		defer esConnection.Close()

		// Init message bus consumer
		log.Info("Setting up message bus consumer...")
		msgBus, err := messagebus.NewEventBusConsumer("cluster-bootstrap-reactor", msgbusPrefix)
		if err != nil {
			return err
		}
		defer msgBus.Close()

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

		// Set up JWT signer
		signer := jwt.NewSigner(jwtPrivateKeyFile)

		// Set up reactor
		reactorEventHandler := eventhandler.NewReactorEventHandler(esClient, reactors.NewClusterBootstrapReactor(issuerURL, signer, certManager))
		defer reactorEventHandler.Stop()

		// Register event handler with event bus
		if err := msgBus.AddHandler(ctx,
			reactorEventHandler,
			msgBus.Matcher().MatchEventType(events.ClusterCreated),
			msgBus.Matcher().MatchEventType(events.ClusterCreatedV2),
			msgBus.Matcher().MatchEventType(events.CertificateRequested),
		); err != nil {
			return err
		}

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
