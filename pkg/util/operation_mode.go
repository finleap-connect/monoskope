package util

import "os"

type OperationMode string

const (
	RELEASE     OperationMode = "release"
	DEVELOPMENT OperationMode = "development"
)

// GetOperationMode returns the operation mode specified via the env var M8_OPERATION_MODE.
// (defaults is RELEASE)
func GetOperationMode() OperationMode {
	operationMode := OperationMode(os.Getenv("M8_OPERATION_MODE"))
	switch operationMode {
	case DEVELOPMENT:
		return DEVELOPMENT
	default:
		return RELEASE
	}
}
