package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"strconv"
	"time"
)

func main() {


	// 客户端配置
	var err error
	Config := clientv3.Config{
		Endpoints:   []string{"192.168.43.26:2379"},
		DialTimeout: 5 * time.Second,
	}
	var Client *clientv3.Client
	// 建立连接
	if Client, err = clientv3.New(Config); err != nil {
		log.Fatalln(err)
		return
	}



	var kv clientv3.KV
	kv = clientv3.NewKV(Client)

	// 模拟KV的变化
	go func() {
		for i:=0;i<5;i++{
			_, err = kv.Put(context.TODO(), "/school/class/students", "helios"+strconv.Itoa(i))
			_, err = kv.Delete(context.TODO(), "/school/class/students")
			time.Sleep(1 * time.Second)
		}
	}()

	// 先GET到当前的值，并监听后续变化
	var getResp *clientv3.GetResponse
	if getResp, err = kv.Get(context.TODO(), "/school/class/students"); err != nil {
		fmt.Println(err)
		return
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值:", string(getResp.Kvs[0].Value))
	}

	// 获得当前revision
	//watchStartRevision := getResp.Header.Revision + 1

	//var watcher clientv3.Watcher
	// 创建一个watcher
	//watcher = clientv3.NewWatcher(Client)
	//fmt.Println("从该版本向后监听:", watchStartRevision)

	//ctx, cancelFunc := context.WithCancel(context.TODO())
	//time.AfterFunc(5*time.Second, func() {
	//	cancelFunc()
	//})

	watchRespChan := Client.Watch(context.Background(), "/school/class/", clientv3.WithPrefix())
	// 处理kv变化事件
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), ", Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}
