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
