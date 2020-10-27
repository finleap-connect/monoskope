package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/coreos/go-oidc"
	dexpb "github.com/dexidp/dex/api"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type Handler struct {
	log         logger.Logger
	dexClient   dexpb.DexClient
	httpClient  *http.Client
	oauthClient *dexpb.Client
	verifier    *oidc.IDTokenVerifier
	provider    *oidc.Provider
	config      *Config
}

func NewHandler(dexClient dexpb.DexClient, config *Config) (*Handler, error) {
	n := &Handler{
		log:        logger.WithName("auth"),
		dexClient:  dexClient,
		config:     config,
		httpClient: http.DefaultClient,
	}
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

	n.verifier = n.provider.Verifier(&oidc.Config{ClientID: n.oauthClient.GetId()})
	return nil
}

func (n *Handler) Verify(ctx context.Context, bearerToken string) (*oidc.IDToken, error) {
	idToken, err := n.verifier.Verify(ctx, bearerToken)
	if err != nil {
		return nil, fmt.Errorf("could not verify bearer token: %v", err)
	}
	return idToken, nil
}

// authorize verifies a bearer token and pulls user information form the claims.
func (n *Handler) Authorize(ctx context.Context, bearerToken string) (*ExtraClaims, error) {
	idToken, err := n.Verify(ctx, bearerToken)
	if err != nil {
		return nil, err
	}

	claims := &ExtraClaims{}
	if err = idToken.Claims(claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}
	if !claims.EmailVerified {
		return nil, fmt.Errorf("email (%q) in returned claims was not verified", claims.Email)
	}

	return claims, nil
}
