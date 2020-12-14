package version

import (
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/cmd/util"
)

func NewVersionCmd(cmdName string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		Run: func(cmd *cobra.Command, args []string) {
			util.PrintVersion(cmdName)
		},
	}
}
