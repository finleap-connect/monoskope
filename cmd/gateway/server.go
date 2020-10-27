package main

import (
	"net"

	dexpb "github.com/dexidp/dex/api"
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"
	"google.golang.org/grpc"
)

var (
	apiAddr     string
	metricsAddr string
	keepAlive   bool
	dexAddr     string
	authConfig  auth.Config
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

		// Connect to dex
		opts := []grpc.DialOption{grpc.WithInsecure()}
		dexConn, err := grpc.Dial(dexAddr, opts...)
		if err != nil {
			return err
		}
		defer dexConn.Close()
		dexClient := dexpb.NewDexClient(dexConn)

		// Create interceptor for auth
		authInterceptor, err := auth.NewAuthInterceptor(dexClient, &authConfig)
		if err != nil {
			return err
		}

		// Create the server
		s, err := gateway.NewServer(keepAlive, authInterceptor)
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
	flags.StringVar(&dexAddr, "dex-addr", "127.0.0.1:5557", "Address of dex gRPC service")
	flags.StringVar(&authConfig.IssuerURL, "issuer-url", "http://127.0.0.1:5556", "Issuer URL")
}
