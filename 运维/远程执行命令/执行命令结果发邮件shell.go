package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"远程执行命令/funcs"
	"远程执行命令/modules"
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
func (cli *sshinfo) connect() error {
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
		defer cli.close_session()
		return err
	}
	cli.Session = session
	return nil
}
func (cli *sshinfo) close_session() {
	cli.Session.Close()
}

var ResultHosts []*modules.HostMonitor = make([]*modules.HostMonitor, 0)

func main() {

	ssh1 := New_ssh(22, []string{"192.168.40.100", "root", "W3b5Ev!c3"}...)
	//fmt.Println(ssh1)
	err := ssh1.connect()
	if err != nil {
		log.Fatal(err)
	}
	//ssh.Session.Stdout=os.Stdout
	//ssh.Session.Stderr=os.Stderr
	//ssh.Session.Run("touch /root/1")
	//ssh.Session.Run("ls /; ls /tmp")
	//ssh.close_session() //todo session一次运行一次run
	hm := new(modules.HostMonitor)
	ssh1.terminal_run(hm)

	hm.Host = "192.168.40.100"
	ssh1.close_session()

	// xxx mail tpl
	funcs.MailHtml(ResultHosts)

}

func (cli *sshinfo) terminal_run(hm *modules.HostMonitor) {

	w, err := cli.Session.StdinPipe()
	if err != nil {
		panic(err)
	}
	r, err := cli.Session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	e, err := cli.Session.StderrPipe()
	if err != nil {
		panic(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal 建立终端
	if err := cli.Session.RequestPty("vt100", 40, 80, modes); err != nil { //term:xerm 是彩色显示
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	in, out := MuxShell(w, r, e, 4)             // 4条命令 包含第一个登录
	if err := cli.Session.Shell(); err != nil { //打开仿真shell
		log.Fatal(err)
	}

	// xxx login
	if k, ok := <-out; ok {
		//print("================================")
		fmt.Printf("login:\n %s\n", k) //todo 所有out通道中记录的 返回信息 打开出来
	}

	//xxx 磁盘使用率
	in <- "df -h|awk 'NR!=1 {print $6 ,$5}'"

	if result, ok := <-out; ok {
		print("===========================================\n")
		hm.DiskUsage = funcs.DiskFormat(result)
	}

	//xxx 内存使用率
	in <- "free -m|sed -n 2p|awk '{print $2,$3}'"

	if result, ok := <-out; ok {

		//fmt.Println(result)

		memlist := strings.Split(result, " ")
		//fmt.Println(memlist[0],memlist[1])
		usage := strings.Split(memlist[1], "\n")[0]
		m := new(modules.MemStatus)
		m.Total, _ = strconv.Atoi(strings.TrimSpace(memlist[0]))
		m.Usage, _ = strconv.Atoi(strings.TrimSpace(usage))
		//fmt.Println(*m)
		m.UsageRate = fmt.Sprintf("%.2f", float64(m.Usage)/float64(m.Total)*100)
		hm.Mem = m
	}

	//xxx date
	in <- "date '+%Y-%m-%d %H:%M:%S'"

	if result, ok := <-out; ok {

		//fmt.Println(result)
		dt := strings.Split(result, " ")
		t := strings.Split(dt[1], "\n")
		hm.CurrTime = fmt.Sprintf("%s %s", strings.TrimSpace(dt[0]), strings.TrimSpace(t[0]))
	}

	//log.Printf("HostMonitor:%+v\n", hm)
	a, err := json.Marshal(hm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(a))
	ResultHosts = append(ResultHosts, hm)
	//cli.Session.Close()
	//session.Wait()
	//session.Close()

}

func MuxShell(w io.Writer, r, e io.Reader, maxcmd int) (chan<- string, <-chan string) {
	in := make(chan string, 0)
	out := make(chan string, maxcmd)
	var wg sync.WaitGroup
	wg.Add(1) //shell 退出前，Shell的进程
	go func() {
		for cmd := range in { //todo in通道中 所有 需要执行的命令 依次执行
			wg.Add(1)
			w.Write([]byte(cmd + "\n")) //w 是管道输入
			wg.Wait()                   //等待命令完成
		}
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:]) //todo 标准输出管道的 内容，stdoutpipe,是io.Reader接口有reader方法，将传入的[]byte 写入
			if err != nil {
				if err == io.EOF { //如果EOF 退出
					fmt.Println("exit")
					//wg.Done()
				}
				//fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			t += n //每次命令结果 追加至buf
			result := string(buf[:t])
			if strings.Contains(result, "password:") || strings.Contains(result, "#") { //匹配是否执行完成
				out <- result
				t = 0 //t是临时存 当前命令返回的结果，清空
				wg.Done()
			}
		}
	}()
	return in, out
}
