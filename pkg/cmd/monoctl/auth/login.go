package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	usecases "gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/usecases"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
)

var (
	timeout time.Duration
)

func NewAuthLoginCmd(configLoader *config.ClientConfigLoader) *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Start authentication flow",
		Long:  `Starts the authentication flow against Monoskope.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configLoader.LoadAndStoreConfig(); err != nil {
				return fmt.Errorf("Failed loading monoconfig: %w", err)
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
			defer cancel()

			err := usecases.NewAuthUsecase(ctx, configLoader).Run()
			if err != nil {
				return fmt.Errorf("Failed to authenticate: %w", err)
			} else {
				fmt.Printf("Successfully authenticated as %s!", configLoader.GetConfig().AuthInformation.Subject)
			}
			return nil
		},
	}

	flags := loginCmd.Flags()
	flags.DurationVar(&timeout, "timeout", 120*time.Second, "Timeout for the auth process, defaults to 60s")

	return loginCmd
}
