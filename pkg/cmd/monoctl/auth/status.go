package auth

import (
	"fmt"

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
				return fmt.Errorf("Failed loading monoconfig: %w", err)
			}

			conf := configLoader.GetConfig()
			if conf.HasToken() {
				fmt.Printf("Authenticated against '%v'\n", conf.Server)
				fmt.Printf("Auth token valid until %v\n", conf.AuthInformation.Expiry)
			} else {
				fmt.Printf("Not authenticated\n")
			}
			return nil
		},
	}
}
