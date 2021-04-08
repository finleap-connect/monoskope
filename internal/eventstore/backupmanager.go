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
	log           logr.Logger
	store         eventsourcing.Store
	backupHandler backup.BackupHandler
	retention     int
}

// NewBackupManager creates a new backup manager configured by config files taken from eventstore.BackupPath and environment config.
func NewBackupManager(store eventsourcing.Store, retention int) (*BackupManager, error) {
	manager := &BackupManager{
		log:       logger.WithName("backup-manager"),
		store:     store,
		retention: retention,
	}
	if err := manager.configure(); err != nil {
		return nil, err
	}
	return manager, nil
}

func (bm *BackupManager) configure() error {
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
			bm.backupHandler = s3.NewS3BackupHandler(conf, bm.store, bm.retention)
		}
	}

	return nil
}

func (bm *BackupManager) RunBackup(ctx context.Context) (*backup.BackupResult, error) {
	return bm.backupHandler.RunBackup(ctx)
}

func (bm *BackupManager) RunPurge(ctx context.Context) (*backup.PurgeResult, error) {
	return bm.backupHandler.RunPurge(ctx)
}

func (bm *BackupManager) RunRestore(ctx context.Context, backupIdentifier string) (*backup.RestoreResult, error) {
	return bm.backupHandler.RunRestore(ctx, backupIdentifier)
}
