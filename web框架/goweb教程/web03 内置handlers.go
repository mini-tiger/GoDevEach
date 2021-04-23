package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  web03 内置handlers
 * @Version: 1.0.0
 * @Date: 2021/4/20 下午12:12
 */

func helloworld3(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {

	//http.NotFoundHandler()  // 全局404

	// xxx 超时handler
	timeouthandler := http.TimeoutHandler(http.HandlerFunc(helloworld3), time.Duration(1*time.Second), "timeout handle")

	http.Handle("/timeouthandler", timeouthandler) // xxx 超时报错

	// xxx 跳转handler
	redirecthandler := http.RedirectHandler("/timeouthandler", 201) // 跳转url 在页面生成一个 超链接

	http.Handle("/redirecthandler", redirecthandler) // xxx 跳转url

	// xxx 去掉前缀
	//    /file/index03.html  访问  /static/index03.html 文件
	http.Handle("/", http.StripPrefix("/file", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8081", nil))

}
