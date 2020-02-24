package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
        "stash.us.cray.com/uan/switchboard/cmd/craycli"
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
	if ! craycli.CraycliCmdExists() {
		fmt.Println("Could not find the cray cli command")
		os.Exit(1)
	}
	if craycli.CraycliCheckOutput("No configuration exists") {
		fmt.Println("The Cray CLI has not been initialized, running 'cray init':")
		craycli.CraycliInteractive("cray", "init")
	}
	if craycli.CraycliCheckOutput("401 Unauthorized") {
		fmt.Println("The Cray CLI has not been authorized, running 'cray auth login':")
		craycli.CraycliInteractive("cray", "auth", "login")
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

