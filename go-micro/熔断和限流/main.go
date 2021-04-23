package main

import (
	"breakerlimiter/message"
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro/v2"
	limiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"

	"log"
)

//学生服务管理实现
type StudentManager struct {
}

//获取学生信息的服务接口实现
func (sm *StudentManager) GetStudent(ctx context.Context, request *message.StudentRequest, response *message.Student) error {
	log.Printf("请求参数:%+v\n", request)
	//tom
	studentMap := map[string]message.Student{
		"davie":  message.Student{Name: "davie", Classes: "软件工程专业", Grade: 80},
		"steven": message.Student{Name: "steven", Classes: "计算机科学与技术", Grade: 90},
		"tony":   message.Student{Name: "tony", Classes: "计算机网络工程", Grade: 85},
		"jack":   message.Student{Name: "jack", Classes: "工商管理", Grade: 96},
	}

	if request.Name == "" {
		return errors.New(" 请求参数错误,请重新请求。")
	}

	//获取对应的student
	student := studentMap[request.Name]

	if student.Name != "" {
		fmt.Println(student.Name, student.Classes, student.Grade)
		*response = student
		return nil
	}

	return errors.New(" 未查询当相关学生信息 ")
}

func main() {

	//创建一个新的服务对象实例
	service := micro.NewService(
		micro.Name("student_service"),
		micro.Version("v1.0.0"),
		micro.WrapHandler(limiter.NewHandlerWrapper(2)), // todo 限流处理器加载,并发2
		// WrapHandler 用于包装外界发来的请求，服务端包装
	)

	//服务初始化
	service.Init()



	//注册
	message.RegisterStudentServiceHandler(service.Server(), new(StudentManager))

	//运行
	err := service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
