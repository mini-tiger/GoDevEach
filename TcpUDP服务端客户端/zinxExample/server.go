package main

import (
	"fmt"
	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

// xxx https://github.com/aceld/zinx
// xxx https://aceld.gitbooks.io/zinx/content/san-3001-zinx-kuang-jia-ji-chu-lu-you-mo-kuai/34-zinx-v03dai-ma-shi-xian.html
//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	//先读取客户端的数据
	fmt.Println("Recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//xxx 再向client回写ping...ping...ping
	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type PingRouter1 struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter1) Handle(request ziface.IRequest) {
	//先读取客户端的数据
	fmt.Println("Recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//xxx 再向client回写ping...ping...ping
	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server句柄

	utils.GlobalObject = &utils.GlobalObj{ConfFilePath: "/home/go/GoDevEach/TcpUDP服务端客户端/zinxExample/cfg.json"}
	utils.GlobalObject.Reload()

	s := znet.NewServer()

	//2 配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &PingRouter{})
	//3 开启服务
	s.Serve()
}
