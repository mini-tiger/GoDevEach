package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
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

func FatherServer(ctx context.Context, wg *sync.WaitGroup, num int) {
	defer wg.Done()
	fatherWg := &sync.WaitGroup{}
	i := 0
	for {
		select {
		case <-ctx.Done():
			log.Printf("FatherServer %d wait exit\n", num)
			fatherWg.Wait()
			log.Printf("FatherServer %d exit success\n", num)
			return
		case <-time.After(time.Second * 5):
			i++
			fatherWg.Add(1)
			go ChildServer(fatherWg, fmt.Sprintf("FatherServer num:%d-%d", num, i))
		}
	}
}

func ChildServer(wg *sync.WaitGroup, args string) {
	defer wg.Done()
	fmt.Printf("ChildServer will process Form %s\n", args)
	ExecuteCmd(args)
	fmt.Printf("ChildServer process success%s\n", args)
}

func ExecuteCmd(args string) {
	cmd := exec.Command("/bin/bash", "/home/go/GoDevEach/程序自更新/hello.sh", args)
	// xxx 拥有自己的进程组,子进程 独立于 父进程
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("err:", err)
	}

}

// xxx https://segmentfault.com/a/1190000039805354

func main() {
	exitChan := make(chan os.Signal, 1)
	// 使用signal.Notify来捕捉退出信号，一般是使用term int信号来作为关闭信号，可以根据自己需要选择
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	// 准备好wg和ctx来记录服务执行逻辑状态
	//wg := &sync.WaitGroup{}
	////  cancel用于取消，ctx用于通知
	//ctx, cancel := context.WithCancel(context.Background())
	//// 开启业务逻辑
	//wg.Add(2)
	//go FatherServer(ctx, wg,1)
	//go FatherServer(ctx, wg,2)

	// xxx 1.接收 或者 获取， 当前程序版本 和 最新程序地址,如果有新版本 下载

	// xxx 2.关闭程序占用的端口

	// xxx 3.启动新程序

	go func() {
		cmd := exec.Command("/bin/bash", "/home/go/GoDevEach/程序自更新/hello1.sh")
		// xxx 拥有自己的进程组,子进程 独立于 父进程
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		_, err := cmd.CombinedOutput()

		if err != nil {
			log.Println("err:", err)
		}
	}()

	time.Sleep(3 * time.Second)
	// xxx 4. 删除 旧程序文件

	os.Exit(0)

	//for {
	//	select{
	//	// 等待退出信号
	//	case <-exitChan:
	//		log.Println("get exit signal")
	//		// 告知业务处理函数该退出了
	//		cancel()
	//		// 等待业务处理函数全都退出
	//		wg.Wait()
	//		log.Println("exit main success")
	//		return
	//	}
	//}
}
