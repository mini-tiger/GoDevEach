package main

import (
	"context"

	"log"
	"strconv"


	"google.golang.org/grpc"

	"github.com/yakaa/grpcx"
	"github.com/yakaa/grpcx/config"
	"github.com/yakaa/grpcx/example/proto"
)

var Endpoints = []string{"192.168.43.26:2379"}
var Schema = "www.tao.com"

func main() {

	for i := 2; i < 4; i++ {
		go func(j int) {
			Server(j)
		}(i)
	}
	select {}
	//Client()

}

func Server(count int) {
	conf := &config.ServiceConf{
		EtcdAuth:      config.EtcdAuth{},
		Schema:        Schema,
		ServerName:    "knowing",
		Endpoints:     Endpoints,
		ServerAddress: "192.168.43.26:2000" + strconv.Itoa(count),
	}
	demo := &RegionHandlerServer{ServerAddress: conf.ServerAddress}
	rpcServer, err := grpcx.MustNewGrpcxServer(conf, func(server *grpc.Server) {
		proto.RegisterRegionHandlerServer(server, demo)
	})
	if err != nil {
		panic(err)
	}
	log.Fatal(rpcServer.Run()) //xxx 监听端口,服务注册
}

type RegionHandlerServer struct {
	ServerAddress string
}

func (s *RegionHandlerServer) GetListenAudio(ctx context.Context, r *proto.FindRequest) (*proto.HasOptionResponse, error) {

	has := []*proto.HasOption(nil)
	for _, t := range r.Tokens {

		has = append(has, &proto.HasOption{Token: t + " FROM " + s.ServerAddress, Listen: 1})
	}
	res := &proto.HasOptionResponse{
		Items: has,
	}
	log.Printf("%s 接收到的请求%+v\n",s.ServerAddress,r)
	return res, nil
}

//func Client() {
//	conf := &config.ClientConf{
//		EtcdAuth:  config.EtcdAuth{},
//		Target:    Schema + ":///knowing",
//		Endpoints: Endpoints,
//		WithBlock: false,
//	}
//
//	r, err := grpcx.MustNewGrpcxClient(conf)
//	//创建 Etcd key的监控
//	// Builder接口在发起rpc请求的时候会调用Build方法。etcd Resolver的Build方法首先创建一条到etcd服务端的连接。然后启动一个goroutine watch相应的key上是否有变更，如果有，根据不同的event进行不同的处理
//
//	if err != nil {
//		panic(err)
//	}
//	conn, err := r.GetConnection() // 负载均衡
//	if err != nil {
//		panic(err)
//	}
//	regionHandlerClient := proto.NewRegionHandlerClient(conn)
//
//	var in = &proto.FindRequest{Tokens: []string{"a_" + strconv.FormatInt(time.Now().Unix(), 10)}}
//	var out *proto.HasOptionResponse
//	for {
//		out, err = regionHandlerClient.GetListenAudio(
//			context.Background(),
//			in,
//		)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("Client 发送:%+v, 接收:%+v\n", in, out)
//		time.Sleep(1 * time.Second)
//	}
//}
