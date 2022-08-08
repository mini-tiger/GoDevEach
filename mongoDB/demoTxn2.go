package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
)

const (
	user     = "cc"
	password = "cc"
	hosts    = "172.22.50.25:27021"
	mongoOpt = "replicaSet=rs0"
	auth     = "cmdb"
	timeout  = time.Duration(3000) * time.Millisecond
)

/*
xxx 注意： 如果数据库中不存在有name为"无名小角色"的记录，上述操作都会失败，
xxx 尽管在事务中先插入了一条name相同的记录，但事务是一个整体，查找的是事务开始之前的记录，因此查找会失败从而导致插入操作同样会失败。
*/
func main() {
	//设置连接参数
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?connectTimeoutMS=300000&authSource=%s", user, password, hosts, auth)
	opt := options.Client().ApplyURI(uri).SetSocketTimeout(timeout)

	//创建一个context上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//if err != nil {
	//	log.Fatal(err)
	//}
	// collection
	client, err := mongo.Connect(ctx, opt)
	// collection
	collection := client.Database(auth).Collection("person")

	// session 读取策略
	sessOpts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	session, err := client.StartSession(sessOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer session.EndSession(context.TODO())

	// transaction 读取优先级
	transacOpts := options.Transaction().SetReadPreference(readpref.Primary())
	// 插入一条记录、查找一条记录在同一个事务中
	result, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		// insert one
		insertOneResult, err := collection.InsertOne(sessionCtx, bson.M{"name": "无名小角色", "level": 5})
		if err != nil {
			//sessionCtx.AbortTransaction(ctx)
			log.Fatal(err)
		}

		fmt.Println("inserted id:", insertOneResult.InsertedID)

		// find one
		var result struct {
			Name  string `bson:"name,omitempty"`
			Level int    `bson:"level,omitempty"`
		}
		singleResult := collection.FindOne(sessionCtx, bson.M{"name": "无名小角色"})
		if err = singleResult.Decode(&result); err != nil {

			return nil, err
		}

		return result, err

	}, transacOpts)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("find one result: %+v \n", result)
}
