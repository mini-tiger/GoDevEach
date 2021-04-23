package modules

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: modules
 * @File:  modules
 * @Version: 1.0.0
 * @Date: 2021/2/22 上午9:12
 */

type MemStatus struct {
	Total     int
	Usage     int
	UsageRate string
}

type HostMonitor struct {
	Host      string
	DiskUsage map[string]string
	Mem       *MemStatus
	CurrTime  string
}

type SSHConn struct {
	Host      string
	User      string
	Pass      string
	Type      string
	KeyPath   string
	Port      int
	SSHClient *ssh.Client
	Session   *ssh.Session
}

func NewSSHConn() *SSHConn {
	return &SSHConn{
		Host:    "",
		User:    "root",
		Pass:    "",
		Type:    "password",
		KeyPath: "",
		Port:    22,
	}
}
func (s *SSHConn) CloseSession() {

	//s.SSHClient.Close()
	s.Session.Close()
}
func (s *SSHConn) CloseSSH() {

	//s.SSHClient.Close()
	s.SSHClient.Close()
}

func (s *SSHConn) GetConn() error {

	//

	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            s.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if s.Type == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(s.Pass)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(s.KeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	SSHClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	s.SSHClient = SSHClient
	//defer sshClient.Close()

	//defer session.Close()
	//执行远程命令
	//combo,err := session.CombinedOutput("whoami; cd /; ls -al;echo https://github.com/dejavuzhou/felix")
	//if err != nil {
	//	log.Fatal("远程执行cmd 失败",err)
	//}
	return nil

}

func (s *SSHConn) GetSession() error {
	//创建ssh-session
	Session, err := s.SSHClient.NewSession()
	if err != nil {
		return err
	}
	s.Session = Session
	return nil
}

func (s *SSHConn) RunCmd(cmd string) ([]byte, error) {

	return s.Session.CombinedOutput(cmd)

}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
