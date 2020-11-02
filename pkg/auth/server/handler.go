package server

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type Handler struct {
	auth.BaseHandler
	config   *Config
	verifier *oidc.IDTokenVerifier
	log      logger.Logger
}

func NewHandler(config *Config) (*Handler, error) {
	baseHandler, err := auth.NewBaseHandler(&config.BaseConfig)
	if err != nil {
		return nil, err
	}

	n := &Handler{
		BaseHandler: *baseHandler,
		config:      config,
		log:         logger.WithName("auth-server"),
	}
	n.setupVerifier()

	return n, nil
}

func (n *Handler) setupVerifier() {
	n.verifier = n.Provider.Verifier(&oidc.Config{ClientID: n.config.ValidClientId})
}

func (n *Handler) verify(ctx context.Context, bearerToken string) (*oidc.IDToken, error) {
	idToken, err := n.verifier.Verify(ctx, bearerToken)
	if err != nil {
		return nil, fmt.Errorf("could not verify bearer token: %v", err)
	}
	return idToken, nil
}

// authorize verifies a bearer token and pulls user information form the claims.
func (n *Handler) Authorize(ctx context.Context, bearerToken string) (*auth.ExtraClaims, error) {
	if n.config.RootToken != nil && bearerToken == *n.config.RootToken {
		n.log.Info("### user authenticated via root token")
		return &auth.ExtraClaims{EmailVerified: true, Email: "root@monoskope"}, nil
	}

	idToken, err := n.verify(ctx, bearerToken)
	if err != nil {
		return nil, err
	}

	claims := &auth.ExtraClaims{}
	if err = idToken.Claims(claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}
	if !claims.EmailVerified {
		return nil, fmt.Errorf("email (%q) in returned claims was not verified", claims.Email)
	}

	n.log.Info("user authenticated via bearer token", "user", claims.Email)

	return claims, nil
}
