package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"net/http"
)

const (
	consulAddress = "192.168.43.111:8500"
	svcIp         = "192.168.43.111"
	svcPort       = 8080
	checkSuffix   = "jenkins"
)

var svcid string = "337"
var svcname string = fmt.Sprintf("service%s", svcid)

func consulRegister(client *consulapi.Client) (err error) {

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = svcid
	registration.Name = svcname
	registration.Port = svcPort
	registration.Tags = []string{"testService"}
	registration.Address = svcIp

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d/%s/", registration.Address, registration.Port, checkSuffix)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
	return
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you are visiting health check api"))
}

// 从consul中发现服务
func ConsulFindServer(client *consulapi.Client) {
	//// 创建连接consul服务配置
	//config := consulapi.DefaultConfig()
	//config.Address = consulAddress
	//client, err := consulapi.NewClient(config)
	//if err != nil {
	//	fmt.Println("consul client error : ", err)
	//}

	// 获取指定service
	service, _, err := client.Agent().Service("hello", nil)
	if err == nil {
		fmt.Println(service.Address)
		fmt.Println(service.Port)
		//fmt.Printf("svc: %#v\n",service)
	}

	//xxx 只获取健康的service
	serviceHealthy0, _, err := client.Health().Checks("hello", nil)
	if err == nil {
		fmt.Printf("check: %s\n", serviceHealthy0.AggregatedStatus()) // 取check状态最差的
	}

	serviceHealthy, _, err := client.Health().Service("hello", "", false, nil)
	fmt.Println("serviceHealthy: ", len(serviceHealthy))
	if err == nil && len(serviceHealthy) > 0 {
		fmt.Printf("Health svc:%+v\n", serviceHealthy[0].Service.Address)
		fmt.Printf("Health svc:%+v\n", serviceHealthy[1].Service.Address)
	} else {
		fmt.Println("Health svc nil")
	}

}

func main() {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return
	}

	//注册
	//if err:=consulRegister(client);err!=nil{
	//	log.Fatalln(err)
	//}
	//发现
	ConsulFindServer(client)

	//定义一个http接口
	//http.HandleFunc("/", Handler)
	//err := http.ListenAndServe("0.0.0.0:81", nil)
	//if err != nil {
	//	fmt.Println("error: ", err.Error())
	//}
	//删除 svc
	//Dregister(client)

}
func Dregister(client *consulapi.Client) {

	//fmt.Println("test begin .")
	//config := consulapi.DefaultConfig()
	////config.Address = "localhost"
	//fmt.Println("defautl config : ", config)
	//client, err := consulapi.NewClient(config)
	//if err != nil {
	//	log.Fatal("consul client error : ", err)
	//}

	err := client.Agent().ServiceDeregister(svcid)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

}
