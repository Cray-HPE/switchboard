package cmd

import (
	"github.com/spf13/cobra"
	"stash.us.cray.com/uan/switchboard/cmd/uai"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list the running User Access Instances",
	Long: `Show the User Access Instances that are already running for the user.

$ switchboard list
#	Name			Status				Age	Image
1	uai-alanm-b1a72874	Running:Ready			4d18h	bis.local:5000/cray/cray-uas-sles15-slurm:latest
2	uai-alanm-f6a0e079	Running:Ready			19m	bis.local:5000/cray/cray-uas-sles15-slurm:latest`,
	Run: list,
}

func list(cmd *cobra.Command, args []string) {
	var uais []uai.Uai
	uais = uai.UaiList()
	uai.UaiPrettyPrint(uais)
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
