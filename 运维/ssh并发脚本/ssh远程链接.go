package main

import (
	"fmt"
	"log"

	"path"
	"path/filepath"
	"runtime"
	"sshdemo/funcs"
	"sshdemo/g"
)

func main() {
	//解析config
	_, dir, _, _ := runtime.Caller(0)
	currDir := filepath.Dir(dir)

	g.ParseConfig(filepath.Join(currDir, "config.json"))

	funcs.SSHRun() // 测试密码
	log.Printf("密码错误的主机有:%v\n", funcs.FailHosts.ToSlice())
	log.Printf("正确的主机有:%v\n", funcs.HostPass.Keys())
	if funcs.HostPass.IsEmpty() {
		log.Fatalf("没有成功的主机")
	}

	funcs.UploadFileAndRun(filepath.Join(currDir, "1.sh"), path.Join("/tmp/", "1.sh")) // 源脚本， 目标文件,unix 使用path.join
	log.Printf("上传文件错误的主机有:%v\n", funcs.FailSftpHosts)
	log.Printf("执行脚本错误的主机有:%v\n", funcs.FailRun)
	for host, result := range funcs.HostResult.M {
		fmt.Printf("主机:%s ,运行结果:\n%s\n", host, result)
	}
	//ssh.Session.Stdout=os.Stdout
	//ssh.Session.Stderr=os.Stderr
	//ssh.Session.Run("touch /root/1")
	//ssh.Session.Run("ls /; ls /tmp")
	//ssh.close_session() //todo session一次运行一次run

	//terminal_run(ssh.Session)
	//ssh.close_session()
	//select {}
}
