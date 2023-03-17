package main

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

//mongo连接参数
const (
	//user     = "cc"
	//password = "cc"
	//hosts    = "172.22.50.25:27021"
	//mongoOpt = "replicaSet=rs0"
	auth    = "cmdb"
	timeout = time.Duration(30) * time.Second
)

//mongo文档结构体
type student struct {
	Name   string `bson:"name"`
	Gender string `bson:"gender"`
	Age    int    `bson:"age"`
	Time   time.Time
}

/*

 */
func main() {
	//设置连接参数
	uri := fmt.Sprintf("mongodb://root:cc@172.22.50.25:32082,172.22.50.25:32083,172.22.50.25:32084/?connectTimeoutMS=300000&authSource=admin")
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
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logs.Error("ping mongo failed, err:%s", err.Error())
		return
	}

	//设定连接的数据库和集合

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	var maxCommitTime = time.Second * 10
	//txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc).SetMaxCommitTime(&maxCommitTime)
	// 将会话ID转换为bsoncore.Document类型

	sessionOpts := &options.SessionOptions{}
	sessionOpts.SetDefaultWriteConcern(wc).SetDefaultReadConcern(rc).SetDefaultMaxCommitTime(&maxCommitTime)

	session1, _ := client.StartSession(sessionOpts)
	sessionContext := mongo.NewSessionContext(ctx, session1)
	txn(client, session1, sessionContext, ctx)

}
func txn(client *mongo.Client, session1 mongo.Session, sessionContext mongo.SessionContext, ctx context.Context) {
	database := auth
	collection := "student"
	//构造插入的数据

	students := []interface{}{
		student{
			Name:   "Michael",
			Gender: "Male",
			Age:    26,
			Time:   time.Now(),
		},
		student{
			Name:   "Alice",
			Gender: "Female",
			Age:    27,
			Time:   time.Now().Add(1),
		},
	}

	//开启事务
	if err := sessionContext.StartTransaction(); err != nil {
		panic(err)
	}
	//最后的时候，关闭会话
	defer sessionContext.EndSession(ctx)

	//插入数据
	if _, err := client.Database(database).Collection(collection).InsertMany(sessionContext, students[0:1]); err != nil {
		if err := sessionContext.AbortTransaction(context.Background()); err != nil {
			//回滚事务
			logs.Error("mongo transaction rollback failed, %s", err.Error())
			panic(err)
		}
		panic(err)
	}
	//time.Sleep(1 * time.Second)

	if _, err := client.Database(database).Collection(collection).InsertMany(sessionContext, students[1:]); err != nil {
		if err := sessionContext.AbortTransaction(context.Background()); err != nil {
			//回滚事务
			logs.Error("mongo transaction rollback failed, %s", err.Error())
			panic(err)
		}
		panic(err)
	}
	sessionContext.AbortTransaction(context.Background())
	//提交事务
	sessionContext.CommitTransaction(sessionContext)
}
