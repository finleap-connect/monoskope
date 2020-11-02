package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/metadata"
)

func NewVersionCmd(cmdName string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\nversion: %s commit: %s\n", cmdName, metadata.Version, metadata.Commit)
		},
	}
}
