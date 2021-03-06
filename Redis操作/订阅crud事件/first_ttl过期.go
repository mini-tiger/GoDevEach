package main

// xxx Keyspace notifications  必须更改redis配置
// https://redis.io/docs/manual/keyspace-notifications/
// https://blog.csdn.net/zhizhengguan/article/details/90575438
// https://segmentfault.com/a/1190000040683431

//xxx redis-helm
/*
value.yaml
master:
....
  extraFlags:
     - --notify-keyspace-events "KEA"

*/
import (
	"context"
	"fmt"
	"time"
)

var DB int = 0
var redisCli *redis.Client

var key string = "/cc/services/config/common"

func init() {
	// 连接redis
	redisCli = redis.NewClient(&redis.Options{
		Addr:     "172.22.50.25:6379",
		Password: "cc",
		DB:       DB,
	})
}

/*
 * redis key 过期自动通知
 */
func SetExpireEvent() {
	// 设置一个键，并且3秒钟之后过期
	redisCli.Set(context.Background(), key, "测试键值过期通知", 3*time.Second)
}

// http://blog.itpub.net/69955379/viewspace-2792316/
func SubExpireEvent() {
	// 订阅key过期事件
	//sub := redisCli.Subscribe(context.Background(), fmt.Sprintf("__keyevent@%d__:expired", DB))
	sub := redisCli.Subscribe(context.Background(), fmt.Sprintf("__keyspace@%d__:%s", DB, key))
	// 这里通过一个for循环监听redis-server发来的消息。
	// 当客户端接收到redis-server发送的事件通知时，
	// 客户端会通过一个channel告知我们。我们再根据
	// msg的channel字段来判断是不是我们期望收到的消息，
	// 然后再进行业务处理。
	for {
		msg := <-sub.Channel()
		fmt.Println("Channel ", msg.Channel)
		fmt.Println("pattern ", msg.Pattern)
		fmt.Println("pattern ", msg.Payload)
		fmt.Println("PayloadSlice ", msg.PayloadSlice)
		fmt.Printf("%#v\n", msg)
	}
}

func main() {
	SetExpireEvent()
	go SubExpireEvent()

	// 这里sleep是为了防止main方法直接推出
	time.Sleep(100 * time.Second)
}
