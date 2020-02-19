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

func waitForRunningReady(targetUai uai.Uai) {
        var uais []uai.Uai
        var status string
	if targetUai.StatusMessage + targetUai.Status == "Running: Ready" {
		return
	}
	SpinnerStart("Waiting for UAI to be ready")
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-timeout:
			SpinnerStop()
			fmt.Printf("Timeout waiting on %s to be 'Running: Ready'\n", targetUai.Name)
			fmt.Printf("Last status was '%s'\n", status)
			os.Exit(1)
		case <-tick:
			uais = uai.UaiList()
			for _,uai := range uais {
				if (targetUai.Name == uai.Name) {
					status = uai.StatusMessage + uai.Status
					break
				}
			}
			if (status == "Running: Ready") {
				SpinnerStop()
				return
			}
		}
	}
}

// There has to be a better way to do this but
// It was not obvious how to get only the error number
// rather than "exit code 123". (improve this later)
func convertErrorStrToInt(err error) int {
	var errNum int
	if err != nil {
		errS := strings.Fields(err.Error())
		errNum, _ = strconv.Atoi(errS[2])
	} else {
		errNum = 0
	}
	return errNum
}

func runSshCmd(sshCmd string) int {
	sshArgs := strings.Fields(sshCmd)
        if sshOriginalCommand, exists := os.LookupEnv("SSH_ORIGINAL_COMMAND"); exists {
                sshArgs = append(sshArgs, sshOriginalCommand)
        }
	sshExec := exec.Command(sshArgs[0], sshArgs[1:]...)
	sshExec.Stdout = os.Stdout
	sshExec.Stdin = os.Stdin
	sshExec.Stderr = os.Stderr
	sshExec.Start()
	ec := sshExec.Wait()
	return convertErrorStrToInt(ec)
}

func start(cmd *cobra.Command, args []string) {
        var uais []uai.Uai
	var sshCmd string
	var freshUai uai.Uai
	var oneShot bool

	// Get the list of UAIs available
        uais = uai.UaiList()

	// Check for SWITCHBOARD_ONE_SHOT which always creates
	// and deletes the UAI after logging out
	if _, exists := os.LookupEnv("SWITCHBOARD_ONE_SHOT"); exists {
		oneShot = true
	} else {
		oneShot = false
	}
	switch num := len(uais); {

	// No UAI is running so start up a fresh one
	case num == 0 || oneShot:
		freshUai = uai.UaiCreate()
		waitForRunningReady(freshUai)
		sshCmd = freshUai.ConnectionString

	// A single UAI is running, use this one
	case num == 1:
		waitForRunningReady(uais[0])
		sshCmd = uais[0].ConnectionString

	// Multiple UAIs running, prompt the user to pick one
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
		waitForRunningReady(uais[selection-1])
		sshCmd = uais[selection-1].ConnectionString
	}
	ec := runSshCmd(sshCmd)
	if oneShot {
		uai.UaiDelete(freshUai.Name)
	}
	os.Exit(ec)

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
