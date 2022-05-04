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

package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Name    string = "APP"   // set by args
	Version string = "DEV"   // set by args when building, see go.mk
	Commit  string = "DEBUG" // set by args when building, see go.mk
)

func NewVersionCmd(cmdName string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit`,
		Run: func(cmd *cobra.Command, args []string) {
			PrintVersion(cmdName)
		},
	}
}

func PrintVersion(cmdName string) {
	fmt.Println(cmdName)
	fmt.Printf(" version     : %s\n", Version)
	fmt.Printf(" commit      : %s\n", Commit)
	fmt.Printf(" go version  : %s\n", runtime.Version())
	fmt.Printf(" go compiler : %s\n", runtime.Compiler)
	fmt.Printf(" platform    : %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
