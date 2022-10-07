// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/eventstore/backup"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/spf13/cobra"
)

var (
	pushGatewayUrl string
	retention      int
	timeoutBackup  string
)

var backupCmd = &cobra.Command{
	Use:   "backup [flags]",
	Short: "Starts the backup",
	Long:  `Starts the backup of the event store`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log := logger.WithName("backup-cmd")

		timeout, err := time.ParseDuration(timeoutBackup)
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
		defer util.PanicOnErrorFunc(metricsPublisher.CloseAndPush)

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
	result, err := backupManger.RunPurge(ctx)
	if err != nil {
		log.Error(err, "Failed to purge backups.", "BackupsLeft", result.BackupsLeft, "PurgedBackups", result.PurgedBackups)
	} else {
		log.Info("Purging outdated backups has been successful.", "BackupsLeft", result.BackupsLeft, "PurgedBackups", result.PurgedBackups)
	}
	return err
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Local flags
	flags := backupCmd.Flags()
	flags.IntVar(&retention, "retention", 7, "Count of backups to keep, <1 means keep all")
	flags.StringVar(&timeoutBackup, "timeout", "1h", "Timeout after which to cancel the backup job")
	flags.StringVar(&pushGatewayUrl, "prometheus-gateway-url", "", "Url of the gateway to push prometheus metrics to")
}
