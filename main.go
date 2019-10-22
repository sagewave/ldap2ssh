package main

import (
	"fmt"
	"ldap2ssh/cli"
	"ldap2ssh/config"
	"ldap2ssh/utils"
	"ldap2ssh/vault"
	"log"
)

func main() {
	if !config.Exists() {
		log.Fatal("could not find a .ldap2ssh config file")
	}

	main := config.Configuration()
	if !vault.TokenIsValid(main.VaultToken, main.VaultAddress) {
		fmt.Println("Your vault token expired. Enter your JumpCloud credentials to render a new token:")
		creds := cli.PromptJumpCloudCredentials(main.JumpCloudUser)
		token := vault.Login(creds, main.VaultAddress)
		main.VaultToken = token
	}

	// CHOOSE SSH KEY
	sshDir := utils.GetSSHDir()
	sshKeys := utils.ListPrivateKeys(sshDir)
	chosenSSHKey := cli.Select("Private Key to Sign", sshKeys)

	// SIGN SSH KEY
	signedKey := vault.SignSSHKey(chosenSSHKey, main.VaultAddress, main.VaultToken)
	log.Println("Signed Key", signedKey)

	// SAVE SSH KEY TO `chosenSSHKey`-cert.pub
}
