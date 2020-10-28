package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/metadata"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version information",
	Long:  `Prints version information and the commit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gateway\nversion: %s commit: %s\n", metadata.Version, metadata.Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
