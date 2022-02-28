// MIT License
//
// (C) Copyright [2022] Hewlett Packard Enterprise Development LP
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
	"github.com/spf13/cobra"
	"log"
	"os"
	"stash.us.cray.com/uas/switchboard/cmd/keys"
)

var hostkeyCmd = &cobra.Command{
	Use:   "hostkey",
	Short: "Create, Distribute and Install a Broker UAI SSH Host Key",
	Long: `Create an SSH Host Key for a Broker UAI, Register it with
Vault so replicas can share it, and install it, and install it in the Broker UAI.
If a host key for this Broker UAI is already registered with Vault, use that one.
`,
	Run: hostkey,
}

func hostkey(cmd *cobra.Command, args []string) {
	_, err := keys.SetupHostKeys()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(hostkeyCmd)
}
