#gRPC服务发现及负载均衡实现

gRPC开源组件官方并未直接提供服务注册与发现的功能实现，但其设计文档已提供实现的思路，并在不同语言的gRPC代码API中已提供了命名解析和负载均衡接口供扩展。



##其基本实现原理：


1.服务启动后gRPC客户端向命名服务器发出名称解析请求，名称将解析为一个或多个IP地址，每个IP地址标示它是服务器地址还是负载均衡器地址，以及标示要使用那个客户端负载均衡策略或服务配置。

2.客户端实例化负载均衡策略，如果解析返回的地址是负载均衡器地址，则客户端将使用grpclb策略，否则客户端使用服务配置请求的负载均衡策略。

3.负载均衡策略为每个服务器地址创建一个子通道（channel）。

4.当有rpc请求时，负载均衡策略决定那个子通道即grpc服务器将接收请求，当可用服务器为空时客户端的请求将被阻塞。

## 运行方法
两个服务端
1. /home/go/GoDevEach/grpcx-master/example/grpcService1.go
2. /home/go/GoDevEach/grpcx-master/example/grpcService2.go
客户端
3. /home/go/GoDevEach/grpcx-master/example/grpcClient.go