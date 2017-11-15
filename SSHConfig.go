package main

import (
	"os"
	"log"
	"path/filepath"
	"strings"
	"golang.org/x/crypto/ssh"
	"bufio"
	"io/ioutil"
)

/*
SSHConfig ... */
type SSHConfig struct {
	ProxyHost string
	KeyFile string
	ProxyUser string
}

/*
MakeConfig ... */
func (config *SSHConfig) MakeConfig() *ssh.ClientConfig {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")

		if len(fields) != 3 {
			continue
		}

		if strings.Contains(fields[0], config.ProxyHost) {
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())

			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}

			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey for %s", config.ProxyHost)
	}

	key, err := ioutil.ReadFile(config.KeyFile)

	if err != nil {
	    log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: config.ProxyUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	return sshConfig
}
