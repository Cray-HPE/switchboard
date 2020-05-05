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

var image string

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
	var images uai.UaiImages
	var sshCmd string
	var targetUai uai.Uai
	var oneShot bool

	// Check for UAI_ONE_SHOT which always creates
	// and deletes the UAI after logging out
	if _, exists := os.LookupEnv("UAI_ONE_SHOT"); exists {
		oneShot = true
	} else {
		oneShot = false
	}

	// Get the list of UAIs available
	uais = uai.UaiList()

	// Get a list of allowable images
	images = uai.UaiImagesList()
	if image != "" {
		for i,img := range images.List {
			if (img == image) {
				break
			}
			if i == len(images.List)-1 {
				fmt.Printf("Invalid image requested: %s\n", image)
				fmt.Printf("Allowable images are: %s\n", strings.Join(images.List, ", "))
				os.Exit(1)
			}
		}
	}

	switch num := len(uais); {
	case num == 0 || oneShot:
		targetUai = uai.UaiCreate(image)

	case num == 1:
		if image == "" || image == uais[0].Image {
			// If an image wasn't specified or
			// image matches the running UAI use that one
			targetUai = uais[0]
		} else {
			// Start a new UAI if the image doesn't match
			targetUai = uai.UaiCreate(image)
		}

	default:
		if image == "" {
			// If a specific image was not requested, prompt
			// the user to select one of the running UAIs
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
			targetUai = uais[selection-1]
		} else {
			// Attempt to find a UAI of the correct image.
			// Create one if the right image isn't running
			for i,u := range uais {
				if (image == u.Image) {
					targetUai = u
					break
				}
				if i == len(uais)-1 {
					targetUai = uai.UaiCreate(image)
				}
			}
		}
	}
	waitForRunningReady(targetUai)
	sshCmd = targetUai.ConnectionString
	ec := runSshCmd(sshCmd)
	if oneShot {
		uai.UaiDelete(targetUai.Name)
	}
	os.Exit(ec)

}

func init() {
	startCmd.Flags().StringVar(&image, "image", "", "Name of UAI image to start")
	rootCmd.AddCommand(startCmd)
}
