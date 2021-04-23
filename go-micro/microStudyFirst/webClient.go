package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/web"
	"log"
	"messageStudy/message"
	"net/http"
)

var Client message.StudentService

func main() {

	serviceWeb := web.NewService(
		web.Name("student.client"),
		web.Address(":8888"),
		web.StaticDir("html"), //测试接口 ，retful 不需要页面
	)

	serviceWeb.Init()

	srvClient := micro.NewService(
		micro.Name("student.client"),
	)

	Client=message.NewStudentService("student_service",srvClient.Client())


	// http://x:8888/student?name=davie
	serviceWeb.HandleFunc("/student",StuentWebService)

	serviceWeb.Run()
}

func StuentWebService(write http.ResponseWriter, request *http.Request)  {
	msg:=request.URL.Query().Get("name")
	log.Printf("请求参数:%v\n",msg)
	res, _ := Client.GetStudent(context.TODO(), &message.StudentRequest{Name: msg})
	fmt.Println(res)
	fmt.Fprint(write,res)
}