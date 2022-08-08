package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"math/rand"
	"strconv"

	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
xxx 多个事务 相互隔离
（通过 http header sessionID  ）
*/
var Client *mongo.Client
var Collection *mongo.Collection

func GenSessionID() string {
	rand.Seed(2)
	//for i := 0; i < 4; i++  {
	//	println(rand.Intn(100))
	//}
	return strconv.Itoa(rand.Intn(100))
}

type TxnManager struct {
	cache map[string]mongo.SessionContext
}

func (t *TxnManager) Release(sessionid string) {
	sessCtx := t.cache[sessionid]
	sessCtx.EndSession(context.Background())

	time.Sleep(1 * time.Second) //使用外部 缓存  或 chan 删除
	delete(t.cache, sessionid)
}
func (t *TxnManager) GenSess() mongo.SessionContext {
	// session 读取策略
	sessOpts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	session, err := Client.StartSession(sessOpts)
	if err != nil {
		log.Fatal(err)
	}
	sessCtx := mongo.NewSessionContext(context.TODO(), session)

	if err = session.StartTransaction(); err != nil {
		panic(err)
	}
	return sessCtx
}
func (t *TxnManager) reloadSession(sessionid string) mongo.SessionContext {
	var sess mongo.SessionContext
	var ok bool
	if sess, ok = t.cache[sessionid]; !ok {
		sess = t.GenSess()
		t.cache[sessionid] = sess
	}
	return sess
}

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
	var err error
	Client, err = mongo.Connect(ctx, opt)
	// collection
	//collection := client.Database(auth).Collection("person111")
	if err != nil {
		log.Fatalln("mongo conn err:", err)
	}

	t := &TxnManager{cache: make(map[string]mongo.SessionContext, 0)}
	// collection
	Collection = Client.Database(auth).Collection("person")

	for i := 0; i < 5; i++ {
		go BeginTxn(t, i)
	}
	time.Sleep(5 * time.Second)
}
func BeginTxn(t *TxnManager, userid int) {
	// 1. 事务 唯一 id, uuid,
	randId := GenSessionID()
	sess := t.reloadSession(randId)
	defer func() {
		t.Release(randId)
	}()
	//defer sess.EndSession(context.Background())
	// insert one
	_, err := Collection.InsertOne(sess, bson.M{"name": "无名小角色", "level": userid})
	if err != nil {
		sess.AbortTransaction(context.Background())
		log.Fatal(err)
	}

	_, err = Collection.InsertOne(sess, bson.M{"name": "无名小角色", "level": userid})
	if err != nil {
		sess.AbortTransaction(context.Background())
		log.Fatal(err)
	}
	Collection.FindOne(sess, bson.M{"name": "无名小角色", "level": userid})

	//xxx 上面都 没有err 则提交
	_ = sess.CommitTransaction(context.Background())

}
