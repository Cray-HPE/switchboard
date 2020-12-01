// Copyright 2020 Hewlett Packard Enterprise Development LP
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
