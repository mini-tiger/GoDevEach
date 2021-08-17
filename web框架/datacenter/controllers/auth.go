package controllers

import (
	"datacenter/authfunc"
	"datacenter/modules"
	"datacenter/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/models"
	"net/http"
	"strconv"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  auth.go
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午3:39
 */

func Token(c *gin.Context) {

	// 1 通过 username 分解 name  client

	//if u.Id > 0{
	//	// 2.1 clientStore  have  return clientinfo
	//	_, err := authfunc.ClientStore.GetByID(context.Background(), strconv.FormatInt(u.Id, 10))
	//
	//	// 2.2 可能重启 丢失cleintStore user
	//	if err.Error() == "not found" {
	//		err=InsertClientStoreUser(client, strconv.FormatInt(u.Id, 10))
	//	}
	//
	//	if err != nil {
	//		ResponseError(c, err)
	//		return
	//	}
	//}

	err := authfunc.OAuthSrv.HandleTokenRequest(c.Writer, c.Request)
	fmt.Println("1111", err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{

			"status":     http.StatusInternalServerError,
			"statusText": "ok",
		})
	}
}

func Register(c *gin.Context) {

	username := c.PostForm("username")
	//fmt.Println(username)

	source := c.PostForm("source")
	client := c.PostForm("client")
	password := c.PostForm("password")
	role := c.PostForm("role")
	//fmt.Println(source)

	var u modules.User
	//rows, err := modules.MysqlDb.Raw("SELECT  *  from user where name=? and  source=?",[]string{"name","111"}...).Scan(&u)
	modules.MysqlDb.Table("user").Select("*").Where(" name=? and  source=?", username, source).Find(&u)

	if u.Id == 0 {
		u.Name = username
		u.Client = client
		u.Password = utils.Md5V3(password)
		u.Source = source
		u.Roles = role
		db := modules.MysqlDb.Table("user").Create(&u)
		//fmt.Printf("%+v\n", u)
		fmt.Printf("%+v\n", db)
		if db.RowsAffected == 0 {
			ResponseError(c, errors.New("user table insert Fail"))
			return
		}
		err := InsertClientStoreUser(client, strconv.FormatInt(u.Id, 10))
		if err != nil {
			ResponseError(c, err)
			return
		}
	}

	ResponseSuccess(c, u)
	return

}

func InsertClientStoreUser(clientid, userid string) (err error) {
	var clientdetail modules.OauthClientDetails
	modules.MysqlDb.Table("oauth_client_details").Select("*").Where(" client_id=? ", clientid).Find(&clientdetail)
	//fmt.Println(clientdetail)
	if clientdetail.ClientSecret == "" {
		err = errors.New("oauth_client_details table select Fail")
		return
	}

	// clientStore insert
	err = authfunc.ClientStore.Set(userid, &models.Client{
		ID:     clientdetail.ClientId,
		Secret: clientdetail.ClientSecret,
		Domain: clientdetail.WebServerRedirectUri,
	})
	return err
}
