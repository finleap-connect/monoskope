package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

var (
	timeoutRestore   string
	backupIdentifier string
)

var restoreCmd = &cobra.Command{
	Use:   "restore [flags]",
	Short: "Starts the restore",
	Long:  `Starts the restore of the event store from a backup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("restore-cmd")

		timeout, err := time.ParseDuration(timeoutRestore)
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
		log.Info("Setting up backup manager...")
		backupManger, err := eventstore.NewBackupManager(store, retention)
		if err != nil {
			log.Error(err, "Failed to configure backup.")
			return err
		}

		// setup context
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// start backup
		log.Info("Starting restore...")
		return runRestore(ctx, log, backupManger)
	},
}

func runRestore(ctx context.Context, log logger.Logger, backupManger *eventstore.BackupManager) error {
	result, err := backupManger.RunRestore(ctx, backupIdentifier)
	if err != nil {
		log.Error(err, "Restore failed", "ProcessedBytes", result.ProcessedBytes, "ProcessedEvents", result.ProcessedEvents)
	} else {
		log.Info("Restore finished successful", "ProcessedBytes", result.ProcessedBytes, "ProcessedEvents", result.ProcessedEvents)
	}
	return err
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	// Local flags
	flags := restoreCmd.Flags()
	flags.StringVar(&timeoutBackup, "timeout", "1h", "Timeout after which to cancel the restore job")
	flags.StringVar(&backupIdentifier, "identifier", "", "Identifier of the backup to restore")
	util.PanicOnError(restoreCmd.MarkFlagRequired("identifier"))
}
