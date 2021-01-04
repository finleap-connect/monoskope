package main

import (
	"net"
	"os"

	"github.com/go-pg/pg"
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

var (
	apiAddr     string
	metricsAddr string
	keepAlive   bool
	dbAddr      string
	dbUser      string
	dbName      string
	dbPassword  string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Some options can be provided by env variables
		if v := os.Getenv("DB_PASSWORD"); v != "" {
			dbPassword = v
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

		// create event store db connection
		db := pg.Connect(&pg.Options{
			Addr:     dbAddr,
			Database: dbName,
			User:     dbUser,
			Password: dbPassword,
		})

		// init event store
		store, err := storage.NewPostgresEventStore(db)
		if err != nil {
			return err
		}

		// TODO: create message bus connection

		// Create the server
		serverConfig := eventstore.ServerConfig{
			Store: store,
		}
		s, err := eventstore.NewServer(&serverConfig)
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
	flags.StringVar(&dbAddr, "db-addr", "127.0.0.1:26257", "Store db host:port")
	flags.StringVar(&dbUser, "db-user", "eventstore", "Store db user")
	flags.StringVar(&dbName, "db-name", "eventstore", "Store db name")
}
