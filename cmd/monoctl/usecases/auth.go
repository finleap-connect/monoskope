package usecases

import (
	"context"
	"fmt"

	"github.com/pkg/browser"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	monoctlAuth "gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"golang.org/x/sync/errgroup"
)

// AuthUseCase provides the internal use-case of authentication.
type AuthUseCase struct {
	log          logger.Logger
	config       *config.Config
	configLoader *config.ClientConfigManager
}

func NewAuthUsecase(configLoader *config.ClientConfigManager, force bool) *AuthUseCase {
	useCase := &AuthUseCase{
		log:          logger.WithName("auth-use-case"),
		config:       configLoader.GetConfig(),
		configLoader: configLoader,
	}
	return useCase
}

func (a *AuthUseCase) RunAuthenticationFlow(ctx context.Context) error {
	a.log.Info("starting authentication")

	conn, err := gateway.CreateGatewayConnecton(ctx, a.config.Server)
	if err != nil {
		return err
	}
	defer conn.Close()
	gatewayClient := api.NewGatewayClient(conn)

	ready := make(chan string, 1)
	defer close(ready)

	callbackServer, err := monoctlAuth.NewServer(&monoctlAuth.Config{
		LocalServerBindAddress: []string{
			"localhost:8000",
			"localhost:18000",
		},
		RedirectURLHostname:    "localhost",
		LocalServerSuccessHTML: DefaultLocalServerSuccessHTML,
		LocalServerReadyChan:   ready,
	})
	if err != nil {
		return err
	}
	defer callbackServer.Close()

	authState := &api.AuthState{CallbackURL: callbackServer.RedirectURI}
	authInfo, err := gatewayClient.GetAuthInformation(ctx, authState)
	if err != nil {
		return err
	}

	var authCode string
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			a.log.Info("Open " + url)
			if err := browser.OpenURL(url); err != nil {
				a.log.Error(err, "could not open the browser")
				return err
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		var innerErr error
		authCode, innerErr = callbackServer.ReceiveCodeViaLocalServer(ctx, authInfo.AuthCodeURL, authInfo.State)
		return innerErr
	})
	if err := eg.Wait(); err != nil {
		a.log.Error(err, "authorization error: %s")
		return err
	}

	authResponse, err := gatewayClient.ExchangeAuthCode(ctx, &api.AuthCode{Code: authCode, State: authInfo.State, CallbackURL: callbackServer.RedirectURI})
	if err != nil {
		return err
	}

	accessToken := authResponse.GetAccessToken()
	a.config.AuthInformation = &config.AuthInformation{
		Token:        accessToken.GetToken(),
		Expiry:       accessToken.GetExpiry().AsTime(),
		RefreshToken: authResponse.GetRefreshToken(),
		Subject:      authResponse.GetEmail(),
	}
	return a.configLoader.SaveConfig()
}

func (a *AuthUseCase) RunRefreshFlow(ctx context.Context) error {
	a.log.Info("refreshing the token")
	conn, err := gateway.CreateGatewayConnecton(ctx, a.config.Server)
	if err != nil {
		return err
	}
	defer conn.Close()
	gwc := api.NewGatewayClient(conn)

	accessToken, err := gwc.RefreshAuth(ctx, &api.RefreshAuthRequest{RefreshToken: a.config.AuthInformation.RefreshToken})
	if err != nil {
		return err
	}

	a.config.AuthInformation.Token = accessToken.GetToken()
	a.config.AuthInformation.Expiry = accessToken.GetExpiry().AsTime()
	return a.configLoader.SaveConfig()
}

func (a *AuthUseCase) Run(ctx context.Context) error {
	// Check if already authenticated
	if a.config.HasAuthInformation() {
		a.log.Info("checking expiration of existing token")
		authInfo := a.config.AuthInformation
		if authInfo.IsValid() {
			a.log.Info("you have a valid auth token", "expiry", authInfo.Expiry)
			return nil
		}
		a.log.Info("your auth token has expired", "expiry", authInfo.Expiry)

		if authInfo.HasRefreshToken() {
			err := a.RunRefreshFlow(ctx)
			if err == nil {
				return nil
			}
			a.log.Error(err, "Failed to do refresh flow")
		}
	}
	return a.RunAuthenticationFlow(ctx)
}

// DefaultLocalServerSuccessHTML is a default response body on authorization success.
const DefaultLocalServerSuccessHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Authorized</title>
	<script>
		window.close()
	</script>
	<style>
		body {
			background-color: #eee;
			margin: 0;
			padding: 0;
			font-family: sans-serif;
		}
		.placeholder {
			margin: 2em;
			padding: 2em;
			background-color: #fff;
			border-radius: 1em;
		}
	</style>
</head>
<body>
	<div class="placeholder">
		<h1>Authorized</h1>
		<p>You can close this window.</p>
	</div>
</body>
</html>
`
