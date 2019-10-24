package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Credentials for JumpCloud
type Credentials struct {
	Username string
	Password string
}

func makePostRequest(url string, payload []byte, vaultToken string) []byte {
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

	body := makePostRequest(url, jsonBytes, "")
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	token := result["auth"].(map[string]interface{})["client_token"]
	return fmt.Sprintf("%v", token)
}

// TokenIsValid validates a given Vault token and checks its expiration
func TokenIsValid(token string, vaultAddr string) bool {
	if token == "" {
		return false
	}

	url := fmt.Sprintf("%s/v1/auth/token/lookup-self", vaultAddr)
	body := makeGetRequest(url, token)
	var result map[string]interface{}

	err := json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	rawTTL, ok := result["data"].(map[string]interface{})
	if !ok {
		return false
	}
	ttl := fmt.Sprintf("%f", rawTTL["ttl"])
	casted, _ := strconv.ParseFloat(ttl, 32)
	return casted > 0
}

// SignSSHKey returns a signed SSH Key
func SignSSHKey(keyfile string, endpoint string, vaultAddr string, vaultToken string) string {
	jsonStr := `{
		"public_key": "%s",
		"valid_principals": "ec2-user",
		"extension": {
			"permit-pty": "",
			"permit-agent-forwarding": "",
			"permit-port-forwarding": ""
		}
	}`

	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Panic("could not open key file", err)
	}
	keystring := string(key)
	keystring = strings.TrimSpace(keystring)
	url := fmt.Sprintf("%s/v1/%s", vaultAddr, endpoint)
	formatted := fmt.Sprintf(jsonStr, keystring)
	payload := []byte(formatted)

	body := makePostRequest(url, payload, vaultToken)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		log.Fatal("could not parse body", result, err)
	}
	signedKey := data["signed_key"]
	return fmt.Sprintf("%v", signedKey)
}
