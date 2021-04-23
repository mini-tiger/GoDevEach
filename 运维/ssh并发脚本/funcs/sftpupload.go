package funcs

import (
	"fmt"
	"gitee.com/taojun319/tjtools/nmap"
	"gitee.com/taojun319/tjtools/sshtools"
	"log"
	"sync"
)

var wg *sync.WaitGroup = new(sync.WaitGroup)
var FailSftpHosts = make([]string, 0)
var FailRun = make([]string, 0)
var HostResult *nmap.SafeMap = nmap.NewSafeMap()

func UploadFileAndRun(sfile, dfile string) {

	for host, passwd := range HostPass.M {
		wg.Add(1)
		go BeginUpload(host, passwd.(string), sfile, dfile)
	}

	wg.Wait()
}

func BeginUpload(host, passwd, sFile, dFile string) {
	defer func() {
		wg.Done()
	}()
	if err := sshtools.SftpUpload(host, passwd, sFile, dFile); err == nil {
		log.Printf("sftp success upload host:%s begin run script %s", host, dFile)
		result := ScriptRunSub(host, passwd, fmt.Sprintf("bash %s", dFile), true, dFile)
		if result != "" {
			HostResult.Put(host, result)
		}

	} else {
		log.Printf("sftp Fail err:%s\n", err)
		FailSftpHosts = append(FailSftpHosts, host)
	}
}

func ScriptRunSub(h, pass, cmd string, shellFile bool, dfile string) (result string) {
	ssh1 := sshtools.New_ssh(22, []string{h, "root", pass}...)
	//fmt.Println(ssh1)
	err := ssh1.Connect()
	if err != nil {
		log.Printf("RunScript sshconnect err:%s\n ", err)
		FailRun = append(FailRun, h)
		return ""
	}

	if err, result = ssh1.TerminalRun(cmd, dfile, shellFile); err != nil {
		log.Printf("RunScript err:%s\n ", err)
		FailRun = append(FailRun, h)
		return ""
	}
	return
}
