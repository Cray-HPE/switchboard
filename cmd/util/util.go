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
package util

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"stash.us.cray.com/uas/switchboard/cmd/uai"
)

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

// There has to be a better way to do this but
// It was not obvious how to get only the error number
// rather than "exit code 123". (improve this later)
func ConvertErrorStrToInt(err error) int {
	var errNum int
	if err != nil {
		errS := strings.Fields(err.Error())
		errNum, _ = strconv.Atoi(errS[2])
	} else {
		errNum = 0
	}
	return errNum
}

func RunSshCmd(sshCmd string, sshpublickey string) int {
	sshArgs := strings.Fields(sshCmd)
	if sshpublickey != "" {
		sshArgs = append(sshArgs, "-i")
		sshArgs = append(sshArgs, sshpublickey)
	}
	if sshOriginalCommand, exists := os.LookupEnv("SSH_ORIGINAL_COMMAND"); exists {
		sshArgs = append(sshArgs, sshOriginalCommand)
	}
	sshExec := exec.Command(sshArgs[0], sshArgs[1:]...)
	sshExec.Stdout = os.Stdout
	sshExec.Stdin = os.Stdin
	sshExec.Stderr = os.Stderr
	sshExec.Start()
	ec := sshExec.Wait()
	return ConvertErrorStrToInt(ec)
}

func WaitForRunningReady(targetUai uai.Uai, user string, classid string) {
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
			if user == "" && classid == "" {
				uais = uai.UaiList()
			} else {
				uais = uai.UaiAdminList(user, classid)
			}
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
