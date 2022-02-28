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
package keys

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func KeyDirPath(owner string) string {
	return fmt.Sprintf("/tmp/%s", owner)
}

func KeyFilePath(owner string) string {
	return fmt.Sprintf("%s/id_rsa", KeyDirPath(owner))
}

func SetupInternalKeys(uaiOwner string) (map[string]string, error) {
	path := composeInternalVaultPath(uaiOwner)
	keyMap, err := generate(false)
	if err != nil {
		return keyMap, fmt.Errorf("Error generating SSH keys for user '%s' - %s\n", uaiOwner, err)
	}
	// The following will either give us back an existing key map or store the one we propose and return it.
	storedKeyMap, err := storeOrGetKeyMap(path, keyMap)
	if err != nil {
		return storedKeyMap, fmt.Errorf("Error storing / retrieving SSH keys for user '%s' - %s\n", uaiOwner, err)
	}
	keyDir := KeyDirPath(uaiOwner)
	if err = install(keyDir, uaiOwner, storedKeyMap); err != nil {
		return storedKeyMap, fmt.Errorf("Error installing SSH keys for user '%s' - %s\n", uaiOwner, err)
	}
	return storedKeyMap, nil
}

func GetUaiPublicKey(uaiOwner string) (string, error) {
	keyMap, err := SetupInternalKeys(uaiOwner)
	if err != nil {
		return "", err
	}
	if val, ok := keyMap["id_rsa.pub"]; ok {
		return val, nil
	}
	return "", fmt.Errorf("SSH keymap for user '%s' has no 'id_rsa.pub' key -- should not happen", uaiOwner)
}

func GetHostKeys(uaiOwner string, uaiIP string) (string, error) {
	knownHosts := fmt.Sprintf("%s/known_hosts", KeyDirPath(uaiOwner))
	outfile, err := os.Create(knownHosts)
	if err != nil {
		return "", err
	}
	defer outfile.Close()
	if err := outfile.Chmod(0600); err != nil {
		return "", err
	}
	cmd := exec.Command("/usr/bin/ssh-keyscan", uaiIP)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = outfile
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s[%s]", stderr.String(), err)
	}
	return knownHosts, nil
}
