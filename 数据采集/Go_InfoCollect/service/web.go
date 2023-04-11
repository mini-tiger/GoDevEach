package service

import (
	"collect_web/collect"
	"collect_web/conf"
	"collect_web/log"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/6
 * @Desc: web.go
**/

var HttpSrv http.Server

func StartWeb() {
	// windows 无法显示日志颜色
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	} else {
		gin.ForceConsoleColor()
	}
	if conf.RunMode == "dev" {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	// 创建一个默认的没有任何中间件的路由
	r := gin.New()
	LoadRoute(r)

	HttpSrv = http.Server{
		Addr:    ":1988",
		Handler: r,
	}

	go func() {
		if err := HttpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Glog.Sugar().Errorf("Gin StartWeb Fail Err:%v", err)
		}
	}()

}

func LoadRoute(router *gin.Engine) {

	// Hello World
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello this is Collect Agent")
	})

	//v1组路由
	hgApi := router.Group("/api/v1")

	hgApi.GET("/self", viewSelf)
	hgApi.GET("/collect", viewCollect)
}

func ResponseSuccess(c *gin.Context, val interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"code":    http.StatusOK,
			"message": "ok",
			"data":    val,
		},
	})
}

func ResponseError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		},
	})
}

func viewSelf(c *gin.Context) {
	err, pro := collect.GetSelfProcess()
	if err != nil {
		ResponseError(c, err)
		return
	}
	ResponseSuccess(c, pro)
	return

}

func viewCollect(c *gin.Context) {
	//ctx, _ := context.WithTimeout(conf.GlobalCtx, conf.DefaultTimeOut)
	lx := new(LinuxMetric)
	lx.RegMetrics()
	//var lxface collect.GetMetricInter = lx
	data := lx.GetMetrics(conf.GlobalCtx)

	// log
	b, err := json.Marshal(data)
	if err != nil {
		log.Glog.Error(fmt.Sprintf("web jsonMarshal Err:%v", err))
		return
	}
	log.Glog.Debug(fmt.Sprintf("web collect data:%v", string(b)))
	log.Glog.Error(fmt.Sprintf("web collect err:%v", lx.GetErrors()))
	//

	c.JSON(200, lx.FormatData())
	return

}
