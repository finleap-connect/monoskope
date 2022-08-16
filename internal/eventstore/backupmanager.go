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

package eventstore

import (
	"context"
	"os"
	"path"

	"github.com/finleap-connect/monoskope/internal/eventstore/backup"
	"github.com/finleap-connect/monoskope/internal/eventstore/backup/s3"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/go-logr/logr"
)

const BackupPath = "/etc/eventstore/backup"

type BackupManager struct {
	log           logr.Logger
	store         eventsourcing.EventStore
	backupHandler backup.BackupHandler
	retention     int
}

// NewBackupManager creates a new backup manager configured by config files taken from eventstore.BackupPath and environment config.
func NewBackupManager(store eventsourcing.EventStore, retention int) (*BackupManager, error) {
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
	fileInfos, err := os.ReadDir(BackupPath)
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
