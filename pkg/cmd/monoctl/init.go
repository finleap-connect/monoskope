package monoctl

import (
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
	"sigs.k8s.io/kind/pkg/errors"
)

var (
	serverURL string
)

func NewInitCmd(configLoader *config.ClientConfigLoader) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init monoctl",
		Long:  `Init monoctl and create a new monoskope configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURL == "" {
				return errors.Errorf("server-url is required")
			}

			config := config.NewConfig()
			config.Server = serverURL

			if err := configLoader.InitConifg(config); err != nil {
				return err
			}

			return nil
		},
	}
	flags := initCmd.Flags()
	flags.StringVar(&serverURL, "server-url", "", "URL of the monoskope instance")
	return initCmd
}
