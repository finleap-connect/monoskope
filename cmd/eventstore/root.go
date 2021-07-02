package main

import (
	"flag"
	"os"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	_ "go.uber.org/automaxprocs"
)

var rootCmd = &cobra.Command{
	Use:          "eventstore action [flags]",
	Short:        "eventstore",
	Long:         `eventstore`,
	SilenceUsage: true,
}

func init() {
	// Setup global flags
	flags := rootCmd.PersistentFlags()
	flags.AddGoFlagSet(flag.CommandLine)
}

func main() {
	rootCmd.AddCommand(version.NewVersionCmd(rootCmd.Name()))

	if err := rootCmd.Execute(); err != nil {
		log := logger.WithName("root-cmd")
		log.Error(err, "command failed")
		os.Exit(1)
	}
}
