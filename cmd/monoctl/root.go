package main

import (
	"flag"
	"time"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/auth"
	conf "gitlab.figo.systems/platform/monoskope/monoskope/cmd/monoctl/config"
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
	fl.DurationVar(&util.Timeout, "command-timeout", 120*time.Second, "Timeout for long running commands")

	confLoader := config.NewLoaderFromExplicitFile(explicitFile)
	rootCmd.AddCommand(NewVersionCmd(rootCmd.Name(), confLoader))
	rootCmd.AddCommand(conf.NewConfigCmd(confLoader))
	rootCmd.AddCommand(auth.NewAuthCmd(confLoader))

	return rootCmd
}
