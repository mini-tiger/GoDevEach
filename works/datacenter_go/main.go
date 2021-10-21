package main

import (
	funcs "datacenter/funcs"
	"datacenter/g"
	"datacenter/middleware"
	"datacenter/modules"
	"datacenter/routers"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"runtime"
	"strconv"
)

// xxx gin doc https://www.kancloud.cn/shuangdeyu/gin_book/949420

const ConfigJson = "config.json"

func SetupCfg() {
	_ = g.LoadConfig(ConfigJson)

	// CurrentDir 需要在 LoadConfig 设置
	_ = os.Chdir(g.CurrentDir)

	// 初始化 日志
	g.InitLog()
	g.PrintConf()
}
func SetupServer() (r *gin.Engine) {
	// 默认已经连接了 Logger and Recovery 中间件
	//r := gin.Default()
	// xxx创建一个默认的没有任何中间件的路由
	r = gin.New()

	// windows 无法显示日志颜色
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	} else {
		gin.ForceConsoleColor()
	}
	gin.SetMode(g.GetConfig().Level)
	//gin.SetMode(gin.ReleaseMode)

	// xxx 全局中间件

	// todo ip 白名单
	//r.Use(middleware.IPWhiteList())

	// Logger 中间件将写日志到 gin.DefaultWriter 即使你设置 GIN_MODE=release.
	// 默认设置 gin.DefaultWriter = os.Stdout
	// r.Use(gin.Logger())

	// 自定义日志中间件,和django一样,中间件 往返都要执行
	r.Use(middleware.LogMiddleWare())
	// 需要将 r.Use(middlewares.Cors()) 在使用路由前进行设置，否则会导致不生效
	r.Use(middleware.Cors())
	// Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	//r.Use(gin.Recovery())
	r.Use(middleware.Recovery())
	// xxx 加载路由
	routers.LoadRoute(r)
	return
}

func SetupPlugins() {
	// xxx 初始化mysql conn
	err := modules.MysqlInitConn()
	if err != nil {
		panic(err)
	}

	// xxx 初始化oAuth2
	funcs.InitOAuth2()

	// xxx 初始化ES
	funcs.InitES()
}

var (
	GoVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	//GoVersion   string
	Branch    string
	Commit    string
	BuildTime string
	lowercase string // 小写也可以
)

func printInfo() {
	log.Printf("Go Build version: %s\n", GoVersion)
	log.Printf("Branch: %s\n", Branch)
	log.Printf("Commit: %s\n", Commit)
	log.Printf("BuildTime: %s\n", BuildTime)
	log.Printf("lowercase: %s\n", lowercase)
}

func main() {
	/*
		go build -ldflags=" \
		-X main.Branch=`git rev-parse --abbrev-ref HEAD` \
		-X main.Commit=`git rev-parse HEAD` \
		-X main.BuildTime=`date '+%Y-%m-%d_%H:%M:%S'`" \
		-v -o main main.go
	*/

	versionFlag := flag.Bool("version", false, "print the version")
	flag.Parse()

	if *versionFlag {
		printInfo()
		os.Exit(0)
	}

	printInfo()

	SetupCfg()
	r := SetupServer()
	SetupPlugins()

	_ = r.Run(":" + strconv.Itoa(g.GetConfig().Port))
}
