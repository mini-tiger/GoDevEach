package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func MyBenchLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置 example 变量
		c.Set("example", "12345")

		// 请求前
		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Printf("latency=%v\n", latency)

		// 获取发送的 status
		status := c.Writer.Status()
		log.Printf("status=%v\n", status)
		if v, ok := c.Get("example"); ok {
			log.Printf("example=%v\n", v)
		}

	}
}

func benchEndpoint(c *gin.Context) {
	c.String(200, "hello benchmark")
}
func loginEndpoint(c *gin.Context) {
	c.String(200, c.Request.URL.String())
}
func submitEndpoint(c *gin.Context) {
	c.String(200, c.Request.URL.String())
}
func readEndpoint(c *gin.Context) {
	c.String(200, c.Request.URL.String())
}
func analyticsEndpoint(c *gin.Context) {
	c.String(200, c.Request.URL.String())
}

func AuthRequire() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	})
}

func main() {
	// 新建一个没有任何默认中间件的路由
	r := gin.New()

	// 全局中间件
	// Logger 中间件将日志写入 gin.DefaultWriter，即使你将 GIN_MODE 设置为 release。
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery 中间件会 recover 任何 panic。如果有 panic 的话，会写入 500。
	r.Use(gin.Recovery())

	// 你可以为每个路由添加任意数量的中间件。
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// 认证路由组
	// authorized := r.Group("/", AuthRequired())
	// 和使用以下两行代码的效果完全一样:
	authorized := r.Group("/v1")
	// 路由组中间件! 在此例中，我们在 "authorized" 路由组中使用自定义创建的
	// AuthRequired() 中间件
	authorized.Use(AuthRequire())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.GET("/read", readEndpoint)

		// 嵌套路由组
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8081")
}
