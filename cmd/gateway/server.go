// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	api_common "github.com/finleap-connect/monoskope/pkg/api/domain/common"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	ggrpc "google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/finleap-connect/monoskope/internal/common"
	"github.com/finleap-connect/monoskope/internal/gateway"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/internal/queryhandler"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/sync/errgroup"
)

var (
	grpcApiAddr                string
	httpApiAddr                string
	queryHandlerAddr           string
	metricsAddr                string
	keepAlive                  bool
	scopes                     []string
	redirectUris               string
	k8sTokenLifetime           = make(map[string]string)
	authTokenValidity          string
	gatewayURL                 string
	identityProvider           string
	policiesPath               string
	k8sTokenLifetimeConfigPath string
	jwtPath                    string
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.WithName("serverCmd")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		authClientConfig := auth.ClientConfig{
			IdentityProvider: identityProvider,
		}
		authServerConfig := auth.ServerConfig{
			URL: gatewayURL,
		}

		log.Info("Reading environment...")
		// Some options can be provided by env variables
		if v := os.Getenv("OIDC_CLIENT_ID"); v != "" {
			authClientConfig.ClientId = v
		}
		if v := os.Getenv("OIDC_CLIENT_SECRET"); v != "" {
			authClientConfig.ClientSecret = v
		}
		if v := os.Getenv("OIDC_NONCE"); v != "" {
			authClientConfig.Nonce = v
		}

		if len(authClientConfig.Scopes) == 0 {
			return fmt.Errorf("scopes must not be empty")
		}

		if len(redirectUris) > 0 {
			authClientConfig.RedirectURIs = strings.Split(redirectUris, ",")
		}
		if len(authClientConfig.RedirectURIs) == 0 {
			return fmt.Errorf("redirectUris must not be empty")
		}

		// Create token signer/validator
		log.Info("Configuring JWT signing and verifying...")
		signer := jwt.NewSigner(path.Join(jwtPath, "tls.key"))
		verifier, err := jwt.NewVerifier(path.Join(jwtPath, "tls.crt"))
		if err != nil {
			return err
		}
		defer verifier.Close()

		// Create interceptor for auth
		authTokenValidityDuration, err := time.ParseDuration(authTokenValidity)
		if err != nil {
			return err
		}
		authServerConfig.TokenValidity = authTokenValidityDuration
		client := auth.NewClient(&authClientConfig)
		server := auth.NewServer(&authServerConfig, signer, verifier)

		// Setup OIDC
		if err := client.SetupOIDC(cmd.Context()); err != nil {
			return err
		}

		// Create UserService client
		conn, userSvcClient, err := queryhandler.NewUserClient(ctx, queryHandlerAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		conn, clusterSvcClient, err := queryhandler.NewClusterClient(ctx, queryHandlerAddr)
		if err != nil {
			return err
		}
		defer conn.Close()

		userRepo := repositories.NewRemoteUserRepository(userSvcClient)
		clusterRepo := repositories.NewRemoteClusterRepository(clusterSvcClient)

		// API servers
		authServer, err := gateway.NewAuthServer(ctx, gatewayURL, server, userRepo, policiesPath)
		if err != nil {
			return err
		}

		oidcProviderServer := gateway.NewOIDCProviderServer(server)
		gatewayApiServer := gateway.NewGatewayAPIServer(&authClientConfig, client, server, userRepo)

		// Look for config
		if len(k8sTokenLifetime) == 0 {
			data, err := ioutil.ReadFile(path.Join(ConfigPath, k8sTokenLifetimeConfigPath))
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(data, k8sTokenLifetime)
			if err != nil {
				return err
			}
		}

		// Parse token lifetime
		tokenLifeTimePerRole := make(map[string]time.Duration)
		for k, v := range k8sTokenLifetime {
			if err := k8s.ValidateRole(k); err != nil {
				return err
			}
			k8sTokenValidityDuration, err := time.ParseDuration(v)
			if err != nil {
				return err
			}
			tokenLifeTimePerRole[k] = k8sTokenValidityDuration
		}
		clusterAuthApiServer := gateway.NewClusterAuthAPIServer(gatewayURL, signer, userRepo, clusterRepo, tokenLifeTimePerRole)

		apiTokenServer := gateway.NewAPITokenServer(gatewayURL, signer, userRepo)

		// Create gRPC server and register implementation
		grpcServer := grpc.NewServer("gateway-grpc", keepAlive)
		grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterGatewayServer(s, gatewayApiServer)
			api.RegisterClusterAuthServer(s, clusterAuthApiServer)
			api.RegisterAPITokenServer(s, apiTokenServer)
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
			api.RegisterGatewayAuthZServer(s, authServer)
		})

		// Finally start the servers
		eg, _ := errgroup.WithContext(cmd.Context())
		eg.Go(func() error {
			return grpcServer.Serve(grpcApiAddr, metricsAddr)
		})
		eg.Go(func() error {
			return oidcProviderServer.Serve(httpApiAddr)
		})
		return eg.Wait()
	},
}

func init() {
	// Local flags
	flags := serverCmd.Flags()
	flags.BoolVar(&keepAlive, "keep-alive", false, "If enabled, gRPC will use keepalive and allow long lasting connections")
	flags.StringVar(&grpcApiAddr, "grpc-api-addr", ":8080", "Address the gRPC service will listen on")
	flags.StringVar(&httpApiAddr, "http-api-addr", ":8081", "Address the HTTP service will listen on")
	flags.StringVar(&queryHandlerAddr, "query-handler-api-addr", ":8081", "Address the queryhandler gRPC service is listening on")
	flags.StringVar(&metricsAddr, "metrics-addr", ":9102", "Address the metrics http service will listen on")
	flags.StringArrayVar(&scopes, "scopes", []string{"openid", "profile", "email"}, "Issuer scopes to request")
	flags.StringVar(&redirectUris, "redirect-uris", "localhost:8000,localhost18000", "Issuer allowed redirect uris")
	flags.StringVar(&k8sTokenLifetimeConfigPath, "k8s-token-lifetime-path", "/etc/gateway/k8s-auth/k8sTokenLifetime.yaml", "YAML containing the token lifetime for k8s token per role. Only used if `k8s-token-lifetime` is not specified")
	flags.StringToStringVar(&k8sTokenLifetime, "k8s-token-lifetime", k8sTokenLifetime, "Token lifetime for k8s token per role")
	flags.StringVar(&authTokenValidity, "auth-token-validity", "12h", "Validity period of m8 auth token")

	flags.StringVar(&identityProvider, "identity-provider-url", "", "Identity provider URL")
	util.PanicOnError(serverCmd.MarkFlagRequired("identity-provider-url"))

	flags.StringVar(&gatewayURL, "gateway-url", "", "URL of the gateway itself")
	util.PanicOnError(serverCmd.MarkFlagRequired("gateway-url"))

	flags.StringVar(&policiesPath, "policies-path", "/etc/gateway/policies", "Path to rego policies to authorize requests against")
	flags.StringVar(&jwtPath, "jwt-signing-verifying-path", "/etc/gateway/jwt", "Path to tls.key and tlss.cert for signing and verifying JWTs")
}
