package auth

import (
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
)

func NewAuthStatusCmd(configLoader *config.ClientConfigLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  `Shows if authenticated against any Monoskope instance and against which one.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configLoader.LoadAndStoreConfig(); err != nil {
				return err
			}
			return nil
		},
	}
}
