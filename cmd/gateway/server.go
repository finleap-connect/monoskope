package main

import (
	"os"

	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/sync/errgroup"
)

var (
	grpcApiAddr string
	httpApiAddr string
	metricsAddr string
	keepAlive   bool
	authConfig  = auth.Config{}
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		// Create interceptor for auth
		authHandler := auth.NewHandler(&authConfig)

		// Setup OIDC
		if err := authHandler.SetupOIDC(cmd.Context()); err != nil {
			return err
		}

		// Gateway API server
		gws := gateway.NewApiServer(&authConfig, authHandler)

		// Create gRPC server and register implementation
		grpcServer := grpc.NewServer("gateway-grpc", keepAlive)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterGatewayServer(s, gws)
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
		})

		authServer := gateway.NewAuthServer(authHandler)

		// Finally start the servers
		eg, _ := errgroup.WithContext(cmd.Context())
		eg.Go(func() error {
			return grpcServer.Serve(grpcApiAddr, metricsAddr)
		})
		eg.Go(func() error {
			return authServer.Serve(httpApiAddr)
		})
		return eg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// Local flags
	flags := serverCmd.Flags()
	flags.BoolVar(&keepAlive, "keep-alive", false, "If enabled, gRPC will use keepalive and allow long lasting connections")
	flags.StringVar(&grpcApiAddr, "grpc-api-addr", ":8080", "Address the gRPC service will listen on")
	flags.StringVar(&httpApiAddr, "http-api-addr", ":8081", "Address the HTTP service will listen on")
	flags.StringVar(&metricsAddr, "metrics-addr", ":9102", "Address the metrics http service will listen on")
	flags.StringVar(&authConfig.IssuerURL, "issuer-url", "http://localhost:6555", "Issuer URL")
}
