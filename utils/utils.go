package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func ListPublicKeys(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}

	var filtered []string
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if IsValidPublicKey(path) &&
			strings.HasSuffix(path, ".pub") &&
			!strings.HasSuffix(path, "-cert.pub") {
			filtered = append(filtered, file.Name())
		}
	}

	return filtered
}

func IsValidPublicKey(filename string) bool {
	publicBytes, err := ioutil.ReadFile(filename)
	_, _, _, _, err = ssh.ParseAuthorizedKey(publicBytes)
	if err != nil {
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

func ValidateCert(certfile string) string {
	publicBytes, err := ioutil.ReadFile(certfile)
	k, _, _, _, err := ssh.ParseAuthorizedKey(publicBytes)
	if err != nil {
		log.Fatal(err)
	}

	cert := k.(*ssh.Certificate)
	validBefore := fmt.Sprintf("%v", cert.ValidBefore)
	i, err := strconv.ParseInt(validBefore, 10, 64)
	tm := time.Unix(i, 0)
	return fmt.Sprintf("%v", tm)
}
