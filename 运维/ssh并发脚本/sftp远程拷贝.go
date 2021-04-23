package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	var (
		err        error
		sftpClient *sftp.Client
	)
	//start := time.Now()
	sftpClient, err = connect("root", "123.com", "192.168.1.106", 22) //远程链接打开
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	_, errStat := sftpClient.Stat("/tmp/") //远程目录确认状态是否存在
	if errStat != nil {
		log.Fatal("/tmp/" + " remote path not exists!")
	}

	//backupDirs, err := ioutil.ReadDir("c:\\1.log")
	//if err != nil {
	//	log.Fatal("c:\\1.log  local path not exists!")
	//}
	uploadFile(sftpClient, "D:\\work\\project-dev\\src\\godev\\mymodels\\ssh并发运行脚本\\1.sh", "/tmp/1.sh")
}
func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath) //打开需要上传的本地文件
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)

	}
	defer srcFile.Close()

	//var remoteFileName = path.Base(localFilePath)

	dstFile, err := sftpClient.Create(remotePath) //创建远程文件
	if err != nil {
		fmt.Println("sftpClient.Create error : ", remotePath)
		log.Fatal(err)

	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile) // 本地文件内容数据
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)

	}
	dstFile.Write(ff) //远程文件 写入 本地文件内容
	fmt.Println(localFilePath + "  copy file to remote server finished!")
}
