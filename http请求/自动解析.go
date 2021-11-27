package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  post
 * @Version: 1.0.0
 * @Date: 2021/11/8 下午5:28
 */

type Library struct {
	Name   string
	Latest string
}

type Libraries struct {
	Results []*Library
}

func main() {
	client := resty.New()

	libraries := &Libraries{}
	client.R().SetResult(libraries).Get("https://api.cdnjs.com/libraries")
	fmt.Printf("%d libraries\n", len(libraries.Results))

	for _, lib := range libraries.Results {
		fmt.Println("first library:")
		fmt.Printf("name:%s latest:%s\n", lib.Name, lib.Latest)
		break
	}

}
