package main

import (
	"net"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
)

var (
	apiAddr     string
	metricsAddr string
	keepAlive   bool
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

		// Gateway API server
		commandHandlerApiServer := commandhandler.NewApiServer()

		// Create gRPC server and register implementation
		grpcServer := grpc.NewServer("commandhandler-grpc", keepAlive)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterCommandHandlerServer(s, commandHandlerApiServer)
			api_common.RegisterServiceInformationServiceServer(s, commandHandlerApiServer)
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
}
