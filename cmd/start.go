package cmd

import (
	"bufio"
	"os"
	"os/exec"
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"

	"github.com/briandowns/spinner"
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
var s *spinner.Spinner

// SpinnerStart will start a spinner/waiting message that will go until we stop it
func SpinnerStart(message string) {
	fmt.Print(fmt.Sprintf("%s...", message))
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
}

// SpinnerStop will stop/cancel a spinner
func SpinnerStop() {
	if s.Active() {
		s.Stop()
	}
	fmt.Println()
}

// TODO introduce a timeout to this endless function
// TODO pass in a uai instead to immediately check
//	status (it may already be Running: Ready)
func waitForRunningReady(uaiName string) {
        var uais []uai.Uai
        var status string
	SpinnerStart("Waiting for UAI to be ready")
	for (status != "Running: Ready") {
		uais = uai.UaiList()
		for _,uai := range uais {
			if (uaiName == uai.Name) {
				status = uai.StatusMessage + uai.Status
			}
		}
		time.Sleep(1 * time.Second)
	}
	SpinnerStop()
}

func runSshCmd(sshCmd string) {
	sshArgs := strings.Fields(sshCmd)
	sshExec := exec.Command(sshArgs[0], sshArgs[1:]...)
	sshExec.Stdout = os.Stdout
	sshExec.Stdin = os.Stdin
	sshExec.Stderr = os.Stderr
	sshExec.Run()
	//TODO return correct exit code from ssh
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
		if uais[0].StatusMessage + uais[0].Status != "Running: Ready" {
			waitForRunningReady(uais[0].Name)
		}
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
	runSshCmd(sshCmd)
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
