package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bndr/gojenkins"

	"log"
)

/**
 * @Author: Tao Jun
 * @Description: jenkinsDemo
 * @File:  first
 * @Version: 1.0.0
 * @Date: 2021/12/28 下午6:02
 */
var jenkins *gojenkins.Jenkins

// https://huangzhongde.cn/post/Golang/gojenkins/
// https://github.com/bndr/gojenkins

func main() {
	ctx := context.Background()
	jenkins = gojenkins.CreateJenkins(nil, "http://172.16.12.40:30002/", "admin", "96631cdef3bd4e30a9ff0c6eaa89220f")
	// Provide CA certificate if server is using self-signed certificate
	// caCert, _ := ioutil.ReadFile("/tmp/ca.crt")
	// jenkins.Requester.CACert = caCert
	_, err := jenkins.Init(ctx)

	if err != nil {
		log.Printf("连接Jenkins失败, %v\n", err)
		return
	}
	log.Println("Jenkins连接成功")
	log.Println("Jenkins 节点")
	GetNodes(ctx)

	log.Println("Jenkins job")
	job, err := jenkins.GetJob(ctx, "devops-test")
	if err != nil {
		panic("Job Does Not Exist")
	}
	b, _ := json.MarshalIndent(job.Raw, "", "\t")
	fmt.Println(string(b))

	log.Println("Jenkins all job")
	jobs, err := jenkins.GetAllJobs(ctx)
	if err != nil {
		panic(err)
	}
	for _, j := range jobs {
		//get xml
		fmt.Println(j.GetConfig(ctx))
		//fmt.Println(index, j.GetName())
	}

}
func GetNodes(ctx context.Context) {
	nodes, _ := jenkins.GetAllNodes(ctx)

	for _, node := range nodes {

		// Fetch Node Data
		node.Poll(ctx)
		if ok, _ := node.IsOnline(ctx); ok {
			nodeName := node.GetName()
			log.Printf("Node %s is Online\n", nodeName)
		}
	}
}
