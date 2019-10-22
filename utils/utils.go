package utils

import (
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func ListPrivateKeys(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}

	var filtered []string
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if IsValidPrivateKey(path) {
			filtered = append(filtered, file.Name())
		}
	}

	return filtered
}

func IsValidPrivateKey(filename string) bool {
	privateBytes, err := ioutil.ReadFile(filename)
	_, err = ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		// log.Print(err)
		return false
	}
	return true
}

func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Could not determine home directory", err)
	}
	return usr.HomeDir
}

func GetSSHDir() string {
	home := GetHomeDir()
	return filepath.Join(home, ".ssh")
}
