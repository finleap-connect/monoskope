package auth

import (
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

func NewAuthCmd(configLoader *config.ClientConfigManager) *cobra.Command {
	authCmd := &cobra.Command{
		Use:                   "auth",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Short:                 "Handle authorization",
		Long:                  `Authenticate with remote Monoskope instance, check status and more.`,
	}

	authCmd.AddCommand(NewAuthStatusCmd(configLoader))
	authCmd.AddCommand(NewAuthLoginCmd(configLoader))

	return authCmd
}
