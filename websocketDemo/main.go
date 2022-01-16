package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"websocketDemo/controllers"
)

/**
 * @Author: Tao Jun
 * @Description: websocketDemo
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/12/16 下午5:03
 */

func main() {
	router := gin.Default()
	wsv1 := router.Group("/ws")
	wsv1.GET("/wx", controllers.WsHandle)

	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.Handle("GET", "/", func(context *gin.Context) {
		// 返回HTML文件，响应状态码200，html文件名为index.html，模板参数为nil
		context.HTML(http.StatusOK, "test.html", nil)
	})

	router.Run(":8888")
}
