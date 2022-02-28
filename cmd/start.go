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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"stash.us.cray.com/uas/switchboard/cmd/craycli"
	"stash.us.cray.com/uas/switchboard/cmd/uai"
	"stash.us.cray.com/uas/switchboard/cmd/util"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "SSH to an existing or newly created User Access Instance",
	Long: `The follow logic will occur with the start command:

Start a UAI if one is not already running and SSH to it once it is available.
SSH to a UAI already running if only one UAI is found
Choose a UAI to SSH to if multiple are found`,
	Run: start,
}
var image string

func start(cmd *cobra.Command, args []string) {
	var uais []uai.Uai
	var images uai.UaiImages
	var targetUai uai.Uai
	var oneShot bool

	craycli.CraycliInitialize()

	// Check for UAI_ONE_SHOT which always creates
	// and deletes the UAI after logging out
	if _, exists := os.LookupEnv("UAI_ONE_SHOT"); exists {
		oneShot = true
	} else {
		oneShot = false
	}

	// Get the list of UAIs available
	uais = uai.UaiList()

	// Get a list of allowable images
	images = uai.UaiImagesList()
	if image != "" {
		for i, img := range images.List {
			if img == image {
				break
			}
			if i == len(images.List)-1 {
				fmt.Printf("Invalid image requested: %s\n", image)
				fmt.Printf("Allowable images are: %s\n", strings.Join(images.List, ", "))
				os.Exit(1)
			}
		}
	}

	switch num := len(uais); {
	case num == 0 || oneShot:
		targetUai = uai.UaiCreate(image)

	case num == 1:
		if image == "" || image == uais[0].Image {
			// If an image wasn't specified or
			// image matches the running UAI use that one
			targetUai = uais[0]
		} else {
			// Start a new UAI if the image doesn't match
			targetUai = uai.UaiCreate(image)
		}

	default:
		if image == "" {
			// If a specific image was not requested, prompt
			// the user to select one of the running UAIs
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
		} else {
			// Attempt to find a UAI of the correct image.
			// Create one if the right image isn't running
			for i, u := range uais {
				if image == u.Image {
					targetUai = u
					break
				}
				if i == len(uais)-1 {
					targetUai = uai.UaiCreate(image)
				}
			}
		}
	}
	util.WaitForRunningReady(targetUai, "", "")
	ec := util.RunSshCmd(targetUai.ConnectionString, "")
	if oneShot {
		uai.UaiDelete(targetUai.Name)
	}
	os.Exit(ec)

}

func init() {
	startCmd.Flags().StringVar(&image, "image", "", "Name of UAI image to start")
	rootCmd.AddCommand(startCmd)
}
