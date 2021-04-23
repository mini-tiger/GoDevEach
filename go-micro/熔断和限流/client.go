package main

import (
	"breakerlimiter/message"
	"context"
	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"log"
	"time"
)

func main() {

	service := micro.NewService(
		micro.Name("student.client"),
		micro.WrapClient(hystrix.NewClientWrapper()), // xxx WrapClient 用于包装对外发布的请求 客户端包装
		// todo 熔断可以 不同服务 分别 设置参数
	)

	service.Init()
	hystrixGo.DefaultMaxConcurrent = 3
	// todo DefaultMaxConcurrent 的作用域是方法级的，与节点数无关。因此， 当某个服务的节点数被扩容一倍， 那么也有必要修改相应hystrix 限制，否则此扩容可能无法发挥全部效果

	hystrixGo.DefaultTimeout = 10 //熔断超时10 毫秒

	// more config      https://github.com/afex/hystrix-go/


	studentService := message.NewStudentServiceClient("student_service",service.Client())
	for i:=0;i<100;i++{
		go func(i int) {
			res, err := studentService.GetStudent(context.TODO(), &message.StudentRequest{Name: "davie"})
			if err != nil {
				log.Printf("ReqNum:%d,err:%s\n",i,err)
				return
			}
			log.Printf("ReqNum:%d,Name:%s,Classes:%s,Grade:%d\n",i,res.Name,res.Classes,res.Grade)
			//fmt.Println(res.Name)
			//fmt.Println(res.Classes)
			//fmt.Println(res.Grade)
		}(i)

	}


	time.Sleep(50 * time.Second)
}
