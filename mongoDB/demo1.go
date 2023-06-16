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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
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

var uri = "mongodb://root:cc@172.22.50.25:32082,172.22.50.25:32083,172.22.50.25:32084/?connectTimeoutMS=300000&authSource=admin"

//var uri string = "mongodb://cc:cc@172.22.50.25:27117,172.22.50.25:27118/?authMechanism=SCRAM-SHA-256&authSource=cmdb"

func main() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(uri)
	var mp uint64 = 3000
	var app = "migrate"
	//var rs = "rs0"
	conOpt := options.ClientOptions{
		MaxPoolSize: &mp,
		AppName:     &app,
	}

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions, &conOpt)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	for _, host := range clientOptions.Hosts {
		// 创建一个新的连接到当前mongos节点
		mongosClientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", host))
		mongosClient, err := mongo.Connect(context.TODO(), mongosClientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// 获取admin数据库的isMaster信息
		adminDB := mongosClient.Database("admin")
		result := adminDB.RunCommand(context.TODO(), map[string]interface{}{
			"hello": 1,
		})

		// 解析并获取mongos节点状态
		var output map[string]interface{}
		if err := result.Decode(&output); err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
		// 获取isMaster命令的返回字段中的"ismaster"值，判断节点是否为mongos节点
		//ismaster, ok := output["ismaster"].(bool)
		//if !ok {
		//	log.Fatal("Unable to determine the status of mongos node")
		//}
		//
		//if ismaster {
		//	fmt.Printf("Mongos node %s is online\n", host)
		//} else {
		//	fmt.Printf("Mongos node %s is offline\n", host)
		//}

		// 关闭当前mongos节点的连接
		if err := mongosClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}

	// 关闭MongoDB连接
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}

}

func insertDemo(client *mongo.Client) {
	//不用提前创建db: cmdb
	collection := client.Database("cmdb").Collection("cc_ObjectBase")
	for i := 10; i < 9999900; i++ {
		_, err := collection.InsertOne(
			context.Background(), // 上下文参数
			bson.M{ // 使用bson.D定义一个JSON文档
				"bk_inst_id":   fmt.Sprintf("cmdb111__%d", i),
				"bk_inst_name": "aaa",
				"bk_obj_id":    "region",
			})
		fmt.Println(i, err)
		time.Sleep(500 * time.Millisecond)
	}

}
