package routers

import (
	"GinDemo/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRoute(router *gin.Engine) {

	// Hello World
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World,Taojun First Deploy for KubeSphere")
	})

	//xxx gin 自带 json parse方式 https://cloud.tencent.com/developer/article/1689928

	// POST 测试JSON 速度
	router.POST("/login1", controllers.PostSpeed1)
	router.POST("/login2", controllers.PostSpeed2)

	//v1组路由
	hgApi := router.Group("/hg/api")

	hgApi.POST("/mongoFind", controllers.GetHisMissionDetail) // mongo example
	hgApi.POST("/mysqlFind", controllers.MysqlFind)           // mysql example  https://gorm.io/zh_CN/docs/sql_builder.html
	hgApi.POST("/singleupload", controllers.SingleUpLoad)     //
}
