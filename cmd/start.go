package cmd

import (
	"bufio"
	"os"
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"
	"github.com/spf13/cobra"
        "stash.us.cray.com/uan/switchboard/cmd/uai"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "SSH to an existing or newly created User Access Instance",
	Long: `The follow logic will occur with the start command:

Start a UAI if one is not already running and SSH to it once it is available.
SSH to a UAI already running if only one UAI is found
Choose a UAI to SSH to if multiple are found`,
	Run: start,
}

/*func runSshCmd(sshCmd string) {
	var args = strings.Fields(sshCmd)
	env := os.Environ()
	fmt.Printf("Running ssh command: '%q'", args)
	cmd := exec.Command(args[0], args[1:len(args)]..., env)
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)

}*/

func waitForRunningReady(uaiName string) {
        var uais []uai.Uai
        var status string
	for (status != "Running: Ready") {
		//fmt.Printf("Running UaiList\n")
		uais = uai.UaiList()
		for _,uai := range uais {
			//fmt.Printf("Searching for UAI\n")
			if (uaiName == uai.Name) {
				status = uai.StatusMessage + uai.Status
				//fmt.Printf("status: '%s'\n", status)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func start(cmd *cobra.Command, args []string) {
        var uais []uai.Uai
	var sshCmd string
        uais = uai.UaiList()
	switch num := len(uais); num {
	case 0:
		freshUai := uai.UaiCreate()
		waitForRunningReady(freshUai.Name)
		sshCmd = freshUai.ConnectionString
	case 1:
		waitForRunningReady(uais[0].Name)
		sshCmd = uais[0].ConnectionString
	default:
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
		waitForRunningReady(uais[selection-1].Name)
		sshCmd = uais[selection-1].ConnectionString
	}
	fmt.Printf("SSH Connection string:\n%s\n", sshCmd)
	//runSshCmd(sshCmd)
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
