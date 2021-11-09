package main

import (
	"bytes"
	"github.com/masterzen/winrm"
	"io"
	"os"
)

func main() {

	stdin := bytes.NewBufferString("ipconfig /all\n")
	endpoint := winrm.NewEndpoint("192.168.6.180", 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, "Administrator", "123.com")
	if err != nil {
		panic(err)
	}
	shell, err := client.CreateShell()
	if err != nil {
		panic(err)
	}
	var cmd *winrm.Command
	cmd, err = shell.Execute("cmd.exe")
	if err != nil {
		panic(err)
	}

	go io.Copy(cmd.Stdin, stdin)
	go io.Copy(os.Stdout, cmd.Stdout)
	go io.Copy(os.Stderr, cmd.Stderr)

	cmd.Wait()

	shell.Close()
}
