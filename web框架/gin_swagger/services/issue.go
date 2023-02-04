package services

import (
	"GinSwagger/models"
	"github.com/gin-gonic/gin"
)

// IssueList function
// @Summary issue列表
// @version 1.0
// @Tags issue
// @Success 200  {object} models.Result{data=[]models.Issue} ”成功后返回"
// @Router /issue/list [get]
func IssueList(c *gin.Context) {
	var issue models.Issue
	//issues := issue.List()

	c.JSON(200, gin.H{
		"stat":    0,
		"message": "列出成功",
		"data":    issue,
	})
}

// IssueAdd function
// @Summary 创建issue summary
// @version 1.0
// @Accept  json
// @Produce  json
// @Description 创建issue desc
// @Tags 创建issue summary tag
// @Param data body models.Issue true "issue params"
// @Param who query string true "人名"
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param BK_NAME header string false "admin"
// @Success 200  object models.Result "成功后返回"
// @Header  200,400,default  {string}  Token     "token"
// @Router /issue/add [post]
// @BasePath  /api/v1
func IssueAdd(c *gin.Context) {
	c.JSON(200, models.Result{
		Stat:    0,
		Message: "",
		Data:    nil,
	})
}
