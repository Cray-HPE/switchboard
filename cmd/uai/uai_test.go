package uai

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"encoding/json"
)

var oneUai = `
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
}`

var twoUai = `
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
	{
		"uai_age": "15s",
		"uai_connect_string": "ssh alanm@172.30.48.15 -p 12345 -i ~/.ssh/id_rsa",
		"uai_host": "ncn-w002",
		"uai_img": "bis.local:5000/cray/cray-uas-sles15-slurm:latest",
		"uai_name": "uai-alanm-deadbeef",
		"uai_portmap": {},
		"uai_status": "Running: Ready",
		"username": "alanm"
	}
]`

var deleteUai = `["Successfully deleted uai-alanm-ea059360"]`

var prettyPrintOutput = `
#       Name                    Status                          Age     Image
1       uai-alanm-b1a72874      Running:Ready                   4d17h   bis.local:5000/cray/cray-uas-sles15-slurm:latest
2       uai-alanm-deadbeef      Running:Ready                   15s     bis.local:5000/cray/cray-uas-sles15-slurm:latest
`

func fakeExecCommand(command string, args...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	switch cs[4] {
	case "create":
		os.Setenv("HELPER_CMD_OUTPUT", oneUai)
	case "list":
		os.Setenv("HELPER_CMD_OUTPUT", twoUai)
	case "delete":
		os.Setenv("HELPER_CMD_OUTPUT", deleteUai)
	}
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, os.Getenv("HELPER_CMD_OUTPUT"))
	os.Exit(0)
}

func TestUai(t *testing.T) {
	var uai Uai
	err := json.Unmarshal([]byte(oneUai), &uai)
	if (err != nil) {
		t.Errorf("Could not decode oneUai")
	}
	var uais []Uai
	err = json.Unmarshal([]byte(twoUai), &uais)
	if (err != nil) {
		t.Errorf("Could not decode twoUai")
	}
	if len(uais) != 2 {
		t.Errorf("Expected two Uais to be returned")
	}
}

func TestUaiCreate(t *testing.T) {
	execCommand = fakeExecCommand
	defer func(){ execCommand = exec.Command }()
	var newUai Uai
	newUai = UaiCreate()
	if (newUai.Name != "uai-alanm-b1a72874") {
		t.Errorf("Failed to decode a Uai from UaiCreate()")
	}
}

func TestUaiList(t *testing.T) {
	execCommand = fakeExecCommand
	defer func(){ execCommand = exec.Command }()
	var uais []Uai
	uais = UaiList()
	if (uais[0].Name != "uai-alanm-b1a72874") {
		t.Errorf("Expected the second Uai to be 'uai-alanm-b1a72874'")
	}
	if (uais[1].Name != "uai-alanm-deadbeef") {
		t.Errorf("Expected the second Uai to be 'uai-alanm-deadbeef'")
	}
}

func TestUaiDelete(t *testing.T) {
	execCommand = fakeExecCommand
	defer func(){ execCommand = exec.Command }()
	UaiDelete("uai-alanm-b1a72874")
}

func TestUaiPrettyPrint(t *testing.T) {
	execCommand = fakeExecCommand
        defer func(){ execCommand = exec.Command }()
	UaiPrettyPrint(UaiList())
}
