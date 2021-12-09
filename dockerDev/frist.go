package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	images, err := cli.ImageList(context.TODO(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		fmt.Printf("image:%#v\n", image)
	}
	// 获取本地已启动的Docker容器，如果要查看全部容器，可以配置types.ContainerListOptions{}参数
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("容器ID：%s，容器名称：%s,容器状态：%s\n", container.ID[:10], container.Image, container.State)
	}
}
