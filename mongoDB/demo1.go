/**
 * @Author: Tao Jun
 * @Since: 2022/7/16
 * @Desc: demo1.go
**/
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var uri = "mongodb://cc:cc@172.22.50.25:27021,172.22.50.25:27022,172.22.50.25:27023/?connectTimeoutMS=300000&authSource=cmdb"
var uri = "mongodb://root:abc123@172.22.50.25:32087/?connectTimeoutMS=300000&authSource=admin"

//var uri = "mongodb://cc:cc@172.22.50.25:40017/?connectTimeoutMS=300000&authSource=cmdb"

func main() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(uri)

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}
