package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/util"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

var (
	force bool
)

func NewAuthLoginCmd(configManager *config.ClientConfigManager) *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Start authentication flow",
		Long:  `Starts the authentication flow against Monoskope.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := util.LoadConfigAndAuth(cmd.Context(), configManager, util.Timeout, force)
			if err != nil {
				return fmt.Errorf("failed to authenticate: %w", err)
			} else {
				fmt.Printf("Successfully authenticated!\n")
			}
			return nil
		},
	}

	flags := loginCmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Force login even if already authenticated.")

	return loginCmd
}
