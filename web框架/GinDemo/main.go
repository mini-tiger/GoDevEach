package main

import (
	"GinDemo/middleware"
	"GinDemo/modules"
	"GinDemo/routers"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// xxx gin doc https://www.kancloud.cn/shuangdeyu/gin_book/949420

func main() {
	// xxx 默认已经连接了 Logger and Recovery 中间件
	r := gin.Default()

	// 全局中间件
	// Logger 中间件将写日志到 gin.DefaultWriter 即使你设置 GIN_MODE=release.
	// 默认设置 gin.DefaultWriter = os.Stdout
	// r.Use(gin.Logger())
	// xxx 自定义日志中间件,和django一样,中间件 往返都要执行
	r.Use(middleware.Logmiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	r.Use(gin.Recovery())

	// 初始化mysql conn
	modules.MysqlInitConn()
	// xxx 加载路由
	routers.LoadRoute(r)

	r.Run(":8001")
}
