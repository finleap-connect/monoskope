package main

import (
	"net"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/util"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	_ "go.uber.org/automaxprocs"
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
		store, err := util.NewEventStore()
		if err != nil {
			return err
		}

		// init message bus publisher
		publisher, err := util.NewEventBusPublisher("eventstore", msgbusPrefix)
		if err != nil {
			return err
		}

		// Create the server
		eventStore, err := eventstore.NewApiServer(store, publisher)
		if err != nil {
			return err
		}

		grpcServer := grpc.NewServer("event-store-grpc", keepAlive)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterEventStoreServer(s, eventStore)
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
		})
		grpcServer.RegisterOnShutdown(func() {
			eventStore.Shutdown()
		})

		// Finally start the server
		return grpcServer.Serve(apiLis, metricsLis)
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
