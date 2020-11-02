package monoctl

import (
	"flag"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/cmd/monoctl/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/cmd/version"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
)

var (
	configLoader *config.ClientConfigLoader
)

func NewRootCmd() *cobra.Command {
	configLoader = config.NewLoader()

	rootCmd := &cobra.Command{
		Use:          "monoctl action [flags]",
		Short:        "monoctl",
		Long:         `monoctl`,
		SilenceUsage: true,
	}

	// Setup global flags
	flags := rootCmd.PersistentFlags()
	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&configLoader.ExplicitFile, "monoconfig", "", "Path to the monoskope config file to use for CLI requests")

	rootCmd.AddCommand(version.NewVersionCmd(rootCmd.Name()))
	rootCmd.AddCommand(auth.NewAuthCmd(configLoader))

	return rootCmd
}
