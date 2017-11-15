package main

import (
	"fmt"
	"net"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

/*
SSHTunnel ... */
type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
}

/*
Start ... */
func (tunnel *SSHTunnel) Start(c chan int) error {
	if (tunnel.Local.Proto == "unix") {
		if _, err := os.Stat(tunnel.Local.String()); !os.IsNotExist(err) {
			panic("Socket file exists")
		}
	}
	listener, err := net.Listen(tunnel.Local.Proto, tunnel.Local.String())

	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			return err
		}

		go tunnel.forward(conn)
	}

	// c <- 0
	// return nil
}

/*
Close ... */
func (tunnel *SSHTunnel) Close() {
	fmt.Printf("Cleaning up socket: %s\n", tunnel.Local.Path)
	os.Remove(tunnel.Local.Path)
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)

	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())

	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)

		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
