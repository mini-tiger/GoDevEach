package etcdv3

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

//ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	cli     *clientv3.Client //etcd client
	leaseID clientv3.LeaseID //租约ID
	//租约keepalieve相应chan
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string //key
	val           string //value
}

//NewServiceRegister 新建注册服务
func NewServiceRegister(endpoints []string, serName, addr string, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	// https://blog.csdn.net/hello_ufo/article/details/89343957
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	_, err = cli.Status(timeoutCtx, endpoints[0])
	if err != nil {
		//errors.Wrapf(err, "error checking etcd status: %v", err)
		log.Fatalf("ETCD 链接失败:%s ",err)
	}


	ser := &ServiceRegister{
		cli: cli,
		key: "/" + schema + "/" + serName + "/" + addr,
		val: addr,
	}

	//申请租约设置时间keepalive
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}

	return ser, nil
}

//设置租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	//设置租约时间
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}


	//设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan

	log.Printf("Put key:%s  val:%s  success!", s.key, s.val)



	return nil
}

//ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		log.Println("续约成功", leaseKeepResp)
	}
	log.Println("关闭续租")
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Println("撤销租约")
	return s.cli.Close()
}
