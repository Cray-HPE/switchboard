// Copyright 2020 Hewlett Packard Enterprise Development LP
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
