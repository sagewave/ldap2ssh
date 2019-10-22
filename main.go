package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/rldw/ldap2ssh/cli"
	"github.com/rldw/ldap2ssh/config"
	"github.com/rldw/ldap2ssh/utils"
	"github.com/rldw/ldap2ssh/vault"
)

func main() {
	if !config.Exists() {
		log.Fatal("could not find a .ldap2ssh config file")
	}

	main := config.Configuration()
	if vault.TokenIsValid(main.VaultToken, main.VaultAddress) {
		fmt.Println("Using existing Vault token.")
	} else {
		fmt.Println("Missing or expired Vault token. Enter your JumpCloud credentials to render a new token:")
		creds := cli.PromptJumpCloudCredentials(main.JumpCloudUser)
		main.VaultToken = vault.Login(creds, main.VaultAddress)
		config.SaveMain(main)
	}

	// CHOOSE SSH KEY
	chosenSSHKey := main.DefaultKey
	sshDir := utils.GetSSHDir()
	if main.DefaultKey == "" {
		sshKeys := utils.ListPublicKeys(sshDir)
		chosenSSHKey = cli.Select("Public Key to Sign", sshKeys, main.DefaultKey)
	} else {
		log.Printf("Using default_key = %s", chosenSSHKey)
	}
	keyfile := filepath.Join(sshDir, chosenSSHKey)

	// CHOOSE SECTION
	sections := config.Sections()
	account := cli.Select("Account", sections, "")
	endpoint := config.GetEndpoint(account)

	// SIGN SSH KEY AND SAVE TO CERT
	signedKey := vault.SignSSHKey(keyfile, endpoint, main.VaultAddress, main.VaultToken)
	keyName := strings.TrimSuffix(keyfile, ".pub")
	certfile := keyName + "-cert.pub"
	cert := []byte(signedKey)
	ioutil.WriteFile(certfile, cert, 0600)
	fmt.Println("\nWrote signed key to", certfile)

	validUntil := utils.ValidateCert(certfile)
	fmt.Println("Certificate is valid until", validUntil)
}
