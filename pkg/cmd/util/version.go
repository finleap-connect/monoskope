package util

import (
	"fmt"
	"runtime"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/metadata"
)

func PrintVersion(cmdName string) {
	fmt.Printf(`%s:
		version     : %s
		commit      : %s
		go version  : %s
		go compiler : %s
		platform    : %s/%s
	`, cmdName, metadata.Version, metadata.Commit, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}
