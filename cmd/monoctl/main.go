package main

import (
	"fmt"
	"os"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/cmd/monoctl"
)

func main() {
	if err := monoctl.NewRootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
