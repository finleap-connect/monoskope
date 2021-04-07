package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/backup"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

var (
	pushGatewayUrl string
	retention      int
	timeout        string
)

var backupCmd = &cobra.Command{
	Use:   "backup [flags]",
	Short: "Starts the backup",
	Long:  `Starts the backup of the event store`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("backup-cmd")

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

		// setup metrics publisher
		log.Info("Setting up metrics publisher...")
		var metricsPublisher backup.MetricsPublisher
		if len(pushGatewayUrl) > 0 {
			metricsPublisher, err = backup.NewMetricsPublisher(pushGatewayUrl)
			if err != nil {
				log.Error(err, "Failed to configure metrics publisher.")
				return err
			}
		} else {
			metricsPublisher = backup.NewNoopMetricsPublisher()
		}
		defer util.PanicOnError(metricsPublisher.CloseAndPush())

		// setup backup management
		log.Info("Setting up backup manager...")
		backupManger, err := eventstore.NewBackupManager(metricsPublisher, store, retention)
		if err != nil {
			log.Error(err, "Failed to configure backup.")
			return err
		}

		// setup context
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// start backup
		log.Info("Starting backup...")
		err = runBackup(ctx, log, metricsPublisher, backupManger)
		if err != nil {
			log.Error(err, "Failed to run backup.")
			return err
		}

		// start purge
		log.Info("Starting purge...")
		err = runPurge(ctx, log, backupManger)
		if err != nil {
			log.Error(err, "Failed to run purge.")
			return err
		}

		return err
	},
}

func runBackup(ctx context.Context, log logger.Logger, metricsPublisher backup.MetricsPublisher, backupManger *eventstore.BackupManager) error {
	metricsPublisher.Start()
	result, err := backupManger.RunBackup(ctx)
	metricsPublisher.Finished()
	metricsPublisher.SetBytes(float64(result.ProcessedBytes))
	metricsPublisher.SetEventCount(float64(result.ProcessedEvents))

	if err != nil {
		metricsPublisher.SetFailTime()
		log.Error(err, "Failed to back up eventstore.", "BackupIdentifier", result.BackupIdentifier, "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)
	} else {
		metricsPublisher.SetSuccessTime()
		log.Info("Backing up eventstore has been successful.", "BackupIdentifier", result.BackupIdentifier, "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)
	}
	return err
}

func runPurge(ctx context.Context, log logger.Logger, backupManger *eventstore.BackupManager) error {
	return nil
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Local flags
	flags := backupCmd.Flags()
	flags.IntVar(&retention, "retention", 7, "Count of backups to keep, <1 means keep all")
	flags.StringVar(&timeout, "timeout", "1h", "Timeout after which to cancel the backup job")
	flags.StringVar(&pushGatewayUrl, "prometheus-gateway-url", "", "Url of the gateway to push prometheus metrics to")
}
