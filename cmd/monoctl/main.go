package main

import (
	"os"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/cmd/monoctl"
)

func main() {
	if err := monoctl.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
