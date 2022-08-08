package main

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//mongo连接参数
const (
	user     = "cc"
	password = "cc"
	hosts    = "172.22.50.25:27021"
	mongoOpt = "replicaSet=rs0"
	auth     = "cmdb"
	timeout  = time.Duration(3000) * time.Millisecond
)

//mongo文档结构体
type student struct {
	Name   string `bson:"name"`
	Gender string `bson:"gender"`
	Age    int    `bson:"age"`
}

func main() {
	//设置连接参数
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?connectTimeoutMS=300000&authSource=%s", user, password, hosts, auth)
	opt := options.Client().ApplyURI(uri).SetSocketTimeout(timeout)

	//创建一个context上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	//获得一个mongo client
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		logs.Error("connect mongo failed, err:%s", err.Error())
		return
	}

	//ping一下mongo
	err = client.Ping(ctx, nil)
	if err != nil {
		logs.Error("ping mongo failed, err:%s", err.Error())
		return
	}

	//设定连接的数据库和集合
	database := auth
	collection := "student"

	var operations []mongo.WriteModel

	// Example using User struct
	userA := student{
		Name:   "Michael",
		Gender: "Male",
		Age:    21,
	}

	userB := student{
		Name:   "Michael11",
		Gender: "Male11",
		Age:    23,
	}

	operationA := mongo.NewUpdateOneModel()
	operationA.SetFilter(bson.M{"name": userA.Name})
	operationA.SetUpdate(bson.M{"$set": bson.M{"age": 22}})
	// Set Upsert flag option to turn the update operation to upsert
	operationA.SetUpsert(true)
	operations = append(operations, operationA)

	//Example using bson.M{}
	operationB := mongo.NewReplaceOneModel()
	operationB.SetFilter(bson.M{"name": userA.Name})
	operationB.SetReplacement(userB)
	operationB.SetUpsert(true)

	operations = append(operations, operationB)
	//
	//// Specify an option to turn the bulk insertion in order of operation
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)

	result, err := client.Database(database).Collection(collection).BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", result)
}
