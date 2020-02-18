package cmd

import (
  "fmt"
  "os"
  "os/exec"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
  "github.com/spf13/cobra"
)
var cfgFile string

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
  cobra.OnInitialize(initConfig)
  checkCmdExists("cray")

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.switchboard.yaml)")


  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  //rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".switchboard" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".switchboard")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
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

