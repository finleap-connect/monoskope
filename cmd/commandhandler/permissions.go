package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
)

var (
	formatMarkdown bool
)

func NewReportPermissions() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "permissions",
		Short: "Prints a list of permissions.",
		Long:  `Prints a list of permissions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data := [][]string{}

			commandRegistry := domain.RegisterCommands()
			types := commandRegistry.GetRegisteredCommandTypes()

			for _, cmdType := range types {
				command, err := commandRegistry.CreateCommand(cmdType, nil)
				if err != nil {
					return err
				}
				policies := command.Policies(cmd.Context())

				for _, p := range policies {
					res := p.Resource()
					sub := p.Subject()

					if res == "" {
						res = "self"
					}

					if sub == "" {
						sub = "self"
					}

					data = append(data, []string{
						string(cmdType),
						p.Role().String(),
						p.Scope().String(),
						res,
						sub,
					})
				}
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Command", "Role", "Scope", "Resource", "Subject"})

			if formatMarkdown {
				table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
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

	flags := loginCmd.Flags()
	flags.BoolVarP(&formatMarkdown, "markdown", "m", false, "Print table in markdown format.")

	return loginCmd
}
