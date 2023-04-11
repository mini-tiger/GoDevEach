package conf

import (
	"context"
	"github.com/spf13/viper"
	"log"
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
var SendHttpUrl string = "http://172.22.50.191:8081"
var GlobalCtx context.Context = context.Background()

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
