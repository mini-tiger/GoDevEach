package conf

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/4
 * @Desc: configure.go
**/

var RunMode string = "dev"
var DefaultTimeOut time.Duration = 15 * time.Second
var CurrentDir string

var GlobalCtx context.Context = context.Background()

var configName = "config.yaml"
var ServerPort = "8081"
var ServerAuthUrlSuffix = "/api/v1/identity_verification"
var ServerSendUrlSuffix = "/api/v1/collect"
var SendHttpServer string = ""

func init() {
	_, currentFile, _, _ := runtime.Caller(1)
	//fmt.Println(currentFile)
	CurrentDir = path.Dir(currentFile)
	//fmt.Println(CurrentDir)
	CurrentDir, _ = os.Getwd()
	GetEnv()
}

func GetEnv() {
	log.Printf("当前版本RunMode: %s", RunMode)
	viper.AutomaticEnv()
	mode := viper.GetString("RUNMODE")
	log.Printf("当前环境变量RunMode: %s", mode)
	if mode == "" {
		log.Printf("当前环境变量RunMode: %s", "空")
		return
	}
	RunMode = CheckRunMode(mode)
}
func CheckRunMode(mode string) string {
	switch strings.ToLower(mode) {
	case "dev":
		return "dev"
	case "prod":
		return "prod"
	default:
		return "dev"
	}
}

func ReadYAML() error {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(CurrentDir)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
func GetServerAddr() string {
	//fmt.Println("server: ", viper.GetString("serveraddr"))
	ip := net.ParseIP(viper.GetString("serveraddr"))
	if ip == nil {
		fmt.Println("yaml IP address err")
		return ""
	}
	SendHttpServer = ip.To4().String()
	return ip.To4().String()

}

func SetServerAddr(addr string) error {
	viper.Set("serveraddr", addr)
	SendHttpServer = addr
	return viper.WriteConfigAs(path.Join(CurrentDir, configName))
}
