package main

import (
	"fmt"
	"io/ioutil"
	"syscall"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"golang.org/x/term"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [file]",
	Short: "Decrypt backup",
	Long:  `Decrypt backup events`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		ciphertext, err := ioutil.ReadFile(filename)
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
