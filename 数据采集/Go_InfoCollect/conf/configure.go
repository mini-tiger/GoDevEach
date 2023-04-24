package conf

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
	"log"
	"net"
	"path"
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
var LoopScanIpInterval = 30
var CurrentDir string

var GlobalCtx context.Context = context.Background()
var configName = "config.yaml"

var ServerAuthUrlSuffix = "/lstack-combine-server/v1/collect"
var ServerSendUrlSuffix = "/lstack-combine-server/v1/collect"

//var SendHttpServer string = ""

var ServerPort = "30980" // default
var Version = "debug"    // default

var ServerAddrYamlPath string = "server.addr" // 环境变量 SERVER.ADDR
var ServerPortYamlPath string = "server.port" // 环境变量 SERVER.PORT

func GetEnv() {
	getEnv()
	//if RunMode == "dev" {
	//	CurrentDir, _ = os.Getwd()
	//} else {
	//	CurrentDir = GetCurrentAbPath()
	//}

	CurrentDir = GetCurrentAbPath()

}

func getEnv() {
	log.Printf("当前build编译版本RunMode: %s", RunMode)
	viper.BindEnv("RUNMODE")
	mode := viper.GetString("RUNMODE")

	//log.Printf("当前环境变量RunMode: %s", mode)
	if mode == "" {

		log.Printf("当前环境变量: %s, 默认使用: %s\n", "空", RunMode)
		return
	} else {
		RunMode = CheckRunMode(mode)
		log.Printf("当前环境变量: %s, RunMode: %s \n", mode, RunMode)
	}
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

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Read config file:", viper.ConfigFileUsed())
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault(ServerPortYamlPath, ServerPort)
	//fmt.Println(viper.GetString(ServerPortYamlPath))
	//fmt.Println(viper.BindEnv(ServerAddrYamlPath))

	return nil
}

func SetServerPort(portstr string) bool {
	//p, err := strconv.Atoi(portstr)
	//if err != nil {
	//	return false
	//}
	if !govalidator.InRangeInt(portstr, 1, 65530) {
		return false
	}
	viper.Set(ServerPortYamlPath, portstr)
	return true
}

func GetServerPort() string {

	return viper.GetString(ServerPortYamlPath)

}

func SetServerAddr(addr string) bool {
	if !govalidator.IsIPv4(addr) {
		return false
	}
	viper.Set(ServerAddrYamlPath, addr)
	return true
}

func GetServerAddr() string {
	ip := net.ParseIP(viper.GetString(ServerAddrYamlPath))
	if ip == nil {
		//log.Println("Get yaml IP address null")
		return ""
	}
	//SendHttpServer = ip.To4().String()
	return ip.To4().String()

}

func WriteToConfig() error {
	//return viper.WriteConfig()

	return viper.WriteConfigAs(path.Join(CurrentDir, configName))

}
