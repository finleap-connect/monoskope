package config

import (
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

func NewConfigCmd(configLoader *config.ClientConfigManager) *cobra.Command {
	authCmd := &cobra.Command{
		Use:                   "config",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Short:                 "Modify monoconfig files",
		Long:                  `Everything around monoctl config.`,
	}

	authCmd.AddCommand(NewInitCmd(configLoader))
	authCmd.AddCommand(NewViewCmd(configLoader))

	return authCmd
}
