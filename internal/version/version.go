package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version string = "DEV"   // set by args when building, see go.mk
	Commit  string = "DEBUG" // set by args when building, see go.mk
)

func NewVersionCmd(cmdName string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		Run: func(cmd *cobra.Command, args []string) {
			PrintVersion(cmdName)
		},
	}
}

func PrintVersion(cmdName string) {
	fmt.Printf(`%s:
		version     : %s
		commit      : %s
		go version  : %s
		go compiler : %s
		platform    : %s/%s
	`, cmdName, Version, Commit, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}
