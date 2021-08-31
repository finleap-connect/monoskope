// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"sort"

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
				policies := command.Policies(cmd.Context())

				for _, p := range policies {
					data = append(data, []string{
						cmdType,
						p.Role().String(),
						p.Scope().String(),
					})
				}
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Command", "Role", "Scope"})

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
