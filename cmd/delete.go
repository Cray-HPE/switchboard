// MIT License
//
// (C) Copyright [2020-2022] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"stash.us.cray.com/uas/switchboard/cmd/craycli"
	"stash.us.cray.com/uas/switchboard/cmd/uai"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete User Access Instances currently running",
	Long:  `Delete User Access Instances currently running.`,
	Run:   delete,
}

// TODO Make this work for multiple UAIs
func delete(cmd *cobra.Command, args []string) {
	var uais []uai.Uai

	craycli.CraycliInitialize()

	uais = uai.UaiList()
	if len(uais) > 0 {
		uai.UaiPrettyPrint(uais)
		fmt.Printf("Select a UAI by number: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		selection, err := strconv.Atoi(strings.TrimSuffix(input, "\n"))
		if err != nil {
			log.Fatal(err)
		}
		if (selection <= 0) || (selection > len(uais)) {
			log.Fatal("Number was not valid")
		}
		uai.UaiDelete(uais[selection-1].Name)
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
