// Copyright 2020 Hewlett Packard Enterprise Development LP
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"stash.us.cray.com/uas/switchboard/cmd/uai"
	"stash.us.cray.com/uas/switchboard/cmd/util"
)

var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "SSH to an existing or newly created User Access Instance",
	Long: `The follow logic will occur with the start command:

Start a UAI if one is not already running and SSH to it once it is available.
SSH to a UAI already running if only one UAI is found
Choose a UAI to SSH to if multiple are found`,
	Run: broker,
}

var classid string

func broker(cmd *cobra.Command, args []string) {
	var uais []uai.Uai
	var classes []uai.UaiClasses
	var targetUai uai.Uai

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of allowable classes
	classes = uai.UaiClassesList()
	if classid != "" {
		for i,cls := range classes {
			if (cls.ClassID == classid) {
				break
			}
			if i == len(classes)-1 {
				fmt.Printf("Invalid class requested: %s\n", classid)
				os.Exit(1)
			}
		}
	} else {
		for i,cls := range classes {
			if (cls.Default) {
				classid = cls.ClassID
				break
			}
			if i == len(classes)-1 {
				fmt.Printf("No --class-id was provided and no default class is configured.")
				os.Exit(1)
			}
		}
	}


	// Get the list of UAIs available
	uais = uai.UaiAdminList(user.Username, classid)

	switch num := len(uais); {
	case num == 0:
		targetUai = uai.UaiAdminCreate(user.Username, classid)

	case num == 1:
		targetUai = uais[0]

	case num > 1:
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

	}

	util.WaitForRunningReady(targetUai, user.Username, classid)
	ec := util.RunSshCmd(targetUai.ConnectionString, "~/.ssh/"+classid)
	os.Exit(ec)

}

func init() {
	brokerCmd.Flags().StringVar(&classid, "class-id", "", "Specify a UAI class ID")
	rootCmd.AddCommand(brokerCmd)
}
