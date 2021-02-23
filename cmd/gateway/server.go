package main

import (
	"os"

	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	api_gw "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	api_gwauth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr     string
	metricsAddr string
	keepAlive   bool
	authConfig  = auth.Config{}
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Some options can be provided by env variables
		if v := os.Getenv("OIDC_CLIENT_ID"); v != "" {
			authConfig.ClientId = v
		}
		if v := os.Getenv("OIDC_CLIENT_SECRET"); v != "" {
			authConfig.ClientSecret = v
		}
		if v := os.Getenv("OIDC_NONCE"); v != "" {
			authConfig.Nonce = v
		}

		// Create the server
		// Create interceptor for auth
		authHandler, err := auth.NewHandler(&authConfig)
		if err != nil {
			return err
		}
		authInterceptor, err := auth.NewInterceptor(authHandler)
		if err != nil {
			return err
		}

		// Gateway API server
		gws := gateway.NewApiServer(&authConfig, authHandler)

		// Create gRPC server and register implementation
		grpcServer := grpc.NewServerWithOpts("gateway-grpc", keepAlive,
			[]ggrpc.UnaryServerInterceptor{
				auth.UnaryServerInterceptor(authInterceptor.EnsureValid),
			},
			[]ggrpc.StreamServerInterceptor{
				auth.StreamServerInterceptor(authInterceptor.EnsureValid),
			})
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api_gw.RegisterGatewayServer(s, gws)
			api_gwauth.RegisterAuthServer(s, gws)
			api_common.RegisterServiceInformationServiceServer(s, gws)
		})

		// Finally start the server
		return grpcServer.Serve(apiAddr, metricsAddr)
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
