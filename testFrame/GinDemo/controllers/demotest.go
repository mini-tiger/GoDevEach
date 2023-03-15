package controllers

import (
	"GinDemoTest/modules"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/3/15
 * @Desc: demotest.go
**/

func LoginDemo(c *gin.Context) {

	var login modules.Login

	if errA := c.ShouldBindWith(&login, binding.JSON); errA != nil {
		ResponseError(c, errors.New("params err")) //统一返回
	}
	if login.User == "taojun" {
		c.JSON(http.StatusOK, gin.H{
			"status":     http.StatusOK,
			"statusText": "ok",
		})
		return
	}
	ResponseError(c, errors.New("user not taojun")) //统一返回
}
