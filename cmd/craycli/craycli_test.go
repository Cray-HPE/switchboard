// Copyright 2020 Hewlett Packard Enterprise Development LP
package craycli

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

const craycliNoInit = `
Usage: cray uas mgr-info list [OPTIONS]

Error: No configuration exists. Run cray init
`

const craycliNoAuth = `
Usage: cray uas mgr-info list [OPTIONS]
Try "cray uas mgr-info list --help" for help.

Error: Error received from server: 401 Unauthorized
`

func fakeExecCommand(command string, args...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T){
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }
    fmt.Fprintf(os.Stderr, craycliNoInit)
    os.Exit(0)
}

func TestCraycliCheckOutput(t *testing.T) {
	execCommand = fakeExecCommand
	defer func(){ execCommand = exec.Command }()
	if rc := CraycliCheckOutput("No configuration exists"); ! rc {
		t.Errorf("Expected 'true' searching for 'No configuration exists' but got %t", rc)
	}
	if rc := CraycliCheckOutput("gobbledegook"); rc {
		t.Errorf("Expected 'false' searching for 'gobbledegook' but got %t", rc)
	}
}

func TestCraycliInteractive(t *testing.T) {
	if ! CraycliInteractive("/bin/sh", "-c", ">&2 echo stderr on stdin; exit 0") {
		t.Errorf("Expected a true value from 'exit 0'")
	}
	if CraycliInteractive("/bin/sh", "-c", ">&2 echo stderr on stdin; exit 1") {
		t.Errorf("Expected a false value from 'exit 1'")
	}
}
