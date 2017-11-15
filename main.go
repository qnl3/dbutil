package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"runtime"

	"github.com/go-yaml/yaml"
)

/*
ConfigMember ... */
type ConfigMember struct {
	Name         string
	ProxyHost    string `yaml:"proxy_host"`
	ProxyUser    string `yaml:"proxy_user"`
	KeyFile      string `yaml:"key_file"`
	Host         string
	Port         int
	LocalAddress string `yaml:"local_address"`
	LocalProto   string `yaml:"local_proto"`
}

/*
Log ... */
func (config *ConfigMember) Log() {
	fmt.Printf("\n---\nConnection: %s\n", config.Name)
	fmt.Printf("%s -> %s\n---\n", config.LocalAddress, config.Host)
}

func main() {
	defer os.Exit(0)

	var yamlType []ConfigMember

	source, err := ioutil.ReadFile("/Users/quinn/.go/src/github.com/quinn/dbutil/config.yml")

	if err != nil {
		panic(err)
	}

	err = yaml.UnmarshalStrict(source, &yamlType)

	if err != nil {
	    panic(err)
	}

	c := make(chan int)

	for _, config := range yamlType {
		sshConfig := &SSHConfig{
			ProxyHost: config.ProxyHost,
			ProxyUser: config.ProxyUser,
			KeyFile:   config.KeyFile,
		}

		tunnel := &SSHTunnel{
			Config: sshConfig.MakeConfig(),

			Local: &Endpoint{
				Proto: config.LocalProto,
				Path:  config.LocalAddress,
			},

			Server: &Endpoint{
				Proto: "tcp",
				Host:  config.ProxyHost,
				Port:  22,
			},

			Remote: &Endpoint{
				Proto: "tcp",
				Host:  config.Host,
				Port:  config.Port,
			},
		}

		config.Log()
		defer tunnel.Close()
		go tunnel.Start(c)

		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

		// Block until a signal is received.
		signal := <- signalChannel
		fmt.Printf("\nGot signal: %s\n\n", signal)
		runtime.Goexit()
	}

	for _ = range yamlType { _ = <- c }

	if err != nil {
		fmt.Printf("io.Copy error: %s", err)
	}
}
