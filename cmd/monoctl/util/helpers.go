package util

import (
	"context"
	"fmt"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/usecases"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

func LoadConfig(configManager *config.ClientConfigManager) error {
	if err := configManager.LoadAndStoreConfig(); err != nil {
		return fmt.Errorf("filed loading monoconfig: %w", err)
	}
	return nil
}

func LoadConfigAndAuth(ctx context.Context, configManager *config.ClientConfigManager, timeout time.Duration, force bool) error {
	if err := LoadConfig(configManager); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return usecases.NewAuthUsecase(configManager, force).Run(ctx)
}
