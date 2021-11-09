package main

import (
	"bytes"
	"fmt"
	"github.com/masterzen/winrm"
)

//winrm quickconfig
//y
//winrm set winrm/config/service/Auth '@{Basic="true"}'
//winrm set winrm/config/service '@{AllowUnencrypted="true"}'
//winrm set winrm/config/winrs '@{MaxMemoryPerShellMB="1024"}'

func main() {
	endpoint := winrm.NewEndpoint("192.168.6.180", 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, "Administrator", "123.com")
	if err != nil {
		panic(err)
	}

	//_, err = client.RunWithInput("ipconfig", os.Stdout, os.Stderr, os.Stdin)
	//if err != nil {
	//	panic(err)
	//}

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	_, err = client.Run(winrm.Powershell("ipconfig /all"), &buf, &errBuf)
	if err != nil {
		panic(err)
	}
	if len(errBuf.Bytes()) > 0 {
		panic("zero ")
	}
	fmt.Println(buf.String())
}
