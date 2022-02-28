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
	"stash.us.cray.com/uas/switchboard/cmd/sharedsecret"
	"time"
)

var sharedSecretRoot = composeSharedSecretRoot()

func composeSharedSecretRoot() string {
	ret := os.Getenv("UAI_SHARED_SECRET_PATH")
	if ret == "" {
		// No UAI_SHARED_SECRET_PATH in the environment means
		// that this Broker was created by a version of UAS
		// that does not know about shared secrets and will
		// not clean them up as needed, so we don't want to
		// create them.  Return a reasonable path to use with
		// the shared secret fallback to local secrets and
		// tell 'sharedsecrets' to use the local fallback
		// instead of actually sharing.
		ret = "secrets/broker-uai"
		sharedsecret.UseLocalFallback()
	}
	return ret
}

func composeHostVaultPath() string {
	return fmt.Sprintf("%s/host", sharedSecretRoot)
}

func composeInternalVaultPath(owner string) string {
	return fmt.Sprintf("%s/internal/%s", sharedSecretRoot, owner)
}

func generate(hostKey bool) (map[string]string, error) {
	keyDirPath, err := os.MkdirTemp("/tmp", "ssh-keys-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(keyDirPath)
	if err = os.Chmod(keyDirPath, 0700); err != nil {
		return nil, err
	}
	var cmdArgs []string
	if hostKey {
		cmdArgs = []string{
			"ssh-keygen",
			"-A",
			"-f", keyDirPath,
		}
		// ssh-keygen puts the keys in the directory 'keyDirPath'/etc/ssh so compose that
		keyDirPath = fmt.Sprintf("%s/etc/ssh", keyDirPath)
		if err = os.MkdirAll(keyDirPath, 0700); err != nil {
			return nil, err
		}

	} else {
		keyFile := fmt.Sprintf("%s/id_rsa", keyDirPath)
		cmdArgs = []string{
			"ssh-keygen",
			"-t", "rsa",
			"-q",
			"-N", "",
			"-f", keyFile,
		}
	}
	cmd := exec.Command("/usr/bin/ssh-keygen", cmdArgs[1:]...)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	fileMap := map[string]string{}
	keyFiles, err := os.ReadDir(keyDirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range keyFiles {
		keyPath := fmt.Sprintf("%s/%s", keyDirPath, file.Name())
		data, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
		fileMap[file.Name()] = string(data)
	}
	return fileMap, nil
}

func storeOrGetKeyMap(sharePath string, keyMap map[string]string) (map[string]string, error) {
	storedMap, err := sharedsecret.GetSecret(sharePath)
	if err != nil {
		// Couldn't get from vault, assuming that we don't have access to vault, so
		// we need to fail cleanly for the non-shared key case.  Log a warning and
		// then return the keymap we were given, so the caller can proceed.
		return storedMap, err
	}
	// The above returns an empty map if nothing is found
	if len(storedMap) != 0 {
		// There was already a map present for this path, so we can use that one.
		return storedMap, nil
	}
	// There was no map, so we get to store the one we were given, then wait a bit to
	// make sure it sticks (in case we are racing with some other broker) and retrieve
	// whatever is stored and return it.
	if err := sharedsecret.PostSecret(sharePath, keyMap); err != nil {
		// The returned map should be empty, but err should drive the caller anyway.
		return storedMap, err
	}
	time.Sleep(time.Second)
	storedMap, err = sharedsecret.GetSecret(sharePath)
	if err != nil {
		// The returned map should be empty, but err should drive the caller anyway.
		return storedMap, err
	}
	return storedMap, nil
}

func install(dirPath string, owner string, fileMap map[string]string) error {
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}
	cmd := exec.Command("/usr/bin/chown", owner, dirPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s[%s]", stderr.String(), err)
	}
	for filename, data := range fileMap {
		filePath := fmt.Sprintf("%s/%s", dirPath, filename)
		if err := os.WriteFile(filePath, []byte(data), 0600); err != nil {
			return err
		}
		cmd = exec.Command("/usr/bin/chown", owner, filePath)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s[%s]", stderr.String(), err)
		}
	}
	return nil
}
