package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"matchAllRoutes/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 全局中间件
	// Logger 中间件将写日志到 gin.DefaultWriter 即使你设置 GIN_MODE=release.
	// 默认设置 gin.DefaultWriter = os.Stdout
	// r.Use(gin.Logger())
	// xxx 自定义日志中间件,和django一样,中间件 往返都要执行
	r.Use(middleware.Logmiddleware())

	// Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	r.Use(gin.Recovery())
	r.ForwardedByClientIP = true
	// 路由定义post请求, url路径为：/user/login, 绑定doLogin控制器函数
	r.GET("/*path", doLogin)
	r.POST("/*path", doLogin)
	r.Run(":8081")

}

// 控制器函数
func doLogin(c *gin.Context) {
	path := c.Param("path")
	//fmt.Printf("%+v\n",c.Request.Header)
	var mapbody map[string]interface{}

	if c.Request.Method == "POST" {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			jsonData = []byte(fmt.Sprintf("err post data : %s", err.Error()))
		}
		fmt.Println(string(jsonData))
		mapbody = make(map[string]interface{})

		json.Unmarshal(jsonData, &mapbody)
	}
	c.Set("reqbody", mapbody)
	c.JSON(http.StatusOK, gin.H{

		"code": http.StatusOK,
		"message": gin.H{
			"Method":     c.Request.Method,
			"Host":       c.Request.Host,
			"RemoteAddr": c.Request.RemoteAddr,
			"Path":       path,
		},
		"clientIP": c.ClientIP(),
		"header":   c.Request.Header,
		"PostData": mapbody,
	})
}
