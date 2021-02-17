// MIT License
//
// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
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
	"stash.us.cray.com/uas/switchboard/cmd/craycli"
	"stash.us.cray.com/uas/switchboard/cmd/uai"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the running User Access Instances",
	Long: `Show the User Access Instances that are already running for the user.

$ switchboard list
#	Name			Status				Age	Image
1	uai-alanm-b1a72874	Running:Ready			4d18h	bis.local:5000/cray/cray-uas-sles15-slurm:latest
2	uai-alanm-f6a0e079	Running:Ready			19m	bis.local:5000/cray/cray-uas-sles15-slurm:latest`,
	Run: list,
}

func list(cmd *cobra.Command, args []string) {
	var uais []uai.Uai

	craycli.CraycliInitialize()

	uais = uai.UaiList()
	uai.UaiPrettyPrint(uais)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
