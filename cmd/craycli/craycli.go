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
package craycli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Make this a var so it can be replaced for unit testing
var execCommand = exec.Command

func CraycliInitialize() {
	if _, err := exec.LookPath("cray"); err != nil {
		fmt.Println("Could not find the cray cli command")
		os.Exit(1)
	}
	if CraycliCheckOutput("No configuration exists") {
		fmt.Println("The Cray CLI has not been initialized, running 'cray init':")
		CraycliInteractive("cray", "init")
	}
	if CraycliCheckOutput("401 Unauthorized") {
		fmt.Println("The Cray CLI has not been authorized, running 'cray auth login':")
		CraycliInteractive("cray", "auth", "login")
	}
}

func CraycliCheckOutput(cliOutput string) bool {
	cmd := execCommand("cray", "uas", "mgr-info", "list")
	stdoutStderr, _ := cmd.CombinedOutput()
	return strings.Contains(string(stdoutStderr[:]), cliOutput)
}

func CraycliInteractive(args ...string) bool {
	cmd := execCommand(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()
	ec := cmd.Wait()
	if ec != nil {
		return false
	} else {
		return true
	}
}
