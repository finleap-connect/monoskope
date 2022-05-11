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
	"fmt"
	"time"

	"github.com/finleap-connect/monoskope/internal/scimserver"
	"github.com/finleap-connect/monoskope/internal/test"
	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

var debugCmd = &cobra.Command{
	Use:   "debug [flags]",
	Short: "Starts the debug server",
	Long:  `Starts the SCIM server in debug mode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		baseTestEnv := test.NewTestEnv("scimserver-testenv")
		testEnv, err := scimserver.NewTestEnv(baseTestEnv)

		shutdown := util.NewShutdownWaitGroup()

		// Start routine waiting for signals
		shutdown.RegisterSignalHandler(func() {
			util.PanicOnError(testEnv.Shutdown())
		})
		shutdown.Wait()

		if !shutdown.IsExpected() && err != nil {
			panic(fmt.Sprintf("shutdown unexpected: %v", err))
		}

		// Check if we are expecting shutdown
		// Wait for both shutdown signals and close the channel
		if ok := shutdown.WaitOrTimeout(30 * time.Second); !ok {
			panic("shutting down gracefully exceeded 30 seconds")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
