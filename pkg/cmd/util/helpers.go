package util

import (
	"context"
	"fmt"
	"time"

	usecases "gitlab.figo.systems/platform/monoskope/monoskope/internal/usecases/monoctl"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
)

func LoadConfig(configManager *config.ClientConfigManager) error {
	if err := configManager.LoadAndStoreConfig(); err != nil {
		return fmt.Errorf("Failed loading monoconfig: %w", err)
	}
	return nil
}

func LoadConfigAndAuth(ctx context.Context, configManager *config.ClientConfigManager, timeout time.Duration) error {
	if err := LoadConfig(configManager); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return usecases.NewAuthUsecase(ctx, configManager).Run()
}
