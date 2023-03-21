package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"selfupdate/g"
	"selfupdate/services"
	"sync"
	"syscall"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: 程序自动更新
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/9/14 下午12:02
 */

func cancel(ctx context.Context, wg *sync.WaitGroup, stopInter services.RegSrvInterface) {

	//services.StopWeb(ctx, wg)
	stopInter.StopSrv(ctx, wg)
}

/*
 xxx 1.接收 或者 获取， 当前程序版本 和 最新程序地址,如果有新版本 下载

 xxx 2.关闭程序占用的端口

 xxx 3.启动新程序
 xxx 4. 删除 旧程序文件
*/
func main() {
	exitChan := make(chan os.Signal, 1)
	// 使用signal.Notify来捕捉退出信号，一般是使用term int信号来作为关闭信号，可以根据自己需要选择
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGHUP) // https://colobu.com/2015/10/09/Linux-Signals/  kill -2  kill -1

	ctxTop, _ := context.WithCancel(context.Background())
	//附加值
	valueCtx1 := context.WithValue(ctxTop, "port", ":8081")

	regWebInter := services.RegWebSrv()
	go func() {
		regWebInter.StartSrv(valueCtx1)
	}()

	go func() {

		select {
		// 等待退出信号
		case <-exitChan:
			log.Println("get exit signal")
			ctx, _ := context.WithTimeout(valueCtx1, 2*time.Second)
			// 告知业务处理函数该退出了
			wg := sync.WaitGroup{}
			wg.Add(1)
			cancel(ctx, &wg, regWebInter)
			// 等待业务处理函数全都退出
			wg.Wait()
			g.ExitAppChan <- struct{}{}
			return
		}

	}()

	<-g.ExitAppChan // old app exit

	log.Println("exit main success")
	log.Println("begin start process New App")
	go func() {

		ep := "/data/work/go/GoDevEach/程序自更新/main1"
		dir := filepath.Dir(ep)
		fullpath, err := filepath.Rel(dir, ep)
		if err != nil {
			log.Fatalf("filepath.Rel: %v\n", err)
		}
		logf, err := os.OpenFile("/tmp/22", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
		if err != nil {
			panic(err)
		}
		//bashpath, err := exec.LookPath("nohup")
		//if err != nil {
		//	panic(err)
		//}
		cmd := exec.Command(fmt.Sprintf("./%s", fullpath))
		cmd.Dir = dir
		cmd.Path = fullpath
		cmd.Env = os.Environ()
		cmd.Stdout = logf
		cmd.Stderr = logf
		//cmd.ExtraFiles = []*os.File{logf}

		//cmd.Args = []string{"> ", "/tmp/22 2>&1 ", "&"}

		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // xxx 拥有自己的进程组,子进程 独立于 父进程
		go func() {
			err = cmd.Run() // 阻塞
			if err != nil {
				log.Println("err:", err)
			}
		}()

		//cmd.Wait()
		time.Sleep(2 * time.Second)
		//fmt.Printf(" current process pid: %v\n", cmd.ProcessState.Pid())
		g.NewAppFinishChan <- struct{}{}

	}()
	<-g.NewAppFinishChan
	//time.Sleep(2 * time.Second)
	os.Exit(0)

	//time.Sleep(3 * time.Second)

	//os.Exit(0)
	//select {}

}
