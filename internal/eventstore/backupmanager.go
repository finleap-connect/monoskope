package eventstore

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/go-logr/logr"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/backup"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/backup/s3"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const BackupPath = "/etc/eventstore/backup"

type BackupManager struct {
	log              logr.Logger
	store            eventsourcing.Store
	backupHandler    backup.BackupHandler
	metricsPublisher backup.MetricsPublisher
	retention        int
}

// NewBackupManager creates a new backup manager configured by config files taken from eventstore.BackupPath and environment config.
func NewBackupManager(metricsPublisher backup.MetricsPublisher, store eventsourcing.Store, retention int) (*BackupManager, error) {
	manager := &BackupManager{
		log:              logger.WithName("backup-manager"),
		store:            store,
		retention:        retention,
		metricsPublisher: metricsPublisher,
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

	return nil
}

func (b *BackupManager) RunBackup(ctx context.Context) (*backup.Result, error) {
	return b.backupHandler.RunBackup(ctx)
}
