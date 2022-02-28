// MIT License
//
// (C) Copyright [2020-2022] Hewlett Packard Enterprise Development LP
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
	"stash.us.cray.com/uas/switchboard/cmd/keys"
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
		for i, cls := range classes {
			if cls.ClassID == classid {
				break
			}
			if i == len(classes)-1 {
				fmt.Printf("Invalid class requested: %s\n", classid)
				os.Exit(1)
			}
		}
	} else {
		for i, cls := range classes {
			if cls.Default {
				classid = cls.ClassID
				break
			}
			if i == len(classes)-1 {
				fmt.Printf("No --class-id was provided and no default class is configured.")
				os.Exit(1)
			}
		}
	}

	// Set up the user's internal SSH session keys
	if _, err := keys.SetupInternalKeys(user.Username); err != nil {
		log.Fatal(err)
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
	knownHosts, err := keys.GetHostKeys(user.Username, targetUai.IP)
	if err != nil {
		log.Fatalf("error preloading internal host keys - %s", err)
	}
	sshCmd := fmt.Sprintf("%s -o TCPKeepalive=true -o UserKnownHostsFile=%s", targetUai.ConnectionString, knownHosts)
	ec := util.RunSshCmd(sshCmd, keys.KeyFilePath(user.Username))
	os.Exit(ec)

}

func init() {
	brokerCmd.Flags().StringVar(&classid, "class-id", "", "Specify a UAI class ID")
	rootCmd.AddCommand(brokerCmd)
}
