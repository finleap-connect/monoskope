package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/coreos/go-oidc"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"
)

type Handler struct {
	config           *Config
	httpClient       *http.Client
	Provider         *oidc.Provider
	upstreamVerifier *oidc.IDTokenVerifier
	verifier         jwt.JWTVerifier
	signer           jwt.JWTSigner
	log              logger.Logger
}

func NewHandler(config *Config, signer jwt.JWTSigner, verifier jwt.JWTVerifier) *Handler {
	n := &Handler{
		config:     config,
		signer:     signer,
		verifier:   verifier,
		httpClient: http.DefaultClient,
		log:        logger.WithName("auth"),
	}
	n.log.Info("Auth handler configured.",
		"IdentityProviderName",
		n.config.IdentityProviderName,
		"IdentityProvider",
		n.config.IdentityProvider,
		"Scopes",
		n.config.Scopes,
		"RedirectURIs",
		n.config.RedirectURIs,
	)
	return n
}

func (n *Handler) SetupOIDC(ctx context.Context) error {
	ctx = oidc.ClientContext(ctx, n.httpClient)

	// Using an exponential backoff to avoid issues in development environments
	backoffParams := backoff.NewExponentialBackOff()
	backoffParams.MaxElapsedTime = time.Second * 10
	err := backoff.Retry(func() error {
		var err error
		n.Provider, err = oidc.NewProvider(ctx, n.config.IdentityProvider)
		return err
	}, backoffParams)
	if err != nil {
		return fmt.Errorf("failed to query provider %q: %v", n.config.IdentityProvider, err)
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

	n.upstreamVerifier = n.Provider.Verifier(&oidc.Config{ClientID: n.config.ClientId})

	n.log.Info("Connected to auth provider successful.", "AuthURL", n.Provider.Endpoint().AuthURL, "TokenURL", n.Provider.Endpoint().TokenURL, "AuthStyle", n.Provider.Endpoint().AuthStyle, "SupportedScopes", scopes.Supported)

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

func getClaims(idToken *oidc.IDToken) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}

	if err := idToken.Claims(claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}

	if !claims.EmailVerified {
		return nil, fmt.Errorf("email (%q) in returned claims was not verified", claims.Email)
	}

	return claims, nil
}

// exchange exchanges the auth code with a token of the upstream IDP
func (n *Handler) exchange(ctx context.Context, code, redirectURL string) (*oauth2.Token, error) {
	n.log.Info("Exchanging auth code for token...")
	return n.getOauth2Config(nil, redirectURL).Exchange(n.clientContext(ctx), code)
}

func (n *Handler) redirectUrlAllowed(callBackUrl string) bool {
	for _, validUrl := range n.config.RedirectURIs {
		if strings.EqualFold(strings.ToLower(validUrl), strings.ToLower(callBackUrl)) {
			return true
		}
	}
	return false
}

func (n *Handler) verifyStateAndClaims(ctx context.Context, token *oauth2.Token, encodedState string) (*jwt.StandardClaims, error) {
	n.log.Info("Verifying state and claims...")
	if !token.Valid() {
		return nil, fmt.Errorf("failed to verify ID token")
	}

	rawIDToken := token.Extra("id_token").(string)
	idToken, err := n.upstreamVerifier.Verify(ctx, rawIDToken)
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
		return nil, grpcUtil.ErrInvalidArgument("url is invalid")
	}

	claims, err := getClaims(idToken)
	if err != nil {
		return nil, err
	}

	n.log.Info("Token verified successfully.", "User", claims.Email, "TokenType", token.TokenType)

	return claims, nil
}

// Exchange exchanges the auth code with a token of the upstream IDP and verifies the claims
func (n *Handler) Exchange(ctx context.Context, code, state, redirectURL string) (*jwt.StandardClaims, error) {
	upstreamToken, err := n.exchange(ctx, code, redirectURL)
	if err != nil {
		return nil, err
	}
	n.log.V(logger.DebugLevel).Info("Token received in exchange for auth code.", "Token", upstreamToken)

	upstreamClaims, err := n.verifyStateAndClaims(ctx, upstreamToken, state)
	if err != nil {
		return nil, err
	}
	n.log.V(logger.DebugLevel).Info("Claims verified.", "Claims", upstreamClaims)

	return upstreamClaims, nil
}

// IssueToken wraps the upstream claims in a JWT signed by Monoskope
func (n *Handler) IssueToken(ctx context.Context, upstreamClaims *jwt.StandardClaims, userId string) (string, *jwt.AuthToken, error) {
	if upstreamClaims.FederatedClaims == nil {
		upstreamClaims.FederatedClaims = make(map[string]string)
	}
	upstreamClaims.FederatedClaims["connector_id"] = n.config.IdentityProviderName

	token := jwt.NewAuthToken(upstreamClaims, n.config.URL, userId)
	n.log.V(logger.DebugLevel).Info("Token issued successfully.", "RawToken", token, "Expiry", token.Expiry.Time().String())

	signedToken, err := n.signer.GenerateSignedToken(token)
	if err != nil {
		return "", nil, err
	}
	n.log.V(logger.DebugLevel).Info("Token signed successfully.", "SignedToken", signedToken)

	return signedToken, token, err
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page that asks for permissions for the required scopes explicitly.
func (n *Handler) GetAuthCodeURL(state *api.AuthState, scopes []string) (string, string, error) {
	if !n.redirectUrlAllowed(state.GetCallbackUrl()) {
		return "", "", errors.New("callback url not allowed")
	}

	// Encode state and calculate nonce
	encoded, err := (&State{Callback: state.GetCallbackUrl()}).Encode()
	if err != nil {
		return "", "", err
	}
	nonce := util.HashString(encoded + n.config.Nonce)

	// Construct authCodeURL
	var authCodeURL string
	if n.config.OfflineAsScope {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
		authCodeURL = n.getOauth2Config(scopes, state.GetCallbackUrl()).AuthCodeURL(encoded, oidc.Nonce(nonce))
	} else {
		authCodeURL = n.getOauth2Config(scopes, state.GetCallbackUrl()).AuthCodeURL(encoded, oidc.Nonce(nonce), oauth2.AccessTypeOffline)
	}

	return authCodeURL, encoded, nil
}

// Authorize parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (n *Handler) Authorize(ctx context.Context, token string, claims interface{}) error {
	if err := n.verifier.Verify(token, claims); err != nil {
		return err
	}
	return nil
}

func (n *Handler) Keys() *jose.JSONWebKeySet {
	return n.verifier.JWKS()
}

func (n *Handler) KeyExpiration() time.Duration {
	return n.verifier.KeyExpiration()
}
