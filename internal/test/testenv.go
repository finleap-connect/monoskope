package test

import (
	"crypto/tls"
	"fmt"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
)

const (
	AuthRootToken = "super-secret-root-token"
)

type TestEnv struct {
	pool      *dockertest.Pool
	resources map[string]*dockertest.Resource
	log       logger.Logger
}

type OAuthTestEnv struct {
	*TestEnv
	DexWebEndpoint                    string
	GatewayClientTransportCredentials credentials.TransportCredentials
	GatewayTlsCert                    *tls.Certificate
	AuthConfig                        *auth.Config
}

func SetupGeneralTestEnv() *TestEnv {
	log := logger.WithName("testenv")
	env := &TestEnv{
		log:       log,
		resources: make(map[string]*dockertest.Resource),
	}
	log.Info("Setting up testenv...")
	return env
}

func SetupAuthTestEnv() (*OAuthTestEnv, error) {
	env := &OAuthTestEnv{
		TestEnv: SetupGeneralTestEnv(),
	}
	log := env.log

	log.Info("Creating docker pool...")
	pool, err := dockertest.NewPool("")
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.pool = pool

	log.Info("Spawning dex container")
	options := &dockertest.RunOptions{
		Name:       "dex",
		Repository: "quay.io/dexidp/dex",
		Tag:        "v2.25.0",
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5556": {{HostPort: "5556"}},
		},
		ExposedPorts: []string{"5556", "5000"},
		Cmd:          []string{"serve", "/etc/dex/cfg/config.yaml"},
		Mounts:       []string{fmt.Sprintf("%s:/etc/dex/cfg", DexConfigPath)},
	}
	dexContainer, err := pool.RunWithOptions(options)
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.resources[dexContainer.Container.Name] = dexContainer
	env.DexWebEndpoint = fmt.Sprintf("http://127.0.0.1:%s", dexContainer.GetPort("5556/tcp"))

	clientTransportCredentials, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.GatewayClientTransportCredentials = clientTransportCredentials

	cert, err := tls.LoadX509KeyPair(data.Path("x509/server_cert.pem"), data.Path("x509/server_key.pem"))
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.GatewayTlsCert = &cert

	rootToken := AuthRootToken
	env.AuthConfig = &auth.Config{
		IssuerURL:      env.DexWebEndpoint,
		OfflineAsScope: true,
		RootToken:      &rootToken,
		RedirectURI:    "http://localhost:6555/oauth/callback",
		ClientId:       "gateway",
		ClientSecret:   "app-secret",
		Nonce:          "secret-nonce",
	}

	return env, nil
}

func (env *OAuthTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}

func (env *TestEnv) Shutdown() error {
	log := env.log
	log.Info("Tearing down testenv...")

	for key, element := range env.resources {
		log.Info("Tearing down docker resource", "resource", key)
		if err := env.pool.Purge(element); err != nil {
			return err
		}
	}

	return nil
}
