package main

import (
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

func getzkConnection() *zk.Conn {
	conn, _, err := zk.Connect([]string{"172.22.50.25:32181"}, 5*time.Second)
	if err != nil {
		panic(err)
	}
	return conn
}

func testZookeeperAuth() {
	path, user, pwd := "/auth_test", "panxie", "123456"
	conn := getzkConnection()

	acl := zk.DigestACL(zk.PermAll, user, pwd)
	// 创建节点，带auth
	p, err := conn.Create(path, []byte("hello,world"), 0, zk.WorldACL(zk.PermAll))
	if path != p || err != nil {
		panic(err.Error() + p)
	}

	p, err = conn.Create(path+"/hello", []byte("hello,world"), 0, acl)
	if path+"/hello" != p || err != nil {
		panic(err.Error() + p)
	}
	conn.SetACL(path, acl, -1)

	// 读取节点，不带auth
	_, _, err = conn.Get(path)
	if err == nil {
		panic("read content without auth but no error occured.")
	}
	// 读取节点，带auth
	err = conn.AddAuth("digest", []byte(fmt.Sprintf("%s:%s", user, pwd)))
	if err != nil {
		panic(err)
	}
	cont, _, err := conn.Get(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("content read:%s\n", string(cont))

	conn.Close()
	conn = getzkConnection()
	// 删除节点，不带auth
	err = conn.Delete(path+"/hello", -1)
	if err == nil {
		panic("delete node without auth but no error occured.")
	}
	// 删除节点，带auth
	err = conn.AddAuth("digest", []byte(fmt.Sprintf("%s:%s", user, pwd)))
	if err != nil {
		panic(err)
	}
	err = conn.Delete(path+"/hello", -1)
	if err != nil {
		panic(err)
	}
	// 读取节点的ACL
	conn.SetACL(path, zk.WorldACL(zk.PermAll), -1)
	acls, _, err := conn.GetACL(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("acl get=%v\n", acls)
	// 设置已有节点的ACL
	_, err = conn.SetACL(path, acl, -1)
	if err != nil {
		panic(err)
	}
	// 获取已有节点的ACL
	acls, _, err = conn.GetACL(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("acl get=%v\n", acls)
}

func main() {
	testZookeeperAuth()
}
