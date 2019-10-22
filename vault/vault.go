package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Credentials for JumpCloud
type Credentials struct {
	Username string
	Password string
}

func makeRequest(url string, payload []byte, vaultToken string) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if vaultToken != "" {
		req.Header.Set("X-Vault-Token", vaultToken)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func makeGetRequest(url string, vaultToken string) []byte {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if vaultToken != "" {
		req.Header.Set("X-Vault-Token", vaultToken)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

// Login to vault
func Login(creds Credentials, vaultAddr string) string {
	url := fmt.Sprintf("%s/v1/auth/ldap/login/%s", vaultAddr, creds.Username)
	jsonStr := fmt.Sprintf(`{"password":"%s"}`, creds.Password)
	var jsonBytes = []byte(jsonStr)

	body := makeRequest(url, jsonBytes, "")
	log.Println(body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	token := result["auth"].(map[string]interface{})["client_token"]
	return fmt.Sprintf("%v", token)
}

// TokenIsValid validates a given Vault token and checks its expiration
func TokenIsValid(token string, vaultAddr string) bool {
	url := fmt.Sprintf("%s/v1/auth/token/lookup-self", vaultAddr)
	body := makeGetRequest(url, token)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	rawTTL := result["data"].(map[string]interface{})["ttl"]
	ttl := fmt.Sprintf("%f", rawTTL)
	casted, _ := strconv.ParseFloat(ttl, 32)
	return casted > 0
}

// SignSSHKey returns a signed SSH Key
func SignSSHKey(keyfile string, vaultAddr string, vaultToken string) string {
	return ""
}
