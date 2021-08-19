package main

import (
	funcs "datacenter/funcs"
	"datacenter/g"

	"datacenter/middleware"
	"datacenter/modules"
	"datacenter/routers"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"runtime"
)

// xxx gin doc https://www.kancloud.cn/shuangdeyu/gin_book/949420

const ConfigJson = "config.json"

func main() {

	g.LoadConfig(filepath.Join(g.Basedir, ConfigJson))
	_ = os.Chdir(g.Basedir)
	// 初始化 日志
	g.InitLog()

	// xxx 默认已经连接了 Logger and Recovery 中间件
	//r := gin.Default()

	// windows 无法显示日志颜色
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	} else {
		gin.ForceConsoleColor()
	}

	gin.SetMode("debug")
	// 创建一个默认的没有任何中间件的路由
	r := gin.New()

	// 全局中间件
	// Logger 中间件将写日志到 gin.DefaultWriter 即使你设置 GIN_MODE=release.
	// 默认设置 gin.DefaultWriter = os.Stdout
	// r.Use(gin.Logger())
	// xxx 自定义日志中间件,和django一样,中间件 往返都要执行
	r.Use(middleware.Logmiddleware())
	// xxx 需要将 r.Use(middlewares.Cors()) 在使用路由前进行设置，否则会导致不生效
	r.Use(middleware.Cors())

	// Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	r.Use(gin.Recovery())

	// 初始化mysql conn
	err := modules.MysqlInitConn()
	if err != nil {
		panic(err)
	}

	// 初始化oAuth2
	funcs.InitoAuth2()
	// xxx 加载路由
	routers.LoadRoute(r)

	// xxx 加载clientid
	//funcs.DBLoadClient()

	funcs.InitES()

	r.Run(":8001")
}
