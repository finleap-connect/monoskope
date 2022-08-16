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
	"os"
	"syscall"

	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [file]",
	Short: "Decrypt backup",
	Long:  `Decrypt backup events`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		ciphertext, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		fmt.Print("Enter key: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println()

		plaintext, err := util.DecryptAES(bytePassword, ciphertext)
		if err != nil {
			return err
		}
		fmt.Println(string(plaintext))
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}
