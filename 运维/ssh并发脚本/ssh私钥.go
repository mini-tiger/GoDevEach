package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

func main() {
	//hnowhost文件对应/root/.ssh/known_hosts。
	//var knowhost = []byte("192.168.14.137 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItb...gPfaynABbA/tD1V9pV5w=")
	//
	////只关注pubkey解析与否
	//_, _, pubKey, _, _, err := ssh.ParseKnownHosts(knowhost)
	//if err != nil {
	//	log.Fatalf("parseKnowHost error", err)
	//	return
	//}
	//fmt.Println(pubKey)

	//读取本机的私钥
	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	//获取签名
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	//设置配置文件
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			//证书验证
			ssh.PublicKeys(signer),
			//密码验证
			//ssh.Password("xxxx"),
		},
		//用于加密期间握手验证主机秘钥的回调函数。就比如第一次连接到一台主机，除了验证，还会弹出一堆信息（让你接受对方公钥）。
		// 用于接收特定主机的秘钥
		//HostKeyCallback: ssh.FixedHostKey(pubKey),
		//用于接收任何主机的秘钥
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	//准备建立客户端连接
	client, err := ssh.Dial("tcp", "172.16.8.145:10422", config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()
	//一个正式用于执行远程命令或者shell的连接
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: ", err)
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami;ls /root"); err != nil {
		log.Fatal("Failed to run: " + err.Error())

	}
	fmt.Println(b.String())

}
