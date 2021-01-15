package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/usecases"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/util"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
)

func NewVersionCmd(cmdName string, configManager *config.ClientConfigManager) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		RunE: func(cmd *cobra.Command, args []string) error {
			version.PrintVersion(cmdName)

			if err := util.LoadConfigAndAuth(cmd.Context(), configManager, util.Timeout); err != nil {
				return fmt.Errorf("init failed: %w", err)
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), util.Timeout)
			defer cancel()

			result, err := usecases.NewServerVersionUseCase(ctx, configManager.GetConfig()).Run()
			if err != nil {
				return fmt.Errorf("failed to retrieve server version: %w", err)
			}

			for _, version := range result {
				fmt.Print(version)
				fmt.Println()
			}

			return nil
		},
	}
}
