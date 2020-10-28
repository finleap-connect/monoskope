package main

import (
	"net"
	"os"

	"github.com/spf13/cobra"
	auth_server "gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth/server"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway"
)

var (
	apiAddr     string
	metricsAddr string
	keepAlive   bool
	dexAddr     string
	authConfig  auth_server.Config
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Some options can be provided by env variables
		if v := os.Getenv("AUTH_ROOT_TOKEN"); v != "" {
			authConfig.RootToken = &v
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

		// Create interceptor for auth
		authInterceptor, err := auth_server.NewInterceptor(&authConfig)
		if err != nil {
			return err
		}

		// Create the server
		conf := &gateway.ServerConfig{
			KeepAlive:             false,
			AuthServerInterceptor: authInterceptor,
		}

		s, err := gateway.NewServer(conf)
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
	flags.StringVar(&dexAddr, "dex-addr", "localhost:5000", "Address of dex gRPC service")
	flags.StringVar(&authConfig.IssuerURL, "issuer-url", "http://localhost:5556", "Issuer URL")
}
