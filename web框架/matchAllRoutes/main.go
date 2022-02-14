package main

import (
	"fmt"
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
			"code":     http.StatusOK,
			"message":  fmt.Sprintf("method:%s,host:%v,path:%s", c.Request.Method, c.Request.Host, path),
			"clientIP": c.ClientIP(),
			"header":   c.Request.Header,
		},
	})
}
