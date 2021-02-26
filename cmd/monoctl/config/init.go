package config

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

var (
	serverURL string
	force     bool
)

func NewInitCmd(configLoader *config.ClientConfigManager) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init monoctl config",
		Long:  `Init monoctl and create a new monoskope configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURL == "" {
				return errors.New("failed initializing monoconfig: server-url is required")
			}

			config := config.NewConfig()
			config.Server = serverURL

			if err := configLoader.InitConifg(config, force); err != nil {
				return fmt.Errorf("failed initializing monoconfig: %w", err)
			}

			return nil
		},
	}

	flags := initCmd.Flags()
	flags.StringVarP(&serverURL, "server-url", "u", "", "URL of the monoskope instance")
	flags.BoolVarP(&force, "force", "f", false, "Force overwrite configuraton.")

	err := initCmd.MarkFlagRequired("server-url")
	if err != nil {
		panic(err)
	}

	return initCmd
}
