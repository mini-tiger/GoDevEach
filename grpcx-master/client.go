package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	user "grpcx-master/proto"
)

const (
	address = "localhost:50051"
)

var userClient user.UserClient
var ctx context.Context

func main() {
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	userClient = user.NewUserClient(conn)

	// 设定请求超时时间 3s
	//ctx, _ = context.WithTimeout(context.Background(), time.Second*3)
	ctx = context.Background()
	//defer cancel()
	for i := 0; ; i++ {
		send()
		time.Sleep(5 * time.Second)
	}
}

func send() {
	// UserIndex 请求
	userIndexReponse, err := userClient.UserIndex(ctx, &user.UserIndexRequest{Page: 1, PageSize: 12})
	if err != nil {
		log.Printf("user index could not greet: %v", err)
	}

	if 0 == userIndexReponse.Err {
		log.Printf("user index success: %s", userIndexReponse.Msg)
		// 包含 UserEntity 的数组列表
		userEntityList := userIndexReponse.Data
		for _, row := range userEntityList {
			fmt.Println(row.Name, row.Age)
		}
	} else {
		log.Printf("user index error: %d", userIndexReponse.Err)
	}

	// UserView 请求
	userViewResponse, err := userClient.UserView(ctx, &user.UserViewRequest{Uid: 1})
	if err != nil {
		log.Printf("user view could not greet: %v", err)
	}

	if 0 == userViewResponse.Err {
		log.Printf("user view success: %s", userViewResponse.Msg)
		userEntity := userViewResponse.Data
		fmt.Println(userEntity.Name, userEntity.Age)
	} else {
		log.Printf("user view error: %d", userViewResponse.Err)
	}

	// UserPost 请求
	userPostReponse, err := userClient.UserPost(ctx, &user.UserPostRequest{Name: "big_cat", Password: "123456", Age: 29})
	if err != nil {
		log.Printf("user post could not greet: %v", err)
	}

	if 0 == userPostReponse.Err {
		log.Printf("user post success: %s", userPostReponse.Msg)
	} else {
		log.Printf("user post error: %d", userPostReponse.Err)
	}

	// UserDelete 请求
	userDeleteReponse, err := userClient.UserDelete(ctx, &user.UserDeleteRequest{Uid: 1})
	if err != nil {
		log.Printf("user delete could not greet: %v", err)
	}

	if 0 == userDeleteReponse.Err {
		log.Printf("user delete success: %s", userDeleteReponse.Msg)
	} else {
		log.Printf("user delete error: %d", userDeleteReponse.Err)
	}
}
