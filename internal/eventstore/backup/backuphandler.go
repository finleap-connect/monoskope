package backup

import (
	"context"
)

type BackupHandler interface {
	// RunBackup creates a backup of all events in the store
	RunBackup(context.Context) (*Result, error)
	// RunRestore restores all events stored in the backup with the given identifier
	RunRestore(context.Context, string) (*Result, error)
}

type Result struct {
	ProcessedEvents  uint64
	ProcessedBytes   uint64
	BackupIdentifier string
}
