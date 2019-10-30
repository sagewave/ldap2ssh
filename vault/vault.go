package vault

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/tidwall/gjson"
)

// Credentials for LDAP
type Credentials struct {
	Username string
	Password string
}

// Login to vault
func Login(creds Credentials, vaultAddr string) string {
	url := fmt.Sprintf("%s/v1/auth/ldap/login/%s", vaultAddr, creds.Username)
	jsonStr := fmt.Sprintf(`{"password":"%s"}`, creds.Password)
	var jsonBytes = []byte(jsonStr)

	body := makePostRequest(url, jsonBytes, "")
	token := gjson.Get(body, "auth.client_token")
	err := gjson.Get(body, "errors")
	if err.Exists() {
		log.Fatal("Vault returned an error during the login: ", err.String())
	}

	return token.String()
}

// TokenIsValid validates a given Vault token and checks its expiration
func TokenIsValid(token string, vaultAddr string) bool {
	if token == "" {
		log.Debug("Empty token given for validation")
		return false
	}

	url := buildURL(vaultAddr, "v1/auth/token/lookup-self")
	body := makeGetRequest(url, token)
	err := gjson.Get(body, "errors")
	if err.Exists() {
		log.Debug("Could not validate Vault token: ", err.String())
		return false
	}

	ttl := gjson.Get(body, "data.ttl")
	log.Debug("Vault token TTL: ", ttl.Int())
	return ttl.Int() > 0
}

// SignSSHKey returns a signed SSH Key
func SignSSHKey(keyfile string, sshUser string, endpoint string, vaultAddr string, vaultToken string) string {
	jsonStr := `{
		"public_key": "%s",
		"valid_principals": "%s",
		"extension": {
			"permit-pty": "",
			"permit-agent-forwarding": "",
			"permit-port-forwarding": ""
		}
	}`

	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Panic("Could not read key file", err)
	}

	keystring := string(key)
	keystring = strings.TrimSpace(keystring)
	url := buildURL(vaultAddr, "v1", endpoint)
	formatted := fmt.Sprintf(jsonStr, keystring, sshUser)
	payload := []byte(formatted)

	body := makePostRequest(url, payload, vaultToken)
	vaultErr := gjson.Get(body, "errors")
	if vaultErr.Exists() {
		log.Fatal("Could not create certificate: ", vaultErr.String())
	}

	signedKey := gjson.Get(body, "data.signed_key")
	return signedKey.String()
}

func buildURL(base string, fragments ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		log.Panic("Could not parse base url", err)
	}

	joinedFragments := path.Join(fragments...)
	u.Path = path.Join(u.Path, joinedFragments)
	return u.String()
}

func makePostRequest(url string, payload []byte, vaultToken string) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if vaultToken != "" {
		req.Header.Set("X-Vault-Token", vaultToken)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	log.Debugf("Making post request to %v", url)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making post request: %v", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("Post request response body: %v", string(body))
	return string(body)
}

func makeGetRequest(url string, vaultToken string) string {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if vaultToken != "" {
		req.Header.Set("X-Vault-Token", vaultToken)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	log.Debugf("Making get request to %v", url)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making get request to: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("Get request response body: %v", string(body))
	return string(body)
}
