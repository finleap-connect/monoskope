package main

import (
	"flag"
	"time"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/util"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
)

var (
	explicitFile string
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "monoctl action [flags]",
		Short:        "monoctl",
		Long:         `monoctl`,
		SilenceUsage: true,
	}

	// Setup global flags
	fl := rootCmd.PersistentFlags()
	fl.AddGoFlagSet(flag.CommandLine)
	fl.StringVar(&explicitFile, "monoconfig", "", "Path to the monoskope config file to use for CLI requests")
	fl.DurationVar(&util.Timeout, "command-timeout", 120*time.Second, "Timeout for long running commands, defaults to 120s")

	configManager := config.NewLoaderFromExplicitFile(explicitFile)
	rootCmd.AddCommand(NewVersionCmd(rootCmd.Name(), configManager))
	rootCmd.AddCommand(NewInitCmd(configManager))
	rootCmd.AddCommand(auth.NewAuthCmd(configManager))

	return rootCmd
}
