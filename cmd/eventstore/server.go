package main

import (
	"net"
	"os"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

var (
	apiAddr      string
	metricsAddr  string
	keepAlive    bool
	dbUrl        string
	msgbusPrefix string
	msgbusUrl    string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Some options can be provided by env variables
		if v := os.Getenv("DB_URL"); v != "" {
			dbUrl = v
		}
		if v := os.Getenv("BUS_URL"); v != "test" {
			msgbusUrl = v
		}

		// Setup grpc listener
		apiLis, err := net.Listen("tcp", apiAddr)
		if err != nil {
			return err
		}
		defer apiLis.Close()

		// Setup metrics listener
		metricsLis, err := net.Listen("tcp", metricsAddr)
		if err != nil {
			return err
		}
		defer metricsLis.Close()

		// init event store
		conf, err := storage.NewPostgresStoreConfig(dbUrl)
		if err != nil {
			return err
		}
		err = conf.ConfigureTLS()
		if err != nil {
			return err
		}
		store, err := storage.NewPostgresEventStore(conf)
		if err != nil {
			return err
		}

		// init message bus publisher
		rabbitConf := messaging.NewRabbitEventBusConfig("event-store", msgbusUrl)
		if msgbusPrefix != "" {
			rabbitConf.RoutingKeyPrefix = msgbusPrefix
		}

		err = rabbitConf.ConfigureTLS()
		if err != nil {
			return err
		}
		publisher, err := messaging.NewRabbitEventBusPublisher(rabbitConf)
		if err != nil {
			return err
		}

		// Create the server
		serverConfig := eventstore.NewServerConfig()
		serverConfig.KeepAlive = keepAlive
		serverConfig.Store = store
		serverConfig.Bus = publisher

		s, err := eventstore.NewServer(serverConfig)
		if err != nil {
			return err
		}

		// Finally start the server
		return s.Serve(apiLis, metricsLis)
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
