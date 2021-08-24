package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  web-0405,request参数
 * @Version: 1.0.0
 * @Date: 2021/4/20 下午4:38
 */

func GetReqParams(writer http.ResponseWriter, request *http.Request) {
	log.Println("request header:", request.Header) // map[string][]string
	log.Println(request.Header["User-Agent"])
	var rs strings.Builder
	rs.WriteString("Request Header:\n")
	for key, value := range request.Header {
		rs.WriteString(fmt.Sprintf("\nkey:%s , value:%s\n", key, value))
	}
	fmt.Fprint(writer, rs.String())
}

func GetParams(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	// 第一种方式
	// id := query["id"][0]

	// 第二种方式
	id := query.Get("id")

	fmt.Printf("GET: id=%s\n", id)

	fmt.Fprintf(writer, `{"code":0}`)
}

func PostBodyParams(writer http.ResponseWriter, request *http.Request) {

	//方法一 根据请求body创建一个json解析器实例
	//decoder := json.NewDecoder(request.Body)

	// 用于存放参数key=value数据
	//var params map[string]string

	// 解析参数 存入map
	//decoder.Decode(&params)

	//fmt.Printf("POST json: username=%s, password=%s\n", params["username"], params["password"])
	//fmt.Fprintf(writer, `{"code":0}`)

	//方法二 根据请求body创建struct
	type reqParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var reqP *reqParams = &reqParams{}

	err := json.NewDecoder(request.Body).Decode(reqP)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Set response header
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reqP)
	if err != nil {
		fmt.Fprintf(writer, `{"code":400}`)
	}

}

func PosturlencodedParams(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	// 第一种方式
	// username := request.Form["username"][0]
	// password := request.Form["password"][0]

	// 第二种方式
	username := request.Form.Get("username")
	password := request.Form.Get("password")

	// 第三种方式
	fmt.Println(request.FormValue("username")) //get key value
	fmt.Println(request.PostFormValue("username"))

	fmt.Printf("POST form-urlencoded: username=%s, password=%s\n", username, password)

	fmt.Fprintf(writer, `{"code":0}`)
}

func PosturlformdataParams(writer http.ResponseWriter, request *http.Request) {

	fmt.Println(request.FormValue("username")) //get key value
	fmt.Println(request.PostFormValue("password"))

	//fmt.Printf("POST form-urlencoded: username=%s, password=%s\n", username, password)

	fmt.Fprintf(writer, `{"code":0}`)
}

func PosturlfileParams(writer http.ResponseWriter, request *http.Request) {

	// 第一种方式
	//request.ParseMultipartForm(1024)
	//fileHeader:=request.MultipartForm.File["uploadfile"][0]
	//formFile,_:=fileHeader.Open()

	// 第二种方式(xxx 适合只传文件)
	// 根据字段名获取表单文件
	formFile, _, err := request.FormFile("uploadfile")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}

	defer formFile.Close()

	var buf []byte = make([]byte, 1024)
	by := bytes.NewBuffer(buf)
	formFile.Read(by.Bytes())
	fmt.Fprintln(writer, by.String())

}

func main() {
	// 获取 request header
	http.Handle("/GetReqParams", http.HandlerFunc(GetReqParams))

	// 获取 get params
	http.Handle("/GetParams", http.HandlerFunc(GetParams))

	// 获取 post params 处理application/json类型的POST请求 xxx 并json 返回
	http.Handle("/PostBodyParams", http.HandlerFunc(PostBodyParams))

	// 获取 post params 处理application/x-www-form-urlencoded类型的POST请求
	http.Handle("/PosturlencodedParams", http.HandlerFunc(PosturlencodedParams))

	// 获取 post params 处理form-data类型的POST请求
	http.Handle("/PosturlformdataParams", http.HandlerFunc(PosturlformdataParams))

	// 获取 上传文件 处理form-data类型的POST请求
	http.Handle("/PosturlfileParams", http.HandlerFunc(PosturlfileParams))

	log.Fatal(http.ListenAndServe(":8081", nil))

}
