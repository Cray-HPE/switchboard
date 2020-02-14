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
	Short: "A brief description of your command",
	Long: `Not Implemented.`,
	Run: delete,
}

// TODO Make this work for multiple UAIs
func delete(cmd *cobra.Command, args []string) {
        var uais []uai.Uai
        uais = uai.UaiList()
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

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
