package test

import (
	"fmt"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	gw_auth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	monoctl_auth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/auth"
)

const (
	AuthRootToken       = "super-secret-root-token"
	RedirectURLHostname = "localhost"
	RedirectURLPort     = ":8000"
)

type TestEnv struct {
	pool      *dockertest.Pool
	resources map[string]*dockertest.Resource
	Log       logger.Logger
}

func (t *TestEnv) CreateDockerPool() error {
	t.Log.Info("Creating docker pool...")

	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}
	t.pool = pool
	return nil
}

func (t *TestEnv) RunWithOptions(opts *dockertest.RunOptions) (*dockertest.Resource, error) {
	t.Log.Info(fmt.Sprintf("Starting docker container resource %s/%s...", opts.Repository, opts.Name))
	res, err := t.pool.RunWithOptions(opts)
	if err != nil {
		return nil, err
	}
	t.resources[res.Container.Name] = res
	return res, err
}

type OAuthTestEnv struct {
	*TestEnv
	DexWebEndpoint string
	AuthConfig     *gw_auth.Config
}

func SetupGeneralTestEnv(envName string) *TestEnv {
	log := logger.WithName(envName)
	env := &TestEnv{
		Log:       log,
		resources: make(map[string]*dockertest.Resource),
	}
	log.Info("Setting up testenv...")
	return env
}

func SetupAuthTestEnv(envName string) (*OAuthTestEnv, error) {
	env := &OAuthTestEnv{
		TestEnv: SetupGeneralTestEnv(envName),
	}

	err := env.CreateDockerPool()
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}

	dexContainer, err := env.RunWithOptions(&dockertest.RunOptions{
		Name:       "dex",
		Repository: "quay.io/dexidp/dex",
		Tag:        "v2.25.0",
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5556": {{HostPort: "5556"}},
		},
		ExposedPorts: []string{"5556", "5000"},
		Cmd:          []string{"serve", "/etc/dex/cfg/config.yaml"},
		Mounts:       []string{fmt.Sprintf("%s:/etc/dex/cfg", DexConfigPath)},
	})
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.DexWebEndpoint = fmt.Sprintf("http://127.0.0.1:%s", dexContainer.GetPort("5556/tcp"))

	rootToken := AuthRootToken
	env.AuthConfig = &gw_auth.Config{
		IssuerURL:      env.DexWebEndpoint,
		OfflineAsScope: true,
		RootToken:      &rootToken,
		ClientId:       "gateway",
		ClientSecret:   "app-secret",
		Nonce:          "secret-nonce",
	}

	return env, nil
}

func (env *OAuthTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}

func (env *OAuthTestEnv) NewOidcClientServer(ready chan<- string) (*monoctl_auth.Server, error) {
	serverConf := &monoctl_auth.Config{
		LocalServerBindAddress: []string{
			fmt.Sprintf("%s%s", RedirectURLHostname, RedirectURLPort),
		},
		RedirectURLHostname:  RedirectURLHostname,
		LocalServerReadyChan: ready,
	}
	server, err := monoctl_auth.NewServer(serverConf)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (env *TestEnv) Shutdown() error {
	log := env.Log
	log.Info("Tearing down testenv...")

	if env.resources != nil {
		for key, element := range env.resources {
			log.Info("Tearing down docker resource", "resource", key)
			if err := env.pool.Purge(element); err != nil {
				return err
			}
		}
	}

	return nil
}
