// MIT License
//
// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
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
package uai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

// Make this a var so it can be replaced for unit testing
var execCommand = exec.Command
var uasUrl = "http://cray-uas-mgr.services.svc.cluster.local:8088/v1/admin"

/*
A struct to represent a UAI from uas-mgr which has a json representation of
[
	{
		"uai_age": "4d17h",
		"uai_connect_string": "ssh alanm@172.30.48.15 -p 31370 -i ~/.ssh/id_rsa",
		"uai_host": "ncn-w001",
		"uai_img": "bis.local:5000/cray/cray-uas-sles15-slurm:latest",
		"uai_msg": "ContainerCreating",
		"uai_name": "uai-alanm-b1a72874",
		"uai_portmap": {},
		"uai_status": "Waiting",
		"username": "alanm"
	},
]
*/
type Uai struct {
	Name             string `json:"uai_name"`
	Username         string `json:"username"`
	ConnectionString string `json:"uai_connect_string"`
	Image            string `json:"uai_img"`
	Status           string `json:"uai_status"`
	StatusMessage    string `json:"uai_msg"`
	Host             string `json:"uai_host"`
	Age              string `json:"uai_age"`
}

type UaiImages struct {
	Default string   `json:"default_image"`
	List    []string `json:"image_list"`
}

type UaiClasses struct {
	ClassID          string `json:"class_id"`
	Comment          string `json:"comment"`
	Default          bool   `json:"default"`
	PublicSSH        bool   `json:"public_ssh"`
	UAICreationClass string `json:"uai_creation_class"`
}

// Create a UAI using default parameters
func UaiCreate(image string) Uai {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	createArgs := strings.Fields("cray uas create --format json")
	createArgs = append(createArgs, "--publickey", home+"/.ssh/id_rsa.pub")
	if image != "" {
		createArgs = append(createArgs, "--imagename", image)
	}
	cmd := execCommand(createArgs[0], createArgs[1:]...)
	var uai Uai
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Creating a new UAI...")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(stdout).Decode(&uai)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return uai
}

// Run cray uas list and decode the json into a slice of type Uai
func UaiList() []Uai {
	cmd := execCommand("cray", "uas", "list", "--format", "json")
	var uais []Uai
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(stdout).Decode(&uais); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return uais
}

// Run cray uas images list and decode the json into UaiImages
func UaiImagesList() UaiImages {
	cmd := execCommand("cray", "uas", "images", "list", "--format", "json")
	var uaiImages UaiImages
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(stdout).Decode(&uaiImages); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return uaiImages
}

func UaiClassesList() []UaiClasses {
	req, err := http.NewRequest(http.MethodGet, uasUrl+"/config/classes", nil)
	if err != nil {
		log.Fatal(err)
	}
	uasMgrClient := http.Client{
		Timeout: time.Second * 5,
	}
	res, err := uasMgrClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var uaiClasses []UaiClasses
	err = json.Unmarshal(body, &uaiClasses)
	if err != nil {
		log.Fatal(err)
	}
	return uaiClasses
}

func UaiAdminList(user string, classid string) []Uai {
	query := uasUrl + "/uais?" + "owner=" + user + "&class_id=" + classid
	req, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		log.Fatal(err)
	}
	uasMgrClient := http.Client{
		Timeout: time.Second * 5,
	}
	res, err := uasMgrClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var uais []Uai
	err = json.Unmarshal(body, &uais)
	if err != nil {
		log.Fatal(err)
	}
	return uais
}

func UaiAdminCreate(user string, classid string) Uai {

	// Get user passwd string
	passwd, err := exec.Command("getent", "passwd", user).Output()
	passwdStr := strings.TrimSuffix(string(passwd), "\n")
	if err != nil {
		log.Fatal(err)
	}

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	// Check for an SSH key. Generate one if it does not exist
	sshpublickey := home + "/.ssh/" + classid + ".pub"
	_, err = os.Stat(sshpublickey)
	if os.IsNotExist(err) {
		err := exec.Command("ssh-keygen", "-f", strings.TrimSuffix(sshpublickey, ".pub"), "-N", "").Run()
		if err != nil {
			fmt.Printf("Failed to generate an SSH key\n")
			log.Fatal(err)
		}
	}

	// Read in the public key to a request body
	file, err := os.Open(sshpublickey)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("publickey_str", classid+".pub")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(part, file)
	writer.Close()

	// Request a new UAI with the admin API
	query := uasUrl + "/uais?" + "owner=" + user + "&class_id=" + classid + "&passwd_str=" + url.QueryEscape(passwdStr)
	req, err := http.NewRequest(http.MethodPost, query, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	uasMgrClient := http.Client{
		Timeout: time.Second * 5,
	}
	fmt.Println("Creating a new UAI...")
	res, err := uasMgrClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var uai Uai
	err = json.Unmarshal(resBody, &uai)
	if err != nil {
		log.Fatal(err)
	}
	return uai
}

// Delete a UAI by name
func UaiDelete(uais string) {
	cmd := execCommand("cray", "uas", "delete", "--format", "json",
		"--uai-list", uais)
	_, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleting UAI(s)...")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

/*
Func to tabwriter a slice of type Uai in the format:
#       Name                    Status                          Age     Image
1       uai-alanm-b1a72874      Running:Ready                   4d18h   bis.local:5000/cray/cray-uas-sles15-slurm:latest
2       uai-alanm-f6a0e079      Running:Ready                   19m     bis.local:5000/cray/cray-uas-sles15-slurm:latest
*/
func UaiPrettyPrint(uais []Uai) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	if len(uais) > 0 {
		fmt.Fprintln(w, "#\tName\tStatus\tAge\tImage")
	}
	for i, u := range uais {
		fmt.Fprintf(w, "%d\t%s\t%s%s\t%s\t%s\n", i+1, u.Name,
			u.StatusMessage, u.Status,
			u.Age, u.Image)
	}
	w.Flush()
}
