package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string = "DEV"
	commit  string = "DEBUG"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version information",
	Long:  `Prints version information and the commit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("monoctl\nversion: %s commit: %s\n", version, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
