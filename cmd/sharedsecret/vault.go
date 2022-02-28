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
	"io/ioutil"
	"net/http"
	"os"
)

type VaultAuth struct {
	ClientToken string `json:"client_token"`
}

type VaultLoginResponse struct {
	Auth VaultAuth `json:"auth"`
}

type VaultDataResponse struct {
	Data map[string]string `json:"data"`
}

var useLocalFallback = false

func UseLocalFallback() {
	useLocalFallback = true
}

func authenticate(client *http.Client) (string, error) {
	if useLocalFallback {
		return "", fmt.Errorf("UAS version does not support shared secrets")
	}

	saToken, err := os.ReadFile("/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", err
	}

	authData := map[string]string{
		"jwt":  string(saToken),
		"role": "uas",
	}
	jsonAuthData, err := json.Marshal(authData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "http://cray-vault.vault:8200/v1/auth/kubernetes/login", bytes.NewBuffer(jsonAuthData))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var vaultData VaultLoginResponse
	err = json.Unmarshal([]byte(responseBody), &vaultData)
	if err != nil {
		return "", err
	}
	if vaultData.Auth.ClientToken == "" {
		return "", fmt.Errorf("failed to authenticate with vault")
	}
	return vaultData.Auth.ClientToken, nil
}

func GetSecret(vaultPath string) (map[string]string, error) {
	client := &http.Client{}
	var bodyBytes []byte
	var err error
	if vaultToken, err := authenticate(client); err == nil {
		req, err := http.NewRequest(http.MethodGet, "http://cray-vault.vault:8200/v1/"+vaultPath, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("X-Vault-Token", vaultToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		// Vault auth failed, so we are probably not allowed to talk
		// to vault on this system.  Get the secret from local storage.
		warnUser(err)
		if bodyBytes, err = localGet(vaultPath); err != nil {
			return nil, err
		}
	}
	var response VaultDataResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func PostSecret(vaultPath string, payload map[string]string) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	vaultToken, err := authenticate(client)
	if err != nil {
		// Vault auth failed, probably a system where we are not
		// authorized to use vault.  Log this and then use local
		// storage instead of vault to hold the payload.
		warnUser(err)
		return localPost(vaultPath, bytes.NewBuffer(payloadJson))
	}
	req, err := http.NewRequest(http.MethodPost, "http://cray-vault.vault:8200/v1/"+vaultPath, bytes.NewBuffer(payloadJson))
	if err != nil {
		return err
	}
	req.Header.Add("X-Vault-Token", vaultToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
