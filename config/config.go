package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rldw/ldap2ssh/utils"

	"gopkg.in/ini.v1"
)

// Main section of the ini file
type Main struct {
	User          string `ini:"username"`
	VaultAddress  string `ini:"vault_address"`
	VaultToken    string `ini:"vault_token"`
	VaultEndpoint string `ini:"vault_endpoint"`
	DefaultKey    string `ini:"default_key"`
}

// Exists checks if the config file exists
func Exists() bool {
	if _, err := os.Stat(filename()); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func filename() string {
	result := filepath.Join(utils.GetHomeDir(), ".ldap2ssh")
	return result
}

// Read the config file
func Read() interface{} {
	cfg, err := ini.Load(filename())
	if err != nil {
		log.Println("failed to read config file", err)
	}
	return cfg
}

// Configuration returns the main section in the ini file
func Configuration() Main {
	cfg, _ := ini.Load(filename())
	c := &Main{
		VaultToken: "",
		DefaultKey: "",
	}
	err := cfg.Section("").MapTo(c)
	if err != nil {
		log.Println("error mapping main section", err)
	}
	return *c
}

// Sections returns all section names
func Sections() []string {
	cfg, _ := ini.Load(filename())
	return filter(cfg.SectionStrings(), func(v string) bool {
		return v != "DEFAULT"
	})
}

// GetEndpoint gets endpoint from given section
func GetEndpoint(section string) string {
	cfg, _ := ini.Load(filename())
	return cfg.Section(section).Key("endpoint").String()
}

// SaveMain saves the main config section
func SaveMain(main Main) {
	cfg, _ := ini.Load(filename())
	err := ini.ReflectFrom(cfg, &main)
	cfg.SaveTo(filename())
	if err != nil {
		log.Fatal(err)
	}
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
