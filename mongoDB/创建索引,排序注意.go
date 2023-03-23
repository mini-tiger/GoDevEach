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

//var uri = "mongodb://root:cc@172.22.50.25:32082,172.22.50.25:32083,172.22.50.25:32084/?connectTimeoutMS=300000&authSource=admin"

var uri string = "mongodb://cc:cc@172.22.50.25:27117,172.22.50.25:27118/?authMechanism=SCRAM-SHA-256&authSource=cmdb"

func main() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(uri)
	var mp uint64 = 3000
	var app = "migrate"
	//var rs = "rs0"
	conOpt := options.ClientOptions{
		MaxPoolSize: &mp,
		//MaxPoolSize:     3000,
		//MinPoolSize:     100,
		//ConnectTimeout:  &timeout,
		//SocketTimeout:   &socketTimeout,
		//ReplicaSet: &rs,
		//RetryWrites:     false,
		//MaxConnIdleTime: "1500000000000"
		AppName: &app,
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

	db := client.Database("admin")
	buildInfoCmd := bson.D{bson.E{Key: "buildInfo", Value: 1}}
	var buildInfoDoc bson.M
	if err := db.RunCommand(context.Background(), buildInfoCmd).Decode(&buildInfoDoc); err != nil {
		log.Printf("Failed to run buildInfo command: %v", err)
		return
	}
	log.Println("Database version:", buildInfoDoc["version"])

	// xxx indexes
	coll := client.Database("cmdb").Collection("cc_TopoGraphics")

	//common.BKTableNameTopoGraphics: {
	//	types.Index{Name: "", Keys: map[string]int32{"scope_type": 1, "scope_id": 1, "node_type": 1, common.BKObjIDField: 1, common.BKInstIDField: 1}, Background: true, Unique: true},
	//},
	//
	var b bool = true
	opt := &options.IndexOptions{}
	opt.Unique = &b
	opt.Background = &b
	//var be1 bson.E = bson.E{Key: "scope_type", Value: 1}
	var De []bson.E = bson.D{
		bson.E{Key: "scope_type", Value: 1},
		bson.E{Key: "scope_id", Value: 1},
		bson.E{Key: "node_type", Value: 1},
		bson.E{Key: "bk_obj_id", Value: 1},
		bson.E{Key: "bk_inst_id", Value: 1},
	}
	//{
	//	"scope_type", "scope_id", "node_type", "bk_obj_id", "bk_inst_id":1
	//},
	bb, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    De,
		Options: opt,
	})
	fmt.Println(bb, err)
	// xxx insert
	//insertDemo(client)
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
