package main

import (
	"github.com/spf13/cobra"
)

var (
	formatMarkdown bool
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
