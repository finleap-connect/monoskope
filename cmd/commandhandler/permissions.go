package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

func NewReportPermissions() *cobra.Command {
	reportPermissionsCmd := &cobra.Command{
		Use:   "permissions",
		Short: "Prints a list of permissions.",
		Long:  `Prints a list of permissions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data := [][]string{}

			for _, cmdType := range es.DefaultCommandRegistry.GetRegisteredCommandTypes() {
				command, err := es.DefaultCommandRegistry.CreateCommand(uuid.Nil, cmdType, nil)
				if err != nil {
					return err
				}
				policies := command.Policies(cmd.Context())

				for _, p := range policies {
					data = append(data, []string{
						string(cmdType),
						p.Role().String(),
						p.Scope().String(),
						fmt.Sprintf("%v", p.ResourceMatch()),
					})
				}
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Command", "Role", "Scope", "ResourceMatch"})

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

	flags := reportPermissionsCmd.Flags()
	flags.BoolVarP(&formatMarkdown, "markdown", "m", false, "Print table in markdown format.")

	return reportPermissionsCmd
}
