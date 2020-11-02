package client

import (
	"context"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"golang.org/x/oauth2"
)

type Handler struct {
	auth.BaseHandler
	log        logger.Logger
	httpClient *http.Client
	config     *Config
}

func NewHandler(config *Config) (*Handler, error) {
	baseHandler, err := auth.NewBaseHandler(&config.BaseConfig)
	if err != nil {
		return nil, err
	}

	n := &Handler{
		BaseHandler: *baseHandler,
		config:      config,
		log:         logger.WithName("auth-client"),
	}
	n.log.Info("oidc client setup", "id", n.config.ClientId, "redirectURI", n.config.RedirectURI)

	return n, nil
}

func (n *Handler) getOauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     n.config.ClientId,
		ClientSecret: n.config.ClientSecret,
		Endpoint:     n.Provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  n.config.RedirectURI,
	}
}

func (n *Handler) clientContext(ctx context.Context) context.Context {
	return oidc.ClientContext(ctx, n.httpClient)
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page that asks for permissions for the required scopes explicitly.
func (n *Handler) GetAuthCodeURL(state *auth.State, config *auth.AuthCodeURLConfig) (string, error) {
	// Encode state and calculate nonce
	encoded, err := state.Encode()
	if err != nil {
		return "", err
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

	return authCodeURL, nil
}

func (n *Handler) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	t := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(-time.Hour),
	}
	return n.getOauth2Config(nil).TokenSource(n.clientContext(ctx), t).Token()
}

// Exchange converts an authorization code into a token.
func (n *Handler) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return n.getOauth2Config(nil).Exchange(n.clientContext(ctx), code)
}
