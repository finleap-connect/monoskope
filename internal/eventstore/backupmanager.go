package eventstore

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/robfig/cron"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/backup"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/backup/s3"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const BackupPath = "/etc/eventstore/backup"

type BackupManager struct {
	log           logr.Logger
	store         eventsourcing.Store
	backupHandler backup.BackupHandler
	cron          *cron.Cron
	cronSchedule  string
	retention     int
}

// NewBackupManager creates a new backup manager configured by config files taken from eventstore.BackupPath and environment config.
func NewBackupManager(store eventsourcing.Store) (*BackupManager, error) {
	manager := &BackupManager{
		log:   logger.WithName("backup-manager"),
		store: store,
	}
	if err := manager.configure(); err != nil {
		return nil, err
	}
	return manager, nil
}

func (b *BackupManager) configure() error {
	// Get backup destination configuration
	fileInfos, err := ioutil.ReadDir(BackupPath)
	if err != nil {
		return err
	}

	// Lookup configured backup destinations
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		switch fileInfo.Name() {
		case "s3.yaml":
			conf, err := s3.NewS3ConfigFromFile(path.Join(BackupPath, fileInfo.Name()))
			if err != nil {
				return err
			}
			b.backupHandler = s3.NewS3BackupHandler(conf, b.store)
		}
	}

	b.cronSchedule = os.Getenv("BACKUP_SCHEDULE")
	if _, err := cron.Parse(b.cronSchedule); err != nil {
		return err
	}

	retention, err := strconv.Atoi(os.Getenv("BACKUP_RETENTION_COUNT"))
	if err != nil {
		return err
	}
	b.retention = retention

	return nil
}

func (b *BackupManager) runBackup(ctx context.Context) {
	result, err := b.backupHandler.RunBackup(ctx)
	if err != nil {
		b.log.Error(err, "Failed to run backup.")
		return
	}

	b.log.Info("Backing up eventstore has been successful.", "BackupIdentifier", result.BackupIdentifier, "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)
}

func (b *BackupManager) Schedule(ctx context.Context) error {
	b.cron = cron.NewWithLocation(time.UTC)
	if err := b.cron.AddFunc(b.cronSchedule, func() { b.runBackup(ctx) }); err != nil {
		return err
	}
	b.cron.Start()
	b.log.Info("Configured backup cronjob: %+v", b.cron.Entries())
	return nil
}

func (b *BackupManager) Close() {
	b.log.Info("Stopping...")
	b.cron.Stop()
}
