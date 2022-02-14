package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.ForwardedByClientIP = true
	// 路由定义post请求, url路径为：/user/login, 绑定doLogin控制器函数
	r.GET("/*path", doLogin)
	r.Run()

}

// 控制器函数
func doLogin(c *gin.Context) {
	path := c.Param("path")

	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"code": http.StatusOK,
			"message": gin.H{
				"Method":     c.Request.Method,
				"Host":       c.Request.Host,
				"RemoteAddr": c.Request.RemoteAddr,
				"Path":       path,
			},
			"clientIP": c.ClientIP(),
			"header":   c.Request.Header,
		},
	})
}
