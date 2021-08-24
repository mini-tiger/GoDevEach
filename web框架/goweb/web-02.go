package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Description: Customize
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/4/20 上午9:41
 */

type homehandle struct {
}

func (c *homehandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "home")
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	var home = &homehandle{} // 自定义实现 ServeHTTP
	http.Handle("/home", home)

	http.HandleFunc("/helloworld", helloworld)
	http.Handle("/helloworld1", http.HandlerFunc(helloworld)) // 与helloworld效果一样,包装方法

	log.Fatal(http.ListenAndServe(":8081", nil))

}
