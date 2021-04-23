package main

import (
	"context"
	"fmt"
	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/web"
	"go-micro/RESTfu_API设计标准和实践/MicroApiDemo/proto"
	"log"
)

type Student struct {
}

var (
	cli proto.StudentService
)

func (s *Student) GetStudent(req *restful.Request, rsp *restful.Response) {

	name := req.PathParameter("name")
	fmt.Println(name)
	response, err := cli.GetStudent(context.TODO(), &proto.Request{
		Name: name,
	})

	if err != nil {
		fmt.Println(err.Error())
		rsp.WriteError(500, err)
	}

	rsp.WriteEntity(response)
}

func main() {

	/*
	访问  IP:port/student/davie
	通过API mdns方式 获取server的数据 并返回
	*/
	web.DefaultAddress="0.0.0.0:8083"
	service := web.NewService(
		web.Name("go.micro.api.student"),
	)

	service.Init()

	cli = proto.NewStudentService("go.micro.srv.student", client.DefaultClient)

	student := new(Student)
	ws := new(restful.WebService)
	ws.Path("/student")
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON) //请求头必须要符合任意一种，默认是xml
	ws.Produces(restful.MIME_JSON, restful.MIME_XML) //返回任意一种格式，默认JSON

	ws.Route(ws.GET("/{name}").To(student.GetStudent))

	wc := restful.NewContainer()
	wc.Add(ws)

	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
