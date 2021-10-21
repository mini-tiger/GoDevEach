package controllers

import (
	"context"
	funcs "datacenter/funcs"
	"datacenter/g"
	"datacenter/modules"
	"datacenter/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  auth.go
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午3:39
 */

type ResultData struct {
	PageInfo map[string]interface{}   `json:"pageinfo"`
	List     []map[string]interface{} `json:"list"`
}

func CheckToken(c *gin.Context) {

	srv := funcs.OAuthSrv
	err := c.Request.ParseForm()
	if err != nil {
		ResponseAuthError(c, err)
	}

	// 请求参数 更改为 access_token
	at := c.Request.Form["token"]
	c.Request.Form["access_token"] = at

	token, err := srv.ValidationBearerToken(c.Request)
	if err != nil {
		ResponseAuthError(c, err)
		return
	}

	var u modules.User
	//modules.MysqlDb.Table("user").Select("*").Where(" id=? ", token.GetUserID()).First(&u)
	modules.MysqlDb.Where(" id=? ", token.GetUserID()).First(&u)

	//fmt.Println(u)
	//fmt.Println(id)
	resp := map[string]interface{}{
		"active":      true,
		"scope":       []string{"all"},
		"exp":         int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":   token.GetClientID(),
		"user_id":     token.GetUserID(),
		"user_name":   u.Name,
		"authorities": strings.Split(u.Roles, ","),
	}

	ResponseBasic(c, http.StatusOK, resp)
}

func Token(c *gin.Context) {
	//fmt.Println(c.Request.ParseForm())
	//fmt.Println(c.Request.Form)
	s := funcs.OAuthSrv
	if g.GetConfig().IsDebug() {
		u, p, o := c.Request.BasicAuth()
		g.GetLog().Debug(fmt.Sprintf("Token BasicAuth client: %s, password: %s, ok: %v\n", u, p, o))
	}

	gt, tgr, err := s.ValidationTokenRequest(c.Request)
	if err != nil {
		ResponseError(c, err)
		return
	}
	//fmt.Printf("gt: %+v\n", gt)
	//username, password := c.Request.FormValue("username"), c.Request.FormValue("password")
	//fmt.Println(username,password)
	//fmt.Println(strings.Split(tgr.UserID,"_")[1])
	//fmt.Printf("tgr: %+v\n", tgr)

	var ti oauth2.TokenInfo
	switch gt {

	case oauth2.PasswordCredentials:
		//确认用户存在
		//var u modules.User

		//udb := modules.MysqlDb.Table("user").Select("*").Where(" client=? and  name=?", tgr.ClientID, strings.Split(tgr.UserID, g.SepStr)[1]).First(&u)
		//
		//udb := modules.MysqlDb.Where(" id=? ",  strings.Split(tgr.UserID, g.SepStr)[0]).First(&u)
		//
		//if udb.RowsAffected == 0 {
		//	ResponseError(c, errors.New(fmt.Sprintf("UserName: %s , ClientId:%s  DB Not Found", strings.Split(tgr.UserID, g.SepStr)[1], tgr.ClientID)))
		//	return
		//}
		////确认密码eq
		//if u.Password != utils.Md5V3(c.Request.FormValue("password")) {
		//	ResponseError(c, errors.New(fmt.Sprintf("UserName: %s , ClientId:%s  password Not Match", strings.Split(tgr.UserID, g.SepStr)[1], tgr.ClientID)))
		//	return
		//}

		//fmt.Println(time.Now().Unix())
		var tk models.Token
		tdb := modules.MysqlDb.Table(g.TokenTableName).
			Select("*").
			Where(" ClientID=? and  UserID=? and UNIX_TIMESTAMP(AccessCreateAt)+AccessExpiresIn/1000000000  > ?", tgr.ClientID, tgr.UserID, time.Now().Unix()).
			Take(&tk)
		//fmt.Printf("tkdb:%+v\n", tkdb)
		//fmt.Printf("tk:%+v\n", tk)

		if tdb.RowsAffected > 0 {
			// exist  accessToken
			//var tm models.Token
			//err = jsoniter.Unmarshal([]byte(tk.Data), &tm)
			//if err != nil {
			//	ResponseError(c, err)
			//	return
			//}
			//fmt.Printf("tm:%+v\n", tk.GetAccessExpiresIn())
			//fmt.Printf("tm:%+v\n", tk.GetAccessCreateAt())
			ti = &tk
			ti.SetAccessExpiresIn(ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Sub(time.Now()))

			c.JSON(http.StatusOK, s.GetTokenData(ti))
			return
		} else {
			// create    与 refresh 基本逻辑一样
			if g.GetConfig().IsDebug() {
				g.GetLog().Debug("新建 Token 记录 %+v\n", tgr)
			}
			// xxx 新建记录   删除之前创建的记录
			modules.MysqlDb.Table(g.TokenTableName).Where(" ClientID=? and  UserID=? and Code='' ", tgr.ClientID, tgr.UserID).Delete(&models.Token{})
		}

	case oauth2.Refreshing:
		// 与 create 基本逻辑一样
	default:
		ResponseError(c, errors.New("grant_type invalid"))
		return
	}

	tgr.ClientSecret = utils.Md5V3(tgr.ClientSecret)
	//fmt.Printf("tgr:%+v\n", tgr)

	ti, err = s.GetAccessToken(context.Background(), gt, tgr)
	//fmt.Printf("%#v\n",ti)
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, s.GetTokenData(ti))
	return

}

func InitClient(c *gin.Context) {

	secret := c.PostForm("secret")
	resource_ids := c.PostForm("resource_ids")
	scope := c.PostForm("scope")
	authorized_grant_types := c.PostForm("authorized_grant_types")
	web_server_redirect_uri := c.PostForm("web_server_redirect_uri")
	authorities := c.PostForm("authorities")
	access_token_validity := c.PostForm("access_token_validity")
	refresh_token_validity := c.PostForm("refresh_token_validity")
	additional_information := c.PostForm("additional_information")
	autoapprove := c.PostForm("autoapprove")
	client_id := c.PostForm("client_id")

	var db *gorm.DB
	// 查询存在
	var clientDetail *modules.OauthClientDetails
	//db = modules.MysqlDb.Table("oauth_client_details").Select("*").Where("client_id=?", client_id).First(&clientDetail)

	clientDetail, db = funcs.CStore.GetDetailByID(client_id)

	//fmt.Println(clientDetail)
	var u modules.User // 准备用户数据
	u.Name = client_id
	u.Password = utils.Md5V3(secret)
	u.Source = "client_users"
	u.Roles = "admin"
	u.Status = "1"

	// 1. 存在则创建 项目 默认用户
	if db.RowsAffected > 0 {
		modules.MysqlDb.Table("user").Select("id").Where(" name=? and  source='client_users'", client_id).First(&u)
		if u.Id == 0 {

			db = modules.MysqlDb.Table("user").Create(&u)
			//fmt.Printf("%+v\n", u)
			//fmt.Printf("%+v\n", db)
			if db.RowsAffected == 0 {
				ResponseError(c, errors.New("user table insert Fail"))
				return
			}
		}
	} else {
		clientDetail.ClientSecret = utils.Md5V3(secret)
		clientDetail.ClientId = client_id
		clientDetail.WebServerRedirectUri = web_server_redirect_uri
		atv, _ := strconv.Atoi(access_token_validity)

		clientDetail.AccessTokenValidity = atv
		clientDetail.AdditionalInformation = additional_information
		clientDetail.Authorities = authorities
		clientDetail.AuthorizedGrantTypes = authorized_grant_types
		clientDetail.Autoapprove = autoapprove
		rtv, _ := strconv.Atoi(refresh_token_validity)
		clientDetail.RefreshTokenValidity = rtv
		clientDetail.ResourceIds = resource_ids
		clientDetail.Scope = scope

		//db = modules.MysqlDb.Table("oauth_client_details").Create(&clientDetail)
		db = funcs.CStore.CreateDetail(clientDetail)
		//fmt.Printf("%+v\n", u)
		//fmt.Printf("%+v\n", db)
		if db.RowsAffected == 0 {
			ResponseError(c, errors.New("oauth_client_details table insert Fail"))
			return
		} else {
			db = modules.MysqlDb.Table("user").Create(&u)
			//fmt.Printf("%+v\n", u)
			//fmt.Printf("%+v\n", db)
			if db.RowsAffected == 0 {
				ResponseError(c, errors.New("user table insert Fail"))
				return
			}
		}
	}
	ResponseSuccess(c, clientDetail)
	return

}
func DeleteClient(c *gin.Context) {
	//fmt.Println(c.PostForm("client_secret"))
	secret := utils.Md5V3(c.PostForm("client_secret"))
	client_id := c.PostForm("client_id")
	//fmt.Println(client_id)
	//fmt.Println(secret)
	var db *gorm.DB
	//clientDetailDB := modules.MysqlDb.Table("oauth_client_details")
	userDB := modules.MysqlDb.Table("user")
	//fmt.Println(secret)

	// 查询存在
	var clientDetail *modules.OauthClientDetails
	//db = modules.MysqlDb.Table("oauth_client_details").Select("*").Where("client_id=?", client_id).First(&clientDetail)
	clientDetail, db = funcs.CStore.GetDetailByWhere(fmt.Sprintf("client_id='%s' and client_secret='%s'", client_id, secret))

	if db.RowsAffected > 0 {
		funcs.CStore.Delete(clientDetail)

		userDB.Where("name=? and source='client_users'", client_id).Delete(&modules.User{})
		//modules.MysqlDb.Raw("DELETE from user  where name=? and source='client_users'", client_id).Rows()
		ResponseSuccess(c, "success")
		return
	}
	ResponseError(c, fmt.Errorf("Client not Exist"))

}

func UpdateClient(c *gin.Context) {

	secret := utils.Md5V3(c.PostForm("secret"))
	resource_ids := c.PostForm("resource_ids")
	scope := c.PostForm("scope")
	authorized_grant_types := c.PostForm("authorized_grant_types")
	web_server_redirect_uri := c.PostForm("web_server_redirect_uri")
	authorities := c.PostForm("authorities")
	access_token_validity := c.PostForm("access_token_validity")
	refresh_token_validity := c.PostForm("refresh_token_validity")
	additional_information := c.PostForm("additional_information")
	autoapprove := c.PostForm("autoapprove")
	client_id := c.PostForm("client_id")

	var db *gorm.DB
	//clientDetailDB := modules.MysqlDb.Table("oauth_client_details")
	userDB := modules.MysqlDb.Table("user")

	// 查询存在
	var clientDetail *modules.OauthClientDetails
	//db = modules.MysqlDb.Table("oauth_client_details").Select("*").Where("client_id=?", client_id).First(&clientDetail)

	clientDetail, db = funcs.CStore.GetDetailByID(client_id)

	var u modules.User // 准备用户数据

	//fmt.Println(clientDetail)
	// 1.
	if db.RowsAffected > 0 {
		//fmt.Println(clientDetail)
		clientDetail.ClientSecret = secret
		clientDetail.ClientId = client_id
		clientDetail.WebServerRedirectUri = web_server_redirect_uri
		atv, _ := strconv.Atoi(access_token_validity)

		clientDetail.AccessTokenValidity = atv
		clientDetail.AdditionalInformation = additional_information
		clientDetail.Authorities = authorities
		clientDetail.AuthorizedGrantTypes = authorized_grant_types
		clientDetail.Autoapprove = autoapprove
		rtv, _ := strconv.Atoi(refresh_token_validity)
		clientDetail.RefreshTokenValidity = rtv
		clientDetail.ResourceIds = resource_ids
		clientDetail.Scope = scope

		//fmt.Println(clientDetail)
		db = funcs.CStore.Save(clientDetail)

		if db.RowsAffected > 0 { //没有改动 这里是0
			// 查找默认用户 没有则创建
			db = userDB.Select("id").Where(" name=? and  source='client_users'", client_id).First(&u)

			if db.RowsAffected == 0 {
				u.Name = client_id
				u.Password = secret
				u.Source = "client_users"
				u.Roles = "admin"
				u.Status = "1"
				db = modules.MysqlDb.Table("user").Create(&u)
				if db.RowsAffected == 0 {
					ResponseError(c, errors.New("user table insert Fail"))
					return
				}
			} else {
				fmt.Printf("%+v\n", u)
				db.Model(&u).Select("name", "password").Updates(modules.User{Name: client_id, Password: secret})
				fmt.Printf("%+v\n", db)
			}
		}
	} else {
		ResponseError(c, errors.New("clientDetail Not Found"))
	}
	ResponseSuccess(c, clientDetail)
	return
}

func Register(c *gin.Context) {

	var username, source, client, password, role string

	var err error

	if username, err = utils.ParamsVerify(c, "username"); err != nil {
		ResponseError(c, err)
		return
	}
	if password, err = utils.ParamsVerify(c, "password"); err != nil {
		ResponseError(c, err)
		return
	}

	if source, err = utils.ParamsVerify(c, "source"); err != nil {
		ResponseError(c, err)
		return
	}
	if client, err = utils.ParamsVerify(c, "client"); err != nil {
		ResponseError(c, err)
		return
	}

	role = c.PostForm("role")

	//if username, ok = c.GetPostForm("username"); !ok {
	//	ResponseError(c, errors.New("request Params Miss username"))
	//	return
	//}

	//if source, ok = c.GetPostForm("source"); !ok {
	//	ResponseError(c, errors.New("request Params Miss source"))
	//	return
	//}
	//
	//if client, ok = c.GetPostForm("client"); !ok {
	//	ResponseError(c, errors.New("request Params Miss client"))
	//	return
	//}
	//if password, ok = c.GetPostForm("password"); !ok {
	//	ResponseError(c, errors.New("request Params Miss password"))
	//	return
	//}

	//to := modules.MysqlDb.Table(g.ClientTableName).Select("*").Where(" client_id=? ", client).Find(&modules.OauthClientDetails{})

	_, tdb := funcs.CStore.GetDetailByWhere(fmt.Sprintf("client_id='%s'", client))

	if tdb.RowsAffected == 0 {
		ResponseError(c, errors.New(fmt.Sprintf("Not Found Client : %s", client)))
		return
	}

	var u modules.User
	//rows, err := modules.MysqlDb.Raw("SELECT  *  from user where name=? and  source=?",[]string{"name","111"}...).Scan(&u)
	modules.MysqlDb.Table("user").Select("*").Where(" name=? and  source=? and client=? ", username, source, client).Find(&u)

	if u.Id == 0 {
		u.Name = username
		u.Client = client
		u.Password = utils.Md5V3(password)
		u.Source = source
		u.Roles = role
		db := modules.MysqlDb.Table("user").Create(&u)
		//fmt.Printf("%+v\n", u)
		//fmt.Printf("%+v\n", db)
		if db.RowsAffected == 0 {
			ResponseError(c, errors.New("user table insert Fail"))
			return
		}
		if g.GetConfig().IsDebug() {
			g.GetLog().Debug("新建用户 %+v\n", u)
		}

	}

	ResponseSuccess(c, &u)
	return

}

//func InsertClientStoreUser(clientid, userid string) (err error) {
//	var clientdetail modules.OauthClientDetails
//	modules.MysqlDb.Table("oauth_client_details").Select("*").Where(" client_id=? ", clientid).Find(&clientdetail)
//	//fmt.Println(clientdetail)
//	if clientdetail.ClientSecret == "" {
//		err = errors.New("oauth_client_details table select Fail")
//		return
//	}
//
//	// clientStore insert
//	err = funcs.ClientStore.Set(userid, &models.Client{
//		ID:     clientdetail.ClientId,
//		Secret: clientdetail.ClientSecret,
//		Domain: clientdetail.WebServerRedirectUri,
//	})
//	return err
//}
