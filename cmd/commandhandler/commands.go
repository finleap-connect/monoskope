package main

import (
	"os"
	"sort"

	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

func NewReportCommands() *cobra.Command {
	reportCommandsCmd := &cobra.Command{
		Use:   "commands",
		Short: "Prints a list of commands.",
		Long:  `Prints a list of commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data := [][]string{}

			var cmdTypes []string
			for _, v := range es.DefaultCommandRegistry.GetRegisteredCommandTypes() {
				cmdTypes = append(cmdTypes, v.String())
			}
			sort.Strings(cmdTypes)

			for _, cmdType := range cmdTypes {
				command, err := es.DefaultCommandRegistry.CreateCommand(uuid.Nil, es.CommandType(cmdType), nil)
				if err != nil {
					return err
				}
				data = append(data, []string{
					cmdType,
					command.AggregateType().String(),
				})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Command", "Aggregate"})

			if formatMarkdown {
				table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
				table.SetAutoMergeCellsByColumnIndex([]int{0})
				table.SetCenterSeparator("|")
			} else {
				table.SetAutoWrapText(false)
				table.SetAutoFormatHeaders(true)
				table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				table.SetCenterSeparator("")
				table.SetColumnSeparator("")
				table.SetRowSeparator("")
				table.SetHeaderLine(false)
				table.SetBorder(false)
				table.SetTablePadding("\t") // pad with tabs
				table.SetNoWhiteSpace(true)
			}

			table.AppendBulk(data) // Add Bulk Data
			table.Render()

			return nil
		},
	}

	flags := reportCommandsCmd.Flags()
	flags.BoolVarP(&formatMarkdown, "markdown", "m", false, "Print table in markdown format.")

	return reportCommandsCmd
}
