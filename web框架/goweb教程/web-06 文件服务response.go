package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

func FileServer(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()

	// 第一种方式
	// id := query["id"][0]

	// 第二种方式
	file := query.Get("file")

	fmt.Printf("GET: file=%s\n", file)

	http.Redirect(writer, request, "/filedir/"+file, 201) //页面生成跳转 href

}

func main() {

	// 内置response
	// http.NotFound(writer,request)   404
	// 	http.ServeFile(writer,request,"./static/index03.html")	   读取文件返回内容
	//
	// 	http.ServeContent(writer,request,f.Name(),fileinfo.ModTime(),f) // 数据流返回 seek
	http.Handle("/PostJsonParams", http.HandlerFunc(FileServer))
	currentdir, _ := os.Getwd()

	// 访问http://192.168.43.28:8081/static/index03.html   打开  index03.html
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(currentdir, "static")))))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
