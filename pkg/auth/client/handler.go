package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/coreos/go-oidc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"golang.org/x/oauth2"
)

type Handler struct {
	log        logger.Logger
	httpClient *http.Client
	verifier   *oidc.IDTokenVerifier
	provider   *oidc.Provider
	config     *Config
}

func NewHandler(config *Config) (*Handler, error) {
	n := &Handler{
		log:        logger.WithName("auth-client"),
		config:     config,
		httpClient: http.DefaultClient,
	}
	// Setup the redirect handler before continuing
	n.log.Info("OAuth-Client setup", "id", n.config.ClientId, "redirectURI", n.config.RedirectURI)
	// Setup OIDC
	err := n.setupOIDC()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Handler) setupOIDC() error {
	ctx := oidc.ClientContext(context.Background(), n.httpClient)

	// Using an exponantial backoff to avoid issues in development environments
	backoffParams := backoff.NewExponentialBackOff()
	backoffParams.MaxElapsedTime = time.Second * 10
	err := backoff.Retry(func() error {
		var err error
		n.provider, err = oidc.NewProvider(ctx, n.config.IssuerURL)
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
	if err := n.provider.Claims(&scopes); err != nil {
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

	n.verifier = n.provider.Verifier(&oidc.Config{ClientID: n.config.ClientId})
	return nil
}

func (n *Handler) getOauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     n.config.ClientId,
		ClientSecret: n.config.ClientSecret,
		Endpoint:     n.provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  n.config.RedirectURI,
	}
}

func (n *Handler) clientContext(ctx context.Context) context.Context {
	return oidc.ClientContext(ctx, n.httpClient)
}

func (n *Handler) GetAuthCodeURL(state *auth.State, config *auth.AuthCodeURLConfig) (string, error) {
	// Encode state and calculate nonce
	encoded, err := state.Encode()
	if err != nil {
		return "", err // TODO: wrap?
	}
	nonce := util.HashString(encoded + n.config.Nonce)

	scopes := config.Scopes
	for _, client := range config.Clients {
		scopes = append(scopes, "audience:server:client_id:"+client)
	}
	scopes = append(scopes, oidc.ScopeOpenID, "profile", "email")

	// Construct authCodeURL
	authCodeURL := ""
	if config.OfflineAccess {
		authCodeURL = n.getOauth2Config(scopes).AuthCodeURL(encoded, oidc.Nonce(nonce))
	} else if n.config.OfflineAsScope {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
		authCodeURL = n.getOauth2Config(scopes).AuthCodeURL(encoded, oidc.Nonce(nonce))
	} else {
		authCodeURL = n.getOauth2Config(scopes).AuthCodeURL(encoded, oidc.Nonce(nonce), oauth2.AccessTypeOffline)
	}
	authCodeURL = authCodeURL + "&connector_id=" + state.ConnectorID
	return authCodeURL, nil
}

func (n *Handler) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	t := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(-time.Hour),
	}
	return n.getOauth2Config(nil).TokenSource(n.clientContext(ctx), t).Token()
}

func (n *Handler) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return n.getOauth2Config(nil).Exchange(n.clientContext(ctx), code)
}
