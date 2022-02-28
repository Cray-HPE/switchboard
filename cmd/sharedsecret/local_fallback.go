// MIT License
//
// (C) Copyright [2022] Hewlett Packard Enterprise Development LP
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
package sharedsecret

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"golang.org/x/sys/unix"
)

// If the cluster we are in is not set up with vault authorizations for
// the broker namespace, or if vault is somehow not working, using local
// storage for keys will, at least, let us limp along the way we did
// before we had shared keys.  Not ideal, but it gives us some backward
// compatibility.
const localStorePath string = "/tmp/localFallbackVault"

var userWarned = false

func warnUser(err error) {
	if userWarned {
		return
	}
	replicas := os.Getenv("UAI_REPLICAS")

	replicaCount, _ := strconv.Atoi(replicas)
	if replicaCount < 2 {
		return
	}
	if !userWarned {
		log.Printf("WARNING: falling back to non-shared secrets which will cause problems when multiple Broker UAI instances are used - %s\n", err)
		userWarned = true
	}
}

func localPost(vaultPath string, data *bytes.Buffer) error {
	oldMask := unix.Umask(0)
	defer unix.Umask(oldMask)
	if err := os.MkdirAll(localStorePath, 0777); err != nil {
		return fmt.Errorf("error creating local store base directory '%s' - %s", localStorePath, err)
	}
	dirPath := fmt.Sprintf("%s/%s", localStorePath, path.Dir(vaultPath))
	filePath := fmt.Sprintf("%s/%s", localStorePath, vaultPath)
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		return fmt.Errorf("error creating directory for secret storage '%s' - %s", dirPath, err)
	}
	if err := os.WriteFile(filePath, data.Bytes(), 0600); err != nil {
		return fmt.Errorf("error writing secret to local file '%s' - %s", filePath, err)
	}
	return nil
}

func localGet(vaultPath string) ([]byte, error) {
	var response VaultDataResponse
	filePath := fmt.Sprintf("%s/%s", localStorePath, vaultPath)
	if _, err := os.Stat(filePath); err != nil {
		// Vault path doesn't exist (or can't be stat'ed which is morally
		// equivalent) so return a json string with an empty map in it.
		ret, _ := json.Marshal(response)
		return ret, nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		ret, _ := json.Marshal(response)
		return ret, fmt.Errorf("error reading secret from file '%s' - %s", filePath, err)
	}
	err = json.Unmarshal(data, &response.Data)
	ret, _ := json.Marshal(response)
	return ret, err
}
