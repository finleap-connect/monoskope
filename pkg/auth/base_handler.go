package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/coreos/go-oidc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type BaseHandler struct {
	httpClient *http.Client
	config     *BaseConfig
	Provider   *oidc.Provider
	log        logger.Logger
}

func NewBaseHandler(config *BaseConfig) (*BaseHandler, error) {
	n := &BaseHandler{
		config:     config,
		httpClient: http.DefaultClient,
		log:        logger.WithName("auth-base"),
	}
	// Setup OIDC
	err := n.setupOIDC()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *BaseHandler) setupOIDC() error {
	ctx := oidc.ClientContext(context.Background(), n.httpClient)

	// Using an exponantial backoff to avoid issues in development environments
	backoffParams := backoff.NewExponentialBackOff()
	backoffParams.MaxElapsedTime = time.Second * 10
	err := backoff.Retry(func() error {
		var err error
		n.Provider, err = oidc.NewProvider(ctx, n.config.IssuerURL)
		return err
	}, backoffParams)
	if err != nil {
		return fmt.Errorf("failed to query provider %q: %v", n.config.IssuerURL, err)
	}

	// What scopes does a provider support?
	var scopes struct {
		// See: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
		Supported []string `json:"scopes_supported"`
	}
	if err := n.Provider.Claims(&scopes); err != nil {
		return fmt.Errorf("failed to parse provider scopes_supported: %v", err)
	}
	if len(scopes.Supported) == 0 {
		// scopes_supported is a "RECOMMENDED" discovery claim, not a required
		// one. If missing, assume that the provider follows the spec and has
		// an "offline_access" scope.
		n.config.OfflineAsScope = true
	} else {
		// See if scopes_supported has the "offline_access" scope.
		n.config.OfflineAsScope = func() bool {
			for _, scope := range scopes.Supported {
				if scope == oidc.ScopeOfflineAccess {
					return true
				}
			}
			return false
		}()
	}

	n.log.Info("connected to auth provider", "AuthURL", n.Provider.Endpoint().AuthURL, "TokenURL", n.Provider.Endpoint().TokenURL, "claims", scopes.Supported)

	return nil
}
