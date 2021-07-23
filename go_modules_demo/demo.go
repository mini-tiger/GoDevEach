package main

import (
	"github.com/gin-gonic/gin"
	"runtime"
)

// xxx go mod 命令 https://colobu.com/2018/08/27/learn-go-module/

//1. xxx cd go_modules_demo && go mod init &&  go mod why # 自检
//2. 查找可用版本go list -m -versions github.com/gin-gonic/gin 或者 版本写 latest
//3. xxx go mod download github.com/gin-gonic/gin  生成 go.sum 或者 go mod tidy

func main() {
	// xxx 默认已经连接了 Logger and Recovery 中间件
	//r := gin.Default()

	// windows 无法显示日志颜色
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	} else {
		gin.ForceConsoleColor()
	}

	gin.SetMode("debug")
	// 创建一个默认的没有任何中间件的路由
	r := gin.New()

	r.Run(":8001")
}
