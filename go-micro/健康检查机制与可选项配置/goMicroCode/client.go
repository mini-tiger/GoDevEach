package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"go-micro/健康检查机制与可选项配置/goMicroCode/message"
	"time"
)

func main() {

	consulReg :=consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))

	service := micro.NewService(
		micro.Name("student.client"),
		micro.Registry(consulReg),
	)

	service.Init()

	studentService := message.NewStudentServiceClient("student_service", service.Client())

	res, err := studentService.GetStudent(context.TODO(), &message.StudentRequest{Name: "davie"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Name)
	fmt.Println(res.Classes)
	fmt.Println(res.Grade)
	time.Sleep(50 * time.Second)
}
