package main

import (
	"net"
	"os"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"
)

var (
	apiAddr     string
	metricsAddr string
	keepAlive   bool
	authConfig  auth.Config
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

		// Create the server
		conf := &gateway.ServerConfig{
			KeepAlive:  false,
			AuthConfig: &authConfig,
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
	flags.StringVarP(&apiAddr, "api-addr", "a", ":8080", "Address the gRPC service will listen on")
	flags.StringVar(&metricsAddr, "metrics-addr", ":9102", "Address the metrics http service will listen on")
	flags.StringVar(&authConfig.IssuerURL, "issuer-url", "http://localhost:6555", "Issuer URL")
}
