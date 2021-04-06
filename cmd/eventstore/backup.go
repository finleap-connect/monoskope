package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	_ "go.uber.org/automaxprocs"
)

var (
	pushGatewayUrl string
	retention      int
	timeout        string
)

var backupCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("server-cmd")

		timeout, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}

		// init event store
		log.Info("Setting up event store...")
		store, err := eventstore.NewEventStore()
		if err != nil {
			log.Error(err, "Failed to configure event store.")
			return err
		}
		defer store.Close()

		// setup backup management
		backupManger, err := eventstore.NewBackupManager(store, retention)
		if err != nil {
			log.Error(err, "Failed to configure automatic backups.")
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		result, err := backupManger.RunBackup(ctx)
		log.Info("Backing up eventstore has been successful.", "BackupIdentifier", result.BackupIdentifier, "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)

		return err
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Local flags
	flags := backupCmd.Flags()
	flags.IntVar(&retention, "retention", 7, "Count of backups to keep, <1 means keep all")
	flags.StringVar(&timeout, "timeout", "1h", "Timeout to cancel backup job")
	flags.StringVar(&pushGatewayUrl, "prometheus-gateway-url", "", "Url of the gateway to push prometheus metrics to")
}
