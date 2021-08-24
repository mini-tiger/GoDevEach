package main

import (
	_ "github.com/go-sql-driver/mysql"
	"goweb/controllers"
	"goweb/middleware"
	_ "goweb/models"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/4/21 下午2:05
 */

func main() {

	// 注册路由
	controllers.RegisterControllers()

	//注册middleware
	middlewareHandle := &middleware.TimeoutMiddleware{ // 2s timeout
		Next: &middleware.BasicAuthMiddleware{
			Next: &middleware.LogMiddleware{}, // 顺序 Timeout -> Auth -> log
		},
	}

	server := http.Server{
		Addr:    ":8081",
		Handler: middlewareHandle,
	}

	currdir, _ := os.Getwd()

	log.Println("Server starting...")

	// pprof http://192.168.43.28:8080/debug/pprof/
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	server.ListenAndServe()

	// xxx  https
	/*
		1. go run /usr/local/go/src/crypto/tls/generate_cert.go -host localhost  当前路径执行  生成证书
		2. serverpush 提前发送css等需要加载的文件  使用浏览器 测试 https://192.168.43.28:8081/serverpush
		3. 打开 file server
	*/

	// xxx web file server 见  web-06 文件服务response.go
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(currdir, "static")))))

	//_ = server.ListenAndServeTLS(path.Join(currdir, "cert.pem"), path.Join(currdir, "key.pem"))

}
