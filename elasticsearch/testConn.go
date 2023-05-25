package main

import "fmt"
import "github.com/olivere/elastic/v7"

func main() {
	// 创建ES client用于后续操作ES
	_, err := elastic.NewClient(
		// 设置ES服务地址，支持多个地址
		elastic.SetURL("http://172.22.50.25:32572"),
		// 设置基于http base auth验证的账号和密码
		elastic.SetSniff(false), // docker es
		//elastic.SetBasicAuth("user", "secret")
	)

	if err != nil {
		// Handle error
		fmt.Printf("连接失败: %v\n", err)
	} else {
		fmt.Println("连接成功")
	}
}
