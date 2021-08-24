package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  mysqlSearch
 * @Version: 1.0.0
 * @Date: 2021/4/21 下午2:06
 */

func registerHttpsTestControllers() {

	http.Handle("/serverpush", http.HandlerFunc(serverPushView))
}

func serverPushView(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println(writer.(http.Pusher))
	currdir, _ := os.Getwd()
	if pusher, ok := writer.(http.Pusher); ok {
		log.Println(" use https serverpush")

		pusher.Push(path.Join(currdir, "static/index.css"), &http.PushOptions{
			Header: http.Header{"Content-Type": []string{"text/css"}},
		})
	}
	//ResponseJsonSuccess(writer, request)
	//b:=make([]byte,1024)
	b, _ := ioutil.ReadFile(path.Join(currdir, "static/index03.html"))
	writer.Write(b)
}
