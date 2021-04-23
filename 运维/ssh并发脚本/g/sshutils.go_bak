package g

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type sshinfo struct {
	IP, Username string
	Passwd       string
	Port         int
	client       *ssh.Client
	Session      *ssh.Session
	Result       string
}

func New_ssh(port int, args ...string) *sshinfo {
	temp := new(sshinfo)
	temp.Port = port
	temp.IP = args[0]
	temp.Username = args[1]
	temp.Passwd = args[2]
	return temp

}
func (cli *sshinfo) Connect() error {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(cli.Passwd))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	clientConfig := &ssh.ClientConfig{
		User:            cli.Username,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr := fmt.Sprintf("%s:%d", cli.IP, cli.Port)

	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return err
	}

	// create session
	session, err := client.NewSession()
	if err != nil {
		defer cli.Close_session()
		return err
	}
	cli.Session = session
	return nil
}
func (cli *sshinfo) Close_session() {
	cli.Session.Close()
}
