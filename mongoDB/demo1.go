/**
 * @Author: Tao Jun
 * @Since: 2022/7/16
 * @Desc: demo1.go
**/
package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// xxx helm k8s deploy

/*
xxx 副本 集连接方法
1. k8s nodeport 添加  directConnection=true 性能 下降 但可以 nodeport 端口转发
2. kubectl -n mongo10 port-forward service/mongo-mongodb-headless --address='0.0.0.0' 28015:27017 连接 28015

//var uri = "mongodb://root:abc123@172.22.50.25:28015/?connectTimeoutMS=300000&authSource=admin&directConnection=true"
*/
//var uri = "mongodb://root:pass2022@172.22.50.25:32091/?connectTimeoutMS=300000&authSource=admin&directConnection=true"

/*
xxx 分片 连接方法
//var uri = "mongodb://root:pass2022@172.22.50.25:32092,172.22.50.25:32093,172.22.50.25:32094/?connectTimeoutMS=300000&authSource=admin"
*/

var uri = "mongodb://root:pass2022@172.22.50.25:32092,172.22.50.25:32093,172.22.50.25:32094/?connectTimeoutMS=300000&authSource=admin"

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

	db := client.Database("admin")
	buildInfoCmd := bson.D{bson.E{Key: "buildInfo", Value: 1}}
	var buildInfoDoc bson.M
	if err := db.RunCommand(context.Background(), buildInfoCmd).Decode(&buildInfoDoc); err != nil {
		log.Printf("Failed to run buildInfo command: %v", err)
		return
	}
	log.Println("Database version:", buildInfoDoc["version"])
}
