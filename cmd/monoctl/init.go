package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

var (
	serverURL string
)

func NewInitCmd(configLoader *config.ClientConfigManager) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init monoctl",
		Long:  `Init monoctl and create a new monoskope configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURL == "" {
				return errors.New("failed initializing monoconfig: server-url is required")
			}

			config := config.NewConfig()
			config.Server = serverURL

			if err := configLoader.InitConifg(config); err != nil {
				return fmt.Errorf("failed initializing monoconfig: %w", err)
			}

			return nil
		},
	}
	flags := initCmd.Flags()
	flags.StringVar(&serverURL, "server-url", "", "URL of the monoskope instance")
	return initCmd
}
