package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "switchboard",
	Short: "Switchboard is a tool to redirect users into a User Access Instance.",
	Long: `Switchboard will automate the process of creating, listing, and deleting 
User Access Instances (UAIs). In addition to running the necessary 'cray 
uas' commands, switchboard will make sure the user is authenticated to 
the Shasta system. 

The following commands are supported:
switchboard start
switchboard list
switchboard delete)`,
}

func init() {
	checkCmdExists("cray")
}

// Check for craycli
func checkCmdExists(cmd string) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("Could not find the command %s\n", cmd)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

