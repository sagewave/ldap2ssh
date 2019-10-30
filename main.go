package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/rldw/ldap2ssh/cli"
	"github.com/rldw/ldap2ssh/config"
	"github.com/rldw/ldap2ssh/utils"
	"github.com/rldw/ldap2ssh/vault"
)

var (
	// Version tool version
	Version = "0.2"
)

func init() {
	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	app := kingpin.New("ldap2ssh", "A command line tool for generating SSH certificates with vault.")
	app.Version(Version)
	app.Author("Rene Ludwig")
	app.HelpFlag.Short('h')

	// general settings
	verbose := app.Flag("verbose", "Verbose output mode").Short('v').Bool()

	// `configure` command
	cmdConfigure := app.Command("configure", "Configure a new account.")
	configureFlags := new(ConfigureFlags)
	cmdConfigure.Flag("account", "The account name to save this configuration to.").Short('a').StringVar(&configureFlags.Account)
	cmdConfigure.Flag("ldap-user", "The default user name to use for singing in with Vault.").Short('u').StringVar(&configureFlags.User)
	cmdConfigure.Flag("ssh-user", "The ssh user to login as on the remote host.").Short('s').StringVar(&configureFlags.SSHUser)
	cmdConfigure.Flag("vault-address", "The complete Vault address including the protocol, e.g. https://vault.example.com").Short('d').StringVar(&configureFlags.VaultAddress)
	cmdConfigure.Flag("vault-endpoint", "The Vault endpoint to use to sign the SSH key, e.g. /ssh-client-signer/sign/ca").Short('e').StringVar(&configureFlags.VaultEndpoint)
	cmdConfigure.Flag("default-key", "The default key to sign").Short('k').StringVar(&configureFlags.DefaultKey)

	// `sign` command
	cmdSign := app.Command("sign", "Sign a public key using Vault.")
	signFlags := new(SignFlags)
	cmdSign.Flag("account", "The account to generate a SSH certificate for").Short('a').StringVar(&signFlags.Account)
	cmdSign.Flag("key", "The SSH public key to sign").Short('k').StringVar(&signFlags.Key)
	cmdSign.Flag("token", "Pass a Vault token instead of using one in the config").Short('t').StringVar(&signFlags.Token)
	cmdSign.Flag("force", "Force Vault token renewal").Short('f').BoolVar(&signFlags.Force)
	cmdSign.Flag("outfile", "Path to save the certificate to").Short('o').StringVar(&signFlags.Outfile)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logging activated")
		log.Debug("Home directory: ", utils.GetHomeDir())
		log.Debug("SSH directory: ", utils.GetSSHDir())
	}

	switch command {
	case cmdSign.FullCommand():
		sign(signFlags)
	case cmdConfigure.FullCommand():
		configure(configureFlags)
	}
}

func sign(signFlags *SignFlags) {
	if !config.Exists() {
		log.Fatal("Could not find a config file, run `ldap2ssh configure` first.")
	}

	// CHOOSE ACCOUNT
	account := signFlags.Account
	if account == "" {
		accounts := config.Sections()
		account = cli.Select("Account", accounts, accounts[0])
	}

	if !config.SectionExists(account) {
		log.Fatalf("Account '%s' does not exist.", account)
	}

	sec := config.GetSection(account)
	if vault.TokenIsValid(sec.VaultToken, sec.VaultAddress) && !signFlags.Force {
		fmt.Print("Using existing Vault token found in ~/.ldap2ssh\n\n")
	} else {
		fmt.Print("Missing or expired Vault token. Enter your JumpCloud credentials to render a new token:\n\n")
		creds := cli.PromptLDAPCredentials(sec.User)
		sec.VaultToken = vault.Login(creds, sec.VaultAddress)
		config.SaveSection(account, &sec)
	}

	// CHOOSE SSH KEY
	chosenSSHKey := sec.DefaultKey
	sshDir := utils.GetSSHDir()
	if sec.DefaultKey == "" {
		sshKeys := utils.ListPublicKeys(sshDir)
		chosenSSHKey = cli.Select("Public Key to Sign", sshKeys, sec.DefaultKey)
	} else {
		fmt.Printf("Using configured default key %s\n\n", chosenSSHKey)
	}
	keyfile := filepath.Join(sshDir, chosenSSHKey)

	// SSH USERNAME
	sshUser := sec.SSHUser
	if sshUser == "" {
		sshUser = cli.StringRequired("SSH Username")
	}

	log.Debugf("Using ssh user %s", sshUser)

	// SIGN SSH KEY AND SAVE TO CERT
	endpoint := sec.VaultEndpoint
	signedKey := vault.SignSSHKey(keyfile, sshUser, endpoint, sec.VaultAddress, sec.VaultToken)
	outfile := strings.TrimSuffix(keyfile, ".pub") + "-cert.pub"
	if signFlags.Outfile != "" {
		outfile = signFlags.Outfile
	}

	cert := []byte(signedKey)
	ioutil.WriteFile(outfile, cert, 0600)
	fmt.Println("\nWrote signed key to", outfile)

	validUntil := utils.ValidateCert(outfile)
	fmt.Println("Certificate is valid until", validUntil)
}

func configure(configureFlags *ConfigureFlags) {
	sec := &config.Section{
		User:          configureFlags.User,
		VaultAddress:  configureFlags.VaultAddress,
		VaultEndpoint: configureFlags.VaultEndpoint,
		DefaultKey:    configureFlags.DefaultKey,
	}

	log.Debugf("section: %v", sec)
	config.SaveSection(configureFlags.Account, sec)
}
