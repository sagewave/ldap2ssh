package cli

import (
	"github.com/rldw/ldap2ssh/vault"

	"github.com/AlecAivazis/survey"
)

// Password ask for password
func Password(message string) string {
	val := ""
	prompt := &survey.Password{
		Message: message,
	}
	survey.AskOne(prompt, &val, nil)
	return val
}

// String ask for a string and provide a default value
func String(message string, defaultValue string) string {
	val := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
	}
	survey.AskOne(prompt, &val, nil)
	return val
}

// StringRequired ask for a required string
func StringRequired(message string) string {
	val := ""
	prompt := &survey.Input{
		Message: message,
	}
	survey.AskOne(prompt, &val, survey.WithValidator(survey.Required))
	return val
}

// Select give selection of options
func Select(message string, options []string, defaultValue string) string {
	selected := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
		Default: defaultValue,
	}
	survey.AskOne(prompt, &selected, survey.WithValidator(survey.Required))
	return selected
}

// PromptLDAPCredentials asks for username and password
func PromptLDAPCredentials(defaultUser string) vault.Credentials {
	return vault.Credentials{
		Username: String("LDAP Username", defaultUser),
		Password: Password("LDAP Password"),
	}
}
