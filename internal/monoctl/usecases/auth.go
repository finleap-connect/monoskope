package usecases

import (
	"context"
	"fmt"

	"github.com/pkg/browser"
	api_gw "gitlab.figo.systems/platform/monoskope/monoskope/api/gateway"
	gw_auth "gitlab.figo.systems/platform/monoskope/monoskope/api/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	monoctl_auth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
	"golang.org/x/sync/errgroup"
)

type AuthUseCase struct {
	log          logger.Logger
	ctx          context.Context
	config       *config.Config
	configLoader *config.ClientConfigLoader
}

func NewAuthUsecase(ctx context.Context, configLoader *config.ClientConfigLoader) *AuthUseCase {
	useCase := &AuthUseCase{
		log:          logger.WithName("auth-use-case"),
		config:       configLoader.GetConfig(),
		configLoader: configLoader,
		ctx:          ctx,
	}
	return useCase
}

func (a *AuthUseCase) Run() error {
	var err error

	conn, err := gateway.CreateGatewayConnecton(a.ctx, a.config.Server, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	gwc := api_gw.NewGatewayClient(conn)

	ready := make(chan string, 1)
	defer close(ready)

	serverConf := &monoctl_auth.Config{
		LocalServerBindAddress: []string{
			"localhost:8000",
			"localhost:18000",
		},
		RedirectURLHostname:    "localhost",
		LocalServerSuccessHTML: DefaultLocalServerSuccessHTML,
		LocalServerReadyChan:   ready,
	}
	server, err := monoctl_auth.NewServer(serverConf)
	if err != nil {
		return err
	}
	defer server.Close()

	eg, ctx := errgroup.WithContext(a.ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			a.log.Info("Open " + url)
			if err := browser.OpenURL(url); err != nil {
				a.log.Error(err, "could not open the browser")
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	if err := eg.Wait(); err != nil {
		a.log.Error(err, "authorization error: %s")
		return err
	}

	authState := &gw_auth.AuthState{CallbackURL: server.RedirectURI}
	authInfo, err := gwc.GetAuthInformation(a.ctx, authState)
	if err != nil {
		return err
	}

	authCode, err := server.ReceiveCodeViaLocalServer(a.ctx, authInfo.AuthCodeURL, authInfo.State)
	if err != nil {
		return err
	}

	//TODO implement what to do with the token stuff now
	userInfo, err := gwc.ExchangeAuthCode(a.ctx, &gw_auth.AuthCode{Code: authCode, State: authInfo.State})
	if err != nil {
		return err
	}

	a.config.AuthInformation = &config.AuthInformation{
		Token:        userInfo.GetAccessToken(),
		RefreshToken: userInfo.GetRefreshToken(),
		Expiry:       userInfo.GetExpiry().AsTime(),
		Subject:      userInfo.GetEmail(),
	}
	err = a.configLoader.SaveConfig()

	return nil
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
