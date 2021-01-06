package main

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

var (
	apiAddr      string
	metricsAddr  string
	keepAlive    bool
	dbAddr       string
	dbName       string
	msgbusPrefix string
	msgbusAddr   string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

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
		conf := storage.NewPostgresStoreConfig(dbName, dbAddr)
		err = conf.ConfigureTLS()
		if err != nil {
			return err
		}

		store, err := storage.NewPostgresEventStore(conf)
		if err != nil {
			return err
		}

		// init message bus publisher
		msgbusUrl := fmt.Sprintf("amqps://%s/", msgbusAddr)
		rabbitConf := messaging.NewRabbitEventBusConfig("event-store", msgbusUrl)

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
	flags.StringVar(&dbAddr, "db-addr", "127.0.0.1:26257", "DB host:port")
	flags.StringVar(&msgbusAddr, "msgbus-addr", "127.0.0.1:5672", "MessageBus host:port")
	flags.StringVar(&msgbusPrefix, "msgbus-routing-key-prefix", "m8", "Prefix for all messages emitted to the msg bus")
}
