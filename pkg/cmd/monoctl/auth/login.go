package auth

import (
	"github.com/spf13/cobra"
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
			return nil
		},
	}
}
