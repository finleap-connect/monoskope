package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	ggrpc "google.golang.org/grpc"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/sync/errgroup"
)

var (
	grpcApiAddr      string
	httpApiAddr      string
	queryHandlerAddr string
	metricsAddr      string
	keyCacheDuration string
	keepAlive        bool
	authConfig       = auth.Config{}
	scopes           string
	redirectUris     string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.WithName("serverCmd")

		log.Info("Reading environment...")
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

		if len(scopes) > 0 {
			authConfig.Scopes = strings.Split(scopes, ",")
		}
		if len(authConfig.Scopes) == 0 {
			return fmt.Errorf("scopes must not be empty")
		}

		if len(redirectUris) > 0 {
			authConfig.RedirectURIs = strings.Split(redirectUris, ",")
		}
		if len(authConfig.RedirectURIs) == 0 {
			return fmt.Errorf("redirectUris must not be empty")
		}

		// Create token signer/validator
		keyCacheDuration, err := time.ParseDuration(keyCacheDuration)
		if err != nil {
			return err
		}

		log.Info("Configuring JWT signing and verifying...")
		signer := jwt.NewSigner("/etc/gateway/jwt/tls.key")
		verifier, err := jwt.NewVerifier("/etc/gateway/jwt/tls.crt", keyCacheDuration)
		if err != nil {
			return err
		}

		// Create interceptor for auth
		authHandler := auth.NewHandler(&authConfig, signer, verifier)

		// Setup OIDC
		if err := authHandler.SetupOIDC(cmd.Context()); err != nil {
			return err
		}

		// Create UserService client
		conn, userSvcClient, err := queryhandler.NewUserClient(context.Background(), queryHandlerAddr)
		if err != nil {
			return err
		}
		defer conn.Close()
		userRepo := repositories.NewRemoteUserRepository(userSvcClient)

		authServer := gateway.NewAuthServer(authHandler, userRepo)

		// Gateway API server
		gws := gateway.NewApiServer(&authConfig, authHandler, userRepo)

		// Create gRPC server and register implementation
		grpcServer := grpc.NewServer("gateway-grpc", keepAlive)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterGatewayServer(s, gws)
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
		})

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
	flags.StringVar(&queryHandlerAddr, "query-handler-api-addr", ":8081", "Address the queryhandler gRPC service is listening on")
	flags.StringVar(&metricsAddr, "metrics-addr", ":9102", "Address the metrics http service will listen on")
	flags.StringVar(&authConfig.IdentityProviderName, "identity-provider-name", "", "Identity provider name")
	flags.StringVar(&authConfig.IdentityProvider, "identity-provider-url", "", "Identity provider URL")
	flags.StringVar(&scopes, "scopes", "openid, profile, email", "Issuer scopes to request")
	flags.StringVar(&redirectUris, "redirect-uris", "localhost:8000,localhost18000", "Issuer allowed redirect uris")
	flags.StringVar(&keyCacheDuration, "key-cache-duration", "24h", "Cache duration of public keys for token verification")

	util.PanicOnError(serverCmd.MarkFlagRequired("identity-provider-name"))
	util.PanicOnError(serverCmd.MarkFlagRequired("identity-provider-url"))
}
