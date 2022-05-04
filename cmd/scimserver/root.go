// Copyright 2022 Monoskope Authors
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
	"flag"
	"os"

	"github.com/finleap-connect/monoskope/internal/version"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "scimserver action [flags]",
	Short:        "scimserver",
	Long:         `scimserver`,
	SilenceUsage: true,
}

func init() {
	// Setup global flags
	flags := rootCmd.PersistentFlags()
	flags.AddGoFlagSet(flag.CommandLine)
}

func main() {
	rootCmd.AddCommand(version.NewVersionCmd(rootCmd.Name()))

	if err := rootCmd.Execute(); err != nil {
		log := logger.WithName("root-cmd")
		log.Error(err, "command failed")
		os.Exit(1)
	}
}
