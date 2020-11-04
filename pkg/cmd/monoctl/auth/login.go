package auth

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	usecases "gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/usecases"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
)

func NewAuthLoginCmd(configLoader *config.ClientConfigLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Start authentication flow",
		Long:  `Starts the authentication flow against Monoskope.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configLoader.LoadAndStoreConfig(); err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
			defer cancel()

			err := usecases.NewAuthUsecase(ctx, configLoader.GetConfig()).Run()

			return err
		},
	}
}
