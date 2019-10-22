package cli

import (
	"ldap2ssh/vault"

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
func Select(message string, options []string) string {
	selected := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	survey.AskOne(prompt, &selected, survey.WithValidator(survey.Required))
	return selected
}

// PromptJumpCloudCredentials asks for username and password
func PromptJumpCloudCredentials(defaultUser string) vault.Credentials {
	username := String("JumpCloud Username", defaultUser)
	password := Password("JumpCloud Password")
	return vault.Credentials{
		username,
		password,
	}
}
