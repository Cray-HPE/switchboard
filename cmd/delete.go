package cmd

import (
	"bufio"
	"os"
	"fmt"
	"strconv"
	"strings"
	"log"

	"github.com/spf13/cobra"
	"stash.us.cray.com/uan/switchboard/cmd/uai"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete User Access Instances currently running",
	Long: `Delete User Access Instances currently running.`,
	Run: delete,
}

// TODO Make this work for multiple UAIs
func delete(cmd *cobra.Command, args []string) {
        var uais []uai.Uai
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
