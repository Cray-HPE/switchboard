package uai

import (
        "fmt"
        "log"
        "os"
        "os/exec"
        "text/tabwriter"
        "encoding/json"
)

// A struct to represent a UAI from uas-mgr
type Uai struct {
	Name string `json:"uai_name"`
	Username string `json:"username"`
	ConnectionString string `json:"uai_connect_string"`
	Image string `json:"uai_img"`
	Status string `json:"uai_status"`
	StatusMessage string `json:"uai_msg"`
	Host string `json:"uai_host"`
	Age string `json:"uai_age"`
}

// Create a UAI using default parameters (TODO make it configurable)
func UaiCreate() Uai {
	// TODO fix path to ~
	cmd := exec.Command("cray", "uas", "create", "--format", "json",
			 "--publickey", "/Users/alanm/.ssh/id_rsa.pub")
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
        if (err != nil) {
                log.Fatal(err)
        }
        if err := cmd.Wait(); err != nil {
                log.Fatal(err)
        }
	return uai
}

// Run cray uas list and decode the json into a slice of type Uai
func UaiList() []Uai {
	cmd := exec.Command("cray", "uas", "list", "--format", "json")
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

// Func to tabwriter a slice of type Uai
func UaiPrettyPrint(uais []Uai) {
        w := new(tabwriter.Writer)
        w.Init(os.Stdout, 0, 8, 0, '\t', 0)
        if len(uais) > 0 {
                fmt.Fprintln(w, "#\tName\tStatus\tAge\tImage")
        }
        for i,u := range uais {
                fmt.Fprintf(w, "%d\t%s\t%s%s\t%s\t%s\n", i+1, u.Name,
                                u.StatusMessage, u.Status,
                                u.Age, u.Image)
        }
        w.Flush()
}

/*
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
