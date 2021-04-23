package main

import (
	"io/ioutil"
	"net"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

func getKeyFile() (key ssh.Signer, err error) {
	//usr, _ := user.Current()
	file := "Path to your key file(.pem)"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func main() {
	key, err := getKeyFile()
	if err != nil {
		panic(err)
	}

	// Define the Client Config as :
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", "<remote ip>:22", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	err = scp.CopyPath("local file path", "remote path", session)
	if err != nil {
		panic("Failed to Copy: " + err.Error())
	}
	defer session.Close()
}
