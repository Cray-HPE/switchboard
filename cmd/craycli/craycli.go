package craycli

import (
	"os"
	"os/exec"
	"strings"
)

// Make this a var so it can be replaced for unit testing
var execCommand = exec.Command

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
