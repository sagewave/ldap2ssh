package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/rldw/ldap2ssh/cli"
	"github.com/rldw/ldap2ssh/config"
	"github.com/rldw/ldap2ssh/utils"
	"github.com/rldw/ldap2ssh/vault"
)

var (
	// Version tool version
	Version = "0.1"
)

func main() {
	app := kingpin.New("ldap2ssh", "A command line tool for generating SSH certificates with vault.")
	app.Version(Version)

	// general settings
	verbose := app.Flag("verbose", "Verbose output mode").Bool()

	// `configure` command
	cmdConfigure := app.Command("configure", "Configure a new account.")
	configureFlags := new(ConfigureFlags)
	cmdConfigure.Flag("account", "The account name to save this configuration to.").Short('a').StringVar(&configureFlags.Account)
	cmdConfigure.Flag("user", "The default user name to use for singing in with Vault.").Short('u').StringVar(&configureFlags.User)
	cmdConfigure.Flag("vaultaddress", "The complete Vault address including the protocol, e.g. https://vault.example.com").Short('d').StringVar(&configureFlags.VaultAddress)
	cmdConfigure.Flag("vaultendpoint", "The Vault endpoint to use to sign the SSH key, e.g. /ssh-client-signer/sign/ca").Short('e').StringVar(&configureFlags.VaultEndpoint)
	cmdConfigure.Flag("defaultkey", "The default key to sign").Short('k').StringVar(&configureFlags.DefaultKey)

	// `sign` command
	cmdSign := app.Command("sign", "Sign a public key using Vault.")
	signFlags := new(SignFlags)
	cmdSign.Flag("account", "The account to generate a SSH certificate for (env: LDAP2SSH_ACCOUNT)").Envar("LDAP2SSH_ACCOUNT").Short('a').StringVar(&signFlags.Account)
	cmdSign.Flag("key", "The SSH public key to sign (env: LDAP2SSH_KEY)").StringVar(&signFlags.Key)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *verbose {
		log.Println("debug logging activated")
	}

	switch command {
	case cmdSign.FullCommand():
		sign()
	case cmdConfigure.FullCommand():
		configure(configureFlags)
	}
}

func sign() {
	if !config.Exists() {
		log.Fatal("could not find a .ldap2ssh config file")
	}

	main := config.Configuration()
	if vault.TokenIsValid(main.VaultToken, main.VaultAddress) {
		fmt.Print("Using existing Vault token found in ~/.ldap2ssh\n\n")
	} else {
		fmt.Print("Missing or expired Vault token. Enter your JumpCloud credentials to render a new token:\n\n")
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
		fmt.Printf("Using default key %s\n\n", chosenSSHKey)
	}
	keyfile := filepath.Join(sshDir, chosenSSHKey)

	// CHOOSE ACCOUNT
	accounts := config.Sections()
	account := cli.Select("Account", accounts, accounts[0])
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

func configure(configureFlags *ConfigureFlags) {
	fmt.Println("account: " + configureFlags.Account)
	fmt.Println("user: " + configureFlags.User)
	fmt.Println("vault address: " + configureFlags.VaultAddress)
	fmt.Println("vault endpoint: " + configureFlags.VaultEndpoint)
	fmt.Println("default key: " + configureFlags.DefaultKey)
}
