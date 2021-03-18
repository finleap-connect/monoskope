package main

import (
	"github.com/spf13/cobra"
)

func NewReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "report",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Short:                 "Reporting of implementation details",
		Long:                  `Reporting of implementation details`,
	}
}
