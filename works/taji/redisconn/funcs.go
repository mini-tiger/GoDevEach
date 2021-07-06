package redisconn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"taji/g"
	"taji/modules"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: redisconn
 * @File:  funcs
 * @Version: 1.0.0
 * @Date: 2021/6/29 下午4:53
 */

var RDB *redis.Client

func Conn() error {
	var ctx = context.Background()
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", g.GetConfig().RedisAddr),
		Password: g.GetConfig().RedisPwd, // no password set
		DB:       g.GetConfig().RedisDB,  // use default DB
	})

	return RDB.Ping(ctx).Err()

}

func LoopRecvPubSub() {
	var ctx = context.Background()

	modules.Mgo.SetCollection("tower_crane", "log_deviceRealTime")

	go func() {
		ps := RDB.Subscribe(ctx, g.GetConfig().SubChan)
		_, err := ps.Receive(ctx)
		if err != nil {
			//panic(err)
			g.GetLog().Error("redis SubChan:%s Receive FAIL\n", g.GetConfig().SubChan)
		}
		ch := ps.Channel()
		for {

			for msg := range ch {
				Real := modules.RealDataEntityFree.Get().(*modules.RealDataEntity)

				if err := json.Unmarshal([]byte(msg.Payload), Real); err == nil {
					//fmt.Printf("%+v\n",*Real)
					Real.InsertTime = time.Now().Unix()

					g.GetLog().Debug("json Unmarshal:%+v\n", *Real)
					//fmt.Printf("%+v\n", *Real) //打印结果：{Tom 123456 [Li Fei]}
					//fmt.Println(modules.Mgo.CollectionCount("tower_crane","auto_user"))
					_, err := modules.Mgo.Collection.InsertOne(context.TODO(), Real)

					if err != nil {
						g.GetLog().Error("Insert mongo err:%s\n", err)
					}
				} else {
					g.GetLog().Error("json 解析失败 %s err:%v\n", msg.Payload, err)
				}
				modules.RealDataEntityFree.Put(Real)
			}
		}

	}()
}
