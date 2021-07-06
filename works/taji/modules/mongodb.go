package modules

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"taji/g"
	"time"
)

var Mgo *MongoConn

// 连接设置
type MongoConn struct {
	MongoClient *mongo.Client
	Collection  *mongo.Collection
}

func NewMongoConn() (err error) {
	//uri := "mongodb+srv://用户名:密码@官方给的.mongodb.net"
	uri := g.GetConfig().MongoUri
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri)) // 连接池
	if err != nil {
		return err
	}
	// 检查连接
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	Mgo = &MongoConn{MongoClient: MongoClient}
	return nil
	//log.Println("Connected to MongoDB!")
}

//func (m *MongoConn) CollectionCount(db, coll string) (int64, error) {
//	collection := m.MongoClient.Database(db).Collection(coll)
//	//name := collection.Name()
//	return collection.EstimatedDocumentCount(context.TODO())
//	//return
//}
//
//func (m *MongoConn) CollectionFilter(db, coll string,
//	Filter interface{}, findOptions *options.FindOptions) ([]bson.M, error) {
//
//	collection := m.MongoClient.Database(db).Collection(coll)
//	// 查询多个
//	// 将选项传递给Find()
//
//	// 定义一个切片用来存储查询结果
//	var results []bson.M
//
//	// 把bson.D{{}}作为一个filter来匹配所有文档
//	cur, err := collection.Find(context.TODO(), Filter, findOptions)
//	if err != nil {
//		return results, err
//	}
//
//	// 查找多个文档返回一个光标
//	// 遍历游标允许我们一次解码一个文档
//	for cur.Next(context.TODO()) {
//		// 创建一个值，将单个文档解码为该值
//		var elem bson.M
//		err := cur.Decode(&elem)
//		if err != nil {
//			continue
//		}
//		results = append(results, elem)
//	}
//
//	if err := cur.Err(); err != nil {
//		return results, err
//	}
//
//	// 完成后关闭游标
//	cur.Close(context.TODO())
//	return results, nil
//}

func (m *MongoConn) SetCollection(database, col string) {
	m.Collection = m.MongoClient.Database(database).Collection(col)
}

func (m *MongoConn) InsertOne(value interface{}) (*mongo.InsertOneResult, error) {

	return m.Collection.InsertOne(context.TODO(), value)
}

func (m *MongoConn) DisableConn() error {

	//name := collection.Name()
	return m.MongoClient.Disconnect(context.TODO())
	//return
}
