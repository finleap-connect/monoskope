package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/metadata"
)

func NewVersionCmd(cmdName string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`%s:
		version     : %s
		commit      : %s
		go version  : %s
		go compiler : %s
		platform    : %s/%s
`, cmdName, metadata.Version, metadata.Commit, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
		},
	}
}
