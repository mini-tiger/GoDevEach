package funcs

import (
	"gitee.com/taojun319/tjtools/nmap"
	"gitee.com/taojun319/tjtools/sshtools"
	mapset "github.com/deckarep/golang-set"
	"log"
	"ssh并发脚本/g"
)

var hostSet mapset.Set
var passwdSet mapset.Set
var HostPass *nmap.SafeMap = nmap.NewSafeMap() // 成功后 主机与密码 MAP
var FailHosts mapset.Set = mapset.NewSet()     // SSH失败的主机

var hostchan chan struct{}

func SSHRun() {

	hosts := g.Config().Hosts
	hostSet = mapset.NewSetFromSlice(hosts) // 去重
	hostchan = make(chan struct{}, hostSet.Cardinality())
	passwds := g.Config().PasswdMap

	passwdSet = mapset.NewSetFromSlice(passwds) // 去重
	for _, host := range hostSet.ToSlice() {
		log.Printf("开始测试%s密码", host)
		go SSHSingle(host)
	}

	for i := 0; i < hostSet.Cardinality(); i++ {
		<-hostchan
	}

}
func SSHSingle(host interface{}) {
	defer func() {
		hostchan <- struct{}{}
	}()
	h := host.(string)
	for _, pass := range passwdSet.ToSlice() {
		ssh1 := sshtools.New_ssh(22, []string{h, "root", pass.(string)}...)
		//fmt.Println(ssh1)
		err := ssh1.Connect()
		if err == nil {
			HostPass.Put(h, pass)
			ssh1.Session.Close()
			//pass = pass.(string)
			return
		}
	}
	FailHosts.Add(host)
	return
}
