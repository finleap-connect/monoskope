package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/util"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

func NewViewCmd(configLoader *config.ClientConfigManager) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "view",
		Short: "View monoctl config",
		Long:  `View monoctl configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := util.LoadConfig(configLoader); err != nil {
				return err
			}

			fmt.Printf("%s:\n", configLoader.GetConfigLocation())

			conf, err := configLoader.GetConfig().String()
			if err != nil {
				return err
			}
			fmt.Println(conf)

			return nil
		},
	}

	return initCmd
}
