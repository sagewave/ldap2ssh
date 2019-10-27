package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/rldw/ldap2ssh/utils"

	"gopkg.in/ini.v1"
)

// Section of the config file
type Section struct {
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
	return filepath.Join(utils.GetHomeDir(), ".ldap2ssh")
}

// GetSection returns a section of the Config
func GetSection(name string) Section {
	cfg, _ := ini.Load(filename())
	sec := &Section{
		VaultToken: "",
		DefaultKey: "",
		User:       "",
	}
	err := cfg.Section(name).MapTo(sec)
	if err != nil {
		log.Fatalf("Error mapping section %s to Section", name)
	}
	return *sec
}

// Sections returns all section names
func Sections() []string {
	cfg, _ := ini.Load(filename())
	return filter(cfg.SectionStrings(), func(v string) bool {
		return v != "DEFAULT"
	})
}

// SectionExists checks if a section with the given name exists in the config
func SectionExists(name string) bool {
	sections := Sections()
	return any(sections, func(v string) bool {
		return v == name
	})
}

// CreateEmpty config file
func CreateEmpty(name string) {
	err := ioutil.WriteFile(name, []byte(""), 0600)
	if err != nil {
		log.Fatalf("Unable to write new config file: %v", err)
	}
}

// SaveSection saves a given section to the config file
func SaveSection(name string, sec *Section) {
	if !Exists() {
		CreateEmpty(filename())
	}

	cfg, _ := ini.Load(filename())
	err := cfg.Section(name).ReflectFrom(sec)
	if err != nil {
		log.Fatal("Could not save section to config file: ", err)
	}
	cfg.SaveTo(filename())
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

func any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}
