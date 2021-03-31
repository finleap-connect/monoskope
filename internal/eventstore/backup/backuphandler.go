package backup

import "context"

type BackupHandler interface {
	RunBackup(context.Context) (*BackupResult, error)
}

type BackupResult struct {
	ProcessedEvents uint64
	ProcessedBytes  uint64
}
