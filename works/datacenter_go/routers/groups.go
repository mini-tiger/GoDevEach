package routers

import (
	"datacenter/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRoute(router *gin.Engine) {

	// Hello World
	router.GET("/", func(c *gin.Context) {
		//c.String(http.StatusOK, "Hello, World,TaoJun First Deploy for DataCenter")
		c.JSON(http.StatusOK, gin.H{
			//"status":        http.StatusOK,
			"statusText": "Hello, World,TaoJun First Deploy for DataCenter",
		})
	})

	//xxx gin 自带 json parse方式 https://cloud.tencent.com/developer/article/1689928

	// POST 测试JSON 速度
	router.POST("/login1", controllers.PostSpeed1)
	router.POST("/login2", controllers.PostSpeed2)
	router.POST("/login3", controllers.PostSpeed3)

	//v1组路由
	hgApi := router.Group("/hg/api")

	hgApi.POST("/mongoFind", controllers.GetHisMissionDetail) // mongo example
	hgApi.POST("/mysqlFind", controllers.MysqlFind)           // mysql example  https://gorm.io/zh_CN/docs/sql_builder.html
	hgApi.POST("/singleupload", controllers.SingleUpLoad)     //

	// work auth
	authApi := router.Group("/oauth")
	authApi.POST("/token", controllers.Token)
	authApi.POST("/check_token", controllers.CheckToken)

	cfgApi := router.Group("/config")
	cfgApi.POST("/register", controllers.Register)
	cfgApi.POST("/initClient", controllers.InitClient)
	cfgApi.POST("/updateClient", controllers.UpdateClient)
	cfgApi.POST("/deleteClient", controllers.DeleteClient)

	// work api
	Api := router.Group("/api")
	Api.POST("/queryByEs", controllers.QueryByEs)
	Api.POST("/queryByEsPlus", controllers.QueryByEsPlus)
	Api.POST("/querySourceByEs", controllers.QueryBySourceEs)
}
