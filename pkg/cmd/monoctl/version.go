package monoctl

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	usecases "gitlab.figo.systems/platform/monoskope/monoskope/internal/usecases/monoctl"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/cmd/monoctl/flags"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/cmd/util"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
)

func NewVersionCmd(cmdName string, configManager *config.ClientConfigManager) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		RunE: func(cmd *cobra.Command, args []string) error {
			util.PrintVersion(cmdName)

			if err := util.LoadConfigAndAuth(cmd.Context(), configManager, flags.Timoeut); err != nil {
				return fmt.Errorf("init failed: %w", err)
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), flags.Timoeut)
			defer cancel()

			result, err := usecases.NewServerVersionUseCase(ctx, configManager.GetConfig()).Run()
			if err != nil {
				return fmt.Errorf("failed to retrieve server version: %w", err)
			}
			fmt.Print(result)
			fmt.Println()

			return nil
		},
	}
}
