package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

func NewAuthStatusCmd(configLoader *config.ClientConfigManager) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  `Shows if authenticated against any Monoskope instance and against which one.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configLoader.LoadAndStoreConfig(); err != nil {
				return fmt.Errorf("failed loading monoconfig: %w", err)
			}

			conf := configLoader.GetConfig()
			authenticated := conf.HasAuthInformation() && conf.AuthInformation.HasToken()

			fmt.Printf("Authenticated: %v\n", authenticated)
			if authenticated {
				fmt.Printf("Server: %v\n", conf.Server)
				fmt.Printf("User: %v\n", conf.AuthInformation.Subject)
				fmt.Printf("Token expiry: %v\n", conf.AuthInformation.Expiry)
				fmt.Printf("Token expired: %v\n", conf.AuthInformation.IsTokenExpired())
				fmt.Printf("Refresh possible: %v\n", conf.AuthInformation.HasRefreshToken())
			}

			return nil
		},
	}
}
