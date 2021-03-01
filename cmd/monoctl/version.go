package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/usecases"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/util"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewVersionCmd(cmdName string, configManager *config.ClientConfigManager) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		RunE: func(cmd *cobra.Command, args []string) error {
			version.PrintVersion(cmdName)

			if err := checkVersion(cmd.Context(), configManager, false); err != nil {
				if status, ok := status.FromError(errors.Unwrap(err)); ok && status.Code() == codes.Unauthenticated {
					return checkVersion(cmd.Context(), configManager, true)
				}
				return err
			}
			return nil
		},
	}
}

func checkVersion(ctx context.Context, configManager *config.ClientConfigManager, forceAuth bool) error {
	if err := util.LoadConfigAndAuth(ctx, configManager, util.Timeout, forceAuth); err != nil {
		return fmt.Errorf("init failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, util.Timeout)
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
}
