package main

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"

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

/*

 */
func main() {
	//设置连接参数
	uri := fmt.Sprintf("mongodb://root:abc123@172.22.50.25:32082,172.22.50.25:32083,172.22.50.25:32084/?connectTimeoutMS=300000&authSource=admin")
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
	//构造插入的数据
	students := []interface{}{
		student{
			Name:   "Michael",
			Gender: "Male",
			Age:    21,
		},
		student{
			Name:   "Alice",
			Gender: "Female",
			Age:    19,
		},
	}

	//在会话中使用mongo
	if err = client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		//开启事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}
		//最后的时候，关闭会话
		defer sessionContext.EndSession(ctx)

		//插入数据
		if _, err := client.Database(database).Collection(collection).InsertMany(ctx, students); err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				//回滚事务
				logs.Error("mongo transaction rollback failed, %s", err.Error())
				return err
			}
			return err
		}
		//提交事务
		return sessionContext.CommitTransaction(ctx)
	}); err != nil {
		logs.Error("insert failed, err:%s", err.Error())
	}
}
