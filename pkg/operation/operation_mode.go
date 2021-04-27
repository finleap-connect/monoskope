package operation

import (
	"fmt"
	"os"
)

type OperationMode string

const (
	RELEASE     OperationMode = "release"
	CMDLINE     OperationMode = "cmdline"
	DEVELOPMENT OperationMode = "development"
)

func init() {
	_ = GetOperationMode()
}

// GetOperationMode returns the operation mode specified via the env var M8_OPERATION_MODE.
// (defaults is RELEASE)
func GetOperationMode() OperationMode {
	operationMode := OperationMode(os.Getenv("M8_OPERATION_MODE"))

	switch operationMode {
	case DEVELOPMENT:
		fmt.Print("################ WARNING ###############\n> OPERATION MODE IS SET TO DEVELOPMENT.\n> SENSIBLE INFORMATION MIGHT BE LEAKED!\n########################################\n")
		return DEVELOPMENT
	case CMDLINE:
		return CMDLINE
	default:
		return RELEASE
	}
}
