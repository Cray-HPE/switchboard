// Copyright 2020 Hewlett Packard Enterprise Development LP
package cmd

import (
	"fmt"
	"os"

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
switchboard broker
switchboard list
switchboard delete`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
