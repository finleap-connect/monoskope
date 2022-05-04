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

package backup

import (
	"context"
)

type BackupHandler interface {
	// RunBackup creates a backup of all events in the store
	RunBackup(context.Context) (*BackupResult, error)
	// RunRestore restores all events stored in the backup with the given identifier
	RunRestore(context.Context, string) (*RestoreResult, error)
	// RunPurge cleans up backups according to the retention set
	RunPurge(context.Context) (*PurgeResult, error)
}

type BackupResult struct {
	ProcessedEvents  uint64
	ProcessedBytes   uint64
	BackupIdentifier string
}

type RestoreResult struct {
	ProcessedEvents uint64
	ProcessedBytes  uint64
}

type PurgeResult struct {
	PurgedBackups int
	BackupsLeft   int
}
