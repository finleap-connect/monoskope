package main

import (
	"os"

	monoctl "gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/cmd"
)

func main() {
	if err := monoctl.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
