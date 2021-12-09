package main

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/net/context"
	"log"
)

const (
	imageName     string   = "gongt/glibc:bash"           //镜像名称
	containerName string   = "mygin-latest"               //容器名称
	indexName     string   = "/" + containerName          //容器索引名称，用于检查该容器是否存在是使用
	cmd           string   = "/usr/bin/bash && sleep 120" //运行的cmd命令，用于启动container中的程序
	workDir       string   = "/data"                      //container工作目录
	openPort      nat.Port = "80"                         //container开放端口
	hostPort      string   = "8080"                       //container映射到宿主机的端口
	containerDir  string   = "/data/root"                 //容器挂在目录
	hostDir       string   = "/root/data"                 //容器挂在到宿主机的目录
	n             int      = 5                            //每5s检查一个容器是否在运行

)

func main() {

	// 前台
	//createContainer()
	// 后台
	backend()
}

//创建容器
func createContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		panic(err)
	}
	//创建容器
	cont, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      imageName,     //镜像名称
		Tty:        true,          //docker run命令中的-t选项
		OpenStdin:  true,          //docker run命令中的-i选项
		Cmd:        []string{cmd}, //docker 容器中执行的命令
		WorkingDir: workDir,       //docker容器中的工作目录
		ExposedPorts: nat.PortSet{
			openPort: struct{}{}, //docker容器对外开放的端口
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			openPort: []nat.PortBinding{nat.PortBinding{
				HostIP:   "0.0.0.0", //docker容器映射的宿主机的ip
				HostPort: hostPort,  //docker 容器映射到宿主机的端口
			}},
		},
		Mounts: []mount.Mount{ //docker 容器目录挂在到宿主机目录
			mount.Mount{
				Type:   mount.TypeBind,
				Source: hostDir,
				Target: containerDir,
			},
		},
	}, nil, &specs.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}, containerName)
	if err == nil {
		log.Printf("success create container:%s\n", cont.ID)
	} else {
		log.Println("failed to create container!!!!!!!!!!!!!")
	}
}

func backend() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "daocloud.io/library/nginx:1.12.0-alpine",
	}, &container.HostConfig{ //没实现
		PortBindings: nat.PortMap{
			openPort: []nat.PortBinding{nat.PortBinding{
				//HostIP:   "0.0.0.0", //docker容器映射的宿主机的ip
				HostPort: hostPort, //docker 容器映射到宿主机的端口
			}},
		}}, nil, nil, "nginx1.12")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}
