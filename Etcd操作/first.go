package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

var (
	err error
)

type EtcdClient struct {
	Client *clientv3.Client
	Config clientv3.Config
	//Kv     clientv3.KV
	Lease         clientv3.Lease
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

//func (this *EtcdClient) NewKv() {
//	if this.Kv == nil {
//		this.Kv = clientv3.NewKV(this.Client)
//	}
//}

func (this *EtcdClient) PutKV(key, value string, ops ...clientv3.OpOption) error {

	//this.NewKv()

	var putResp *clientv3.PutResponse
	if putResp, err = this.Client.Put(context.TODO(), key, value, ops...); err != nil {
		//log.Fatalln(err)
		return err
	}
	//fmt.Println(putResp.Header.Revision)
	//if putResp.PrevKv != nil {
	//	fmt.Printf("prev Value: %s \nCreateRevison :%d \n", string(putResp.PrevKv.Value), putResp.PrevKv.CreateRevision)
	fmt.Printf("写入Key 返回%+v\n", putResp)
	//}
	return nil
}

func (this *EtcdClient) GetKey(key string, ops ...clientv3.OpOption) (err error, getResp *clientv3.GetResponse) {
	//this.NewKv()

	if getResp, err = this.Client.Get(context.TODO(), key, ops...); err != nil {
		return
	}
	fmt.Printf("获取Key 返回%+v\n", getResp)
	return
}

func (this *EtcdClient) DelKey(key string, ops ...clientv3.OpOption) error {
	//this.NewKv()
	var delResp *clientv3.DeleteResponse
	if delResp, err = this.Client.Delete(context.TODO(), key, ops...); err != nil {
		return err
	}
	fmt.Printf("删除key 返回%+v\n", delResp)
	if len(delResp.PrevKvs) > 0 {
		for _, kvpair := range delResp.PrevKvs {
			fmt.Printf("delete key %+v\n", kvpair)
		}
	}

	return nil
}

func (this *EtcdClient) GetLease(t int64) (err error, leaseGrantResp *clientv3.LeaseGrantResponse) {
	// 申请一个租约
	this.Lease = clientv3.NewLease(this.Client)

	if leaseGrantResp, err = this.Lease.Grant(context.TODO(), t); err != nil {
		//fmt.Println(err)
		return
	}
	return
}

func (this *EtcdClient) AssignLease(leaseId clientv3.LeaseID, keydatas map[string]interface{}) {
	//this.NewKv()

	for key, value := range keydatas {
		if putResp, err := this.Client.Put(context.TODO(), key, value.(string), clientv3.WithLease(leaseId)); err == nil {
			//fmt.Println(err)
			fmt.Printf("授权租约 %+v\n", putResp)
			return

		}
	}
}

func (this *EtcdClient) KeepAlive(leaseId clientv3.LeaseID) {
	//var keepRespChan <-chan *clientv3.LeaseKeepAliveResponse

	if keepRespChan, err := this.Client.KeepAlive(context.Background(), leaseId); err != nil {
		this.keepAliveChan = keepRespChan
		return
	}

}

//ListenLeaseRespChan 监听 续租情况
func (this *EtcdClient) ListenLeaseRespChan() {
	for leaseKeepResp := range this.keepAliveChan {
		log.Println("续约成功", leaseKeepResp)
	}
	log.Println("关闭续租")
}

func main() {
	EtcCli := &EtcdClient{}
	// 客户端配置
	EtcCli.Config = clientv3.Config{
		Endpoints:   []string{"172.22.50.25:31015"},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "123456",
	}
	// 建立连接
	if EtcCli.Client, err = clientv3.New(EtcCli.Config); err != nil {
		log.Fatalln(err)
		return
	}

	//put
	_ = EtcCli.PutKV("/test1/key1", "value1", clientv3.WithPrevKV())
	_ = EtcCli.PutKV("/test1/key2", "value2", clientv3.WithPrevKV())
	//get
	_, _ = EtcCli.GetKey("/test1/key1")

	//del 删除所有key,并分别返回
	_ = EtcCli.DelKey("/test1/key1", clientv3.WithFromKey(), clientv3.WithPrevKV())

	println("-----------------------------------------------------------------------")

	//put
	_ = EtcCli.PutKV("/test1/key3", "value3", clientv3.WithPrevKV())
	_ = EtcCli.PutKV("/test1/key4", "value4")
	_ = EtcCli.PutKV("/test1/key5", "value5")
	//get
	_, _ = EtcCli.GetKey("/test1/", clientv3.WithPrefix())

	//del 删除key前缀
	//_ = EtcCli.DelKey("/test1/", clientv3.WithPrefix(), clientv3.WithPrevKV())

	println("-----------------------------------------------------------------------")

	//申请租约
	_, lease := EtcCli.GetLease(5)

	//授权给key
	EtcCli.AssignLease(lease.ID, map[string]interface{}{"/test1/key5": "value5"})

	//查看租约剩余时间
	for {
		_, getResp := EtcCli.GetKey("/test1/key5")

		if getResp.Count == 0 {
			fmt.Println("/test1/key5 过期了")
			break
		}
		fmt.Println("/test1/key5 还没过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)

		//以下步骤 永远不过期
		EtcCli.KeepAlive(lease.ID)

	}

	println("-----------------------------------------------------------------------")

}
