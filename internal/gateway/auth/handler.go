package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	api_gwauth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"golang.org/x/oauth2"
)

type Handler struct {
	config     *Config
	httpClient *http.Client
	Provider   *oidc.Provider
	verifier   *oidc.IDTokenVerifier
	log        logger.Logger
}

func NewHandler(config *Config) (*Handler, error) {
	n := &Handler{
		config:     config,
		httpClient: http.DefaultClient,
		log:        logger.WithName("auth"),
	}
	// Setup OIDC
	err := n.setupOIDC()
	if err != nil {
		return nil, err
	}
	n.verifier = n.Provider.Verifier(&oidc.Config{ClientID: n.config.ClientId})

	return n, nil
}

func (n *Handler) setupOIDC() error {
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

func (n *Handler) getOauth2Config(scopes []string, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     n.config.ClientId,
		ClientSecret: n.config.ClientSecret,
		Endpoint:     n.Provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  redirectURL,
	}
}

func (n *Handler) clientContext(ctx context.Context) context.Context {
	return oidc.ClientContext(ctx, n.httpClient)
}

func (n *Handler) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	// Generate a new token with a refresht token and the expiry of the access token set to golang zero date.
	// Setting the access token expired will force the token source to automatically use the refresh token to issue a new token.
	t := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Time{}, // golang zero date
	}
	return n.getOauth2Config(nil, "").TokenSource(n.clientContext(ctx), t).Token()
}

// Exchange converts an authorization code into a token.
func (n *Handler) Exchange(ctx context.Context, code, redirectURL string) (*oauth2.Token, error) {
	return n.getOauth2Config(nil, redirectURL).Exchange(n.clientContext(ctx), code)
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page that asks for permissions for the required scopes explicitly.
func (n *Handler) GetAuthCodeURL(state *api_gwauth.AuthState, config *AuthCodeURLConfig) (string, string, error) {
	// Encode state and calculate nonce
	encoded, err := (&State{Callback: state.GetCallbackURL()}).Encode()
	if err != nil {
		return "", "", err
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
		authCodeURL = n.getOauth2Config(scopes, state.GetCallbackURL()).AuthCodeURL(encoded, oidc.Nonce(nonce))
	} else if n.config.OfflineAsScope {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
		authCodeURL = n.getOauth2Config(scopes, state.GetCallbackURL()).AuthCodeURL(encoded, oidc.Nonce(nonce))
	} else {
		authCodeURL = n.getOauth2Config(scopes, state.GetCallbackURL()).AuthCodeURL(encoded, oidc.Nonce(nonce), oauth2.AccessTypeOffline)
	}

	return authCodeURL, encoded, nil
}

func (n *Handler) VerifyStateAndClaims(ctx context.Context, token *oauth2.Token, encodedState string) (*ExtraClaims, error) {
	if !token.Valid() {
		return nil, fmt.Errorf("failed to verify ID token")
	}

	rawIDToken := token.Extra("id_token").(string)
	idToken, err := n.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}

	if idToken.Nonce != util.HashString(encodedState+n.config.Nonce) {
		return nil, fmt.Errorf("invalid id_token nonce")
	}

	state, err := DecodeState(encodedState)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state")
	}

	if !state.IsValid() {
		return nil, grpc.ErrInvalidArgument(errors.Errorf("callback url invalid"))
	}

	claims, err := getClaims(idToken)
	if err != nil {
		return nil, err
	}

	n.log.Info("verified bearer token", "user", claims.Email)

	return claims, nil
}

// authorize verifies a bearer token and pulls user information form the claims.
func (n *Handler) Authorize(ctx context.Context, bearerToken string) (*ExtraClaims, error) {
	idToken, err := n.verifier.Verify(ctx, bearerToken)
	if err != nil {
		return nil, err
	}

	claims, err := getClaims(idToken)
	if err != nil {
		return nil, err
	}

	n.log.Info("user authenticated via bearer token", "user", claims.Email)

	return claims, nil
}

func getClaims(idToken *oidc.IDToken) (*ExtraClaims, error) {
	claims := &ExtraClaims{}
	if err := idToken.Claims(claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}
	if !claims.EmailVerified {
		return nil, fmt.Errorf("email (%q) in returned claims was not verified", claims.Email)
	}

	return claims, nil
}
