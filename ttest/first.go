package main

import (
	"bytes"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os/exec"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  watch
 * @Version: 1.0.0
 * @Date: 2021/11/15 下午6:30
 */

const defaultConfigFile = "values.yaml"
const harbor = "misharbor.dyxnet.com/erp-lastest/"
const harbor21 = "harbor.dev.21vianet.com/dyxnet-erp/"

var Conf = new(Config)

func main() {
	var configfile = flag.StringP("flagname", "f", defaultConfigFile, "config yaml")

	// 设置非必须选项的默认值
	//flag.Lookup("flagname").NoOptDefVal = "4321"
	flag.Parse()
	fmt.Println(*configfile)
	time.Sleep(1 * time.Second)
	v := viper.New()
	v.SetConfigFile(*configfile) // 指定配置文件路径
	v.SetConfigType("yaml")
	err := v.ReadInConfig() // 读取配置信息
	if err != nil {         // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := v.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	if len(Conf.FnArgs) < 0 {
		log.Fatalln("len is 0")

	}
	var pullCmd, retag, sourceImage, destImage, pushCmd string
	for _, image := range Conf.FnArgs {
		sourceImage = fmt.Sprintf("%s%s:%s", harbor, image.Name, image.Tag)
		destImage = fmt.Sprintf("%s%s:%s", harbor21, image.Name, image.Tag)

		pullCmd = fmt.Sprintf("docker pull %s", sourceImage)
		retag = fmt.Sprintf("docker tag %s %s", sourceImage, destImage)
		pushCmd = fmt.Sprintf("docker push %s", destImage)
		//fmt.Println(pullCmd)
		log.Printf("begin image:%v\n", image)
		err, _ := RunCommand(pullCmd)
		if err != "" {
			log.Printf("pull error :%s\n", image)
			time.Sleep(10 * time.Second)
		}

		//fmt.Println(retag)
		err, _ = RunCommand(retag)
		if err != "" {
			log.Printf("retag error :%s\n", image)
			time.Sleep(10 * time.Second)
		}

		//fmt.Println(pushCmd)
		err, _ = RunCommand(pushCmd)
		if err != "" {
			log.Printf("push error :%s\n", image)
			time.Sleep(10 * time.Second)
		}
		log.Printf("done image:%v\n", image)
	}
}
func RunCommand(command string) (string, string) {
	cmd := exec.Command("/bin/bash", "-c", command)

	var out bytes.Buffer
	cmd.Stdout = &out

	var e bytes.Buffer
	cmd.Stderr = &e
	//cmd.Start()
	//buf, _ := cmd.CombinedOutput()
	cmd.Run()
	if e.Len() != 0 && out.Len() == 0 {
		return e.String(), out.String()
	}
	return "", out.String()
}

type Config struct {
	FnArgs []Args
}
type Args struct {
	Name string `yaml:"name"`
	Tag  string `yaml:"tag"`
}
