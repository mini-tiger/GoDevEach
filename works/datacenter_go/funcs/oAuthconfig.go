package authfunc

import (
	"datacenter/g"
	"datacenter/modules"
	"datacenter/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"strconv"
	"time"

	//"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: authfunc
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午3:26
 */
var (
	OAuthSrv *server.Server
	MStore   *modules.StoreEntry
	CStore   *ClientStore
)

func InitOAuth2() {
	manager := manage.NewDefaultManager()
	//manager.SetAuthorizeCodeTokenCfg(&manage.Config{
	//	AccessTokenExp: time.Hour * 10, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true})

	manager.SetPasswordTokenCfg(&manage.Config{
		AccessTokenExp: time.Hour * 24, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true})
	// token store

	// use mysql token store
	//fmt.Println(g.GetConfig().Mysqldsn)

	// xxx NewStoreWithDB 中 包含 清理过期 数据库记录
	Store := modules.NewStore(
		modules.NewConfig(g.GetConfig().MysqlDsn), g.GetConfig().TokenTableName, 300,
	)

	MStore = modules.NewStoreEntry(Store)
	manager.MapTokenStorage(MStore)

	// client 数据库
	var err error
	CStore, err = NewClientStore(modules.MysqlDb, g.GetConfig().ClientTableName)

	if err != nil {
		g.GetLog().Panicf("ClientStore Fail %+v\n", err)
	}
	// ClientStorage
	manager.MapClientStorage(CStore)

	OAuthSrv = server.NewServer(server.NewConfig(), manager)
	//OAuthSrv.SetUserAuthorizationHandler(userAuthorizeHandler)
	OAuthSrv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		// 1 通过 username 分解 name  client
		sDec, err := base64.StdEncoding.DecodeString(username)
		if err != nil {
			return "", errors.New(fmt.Sprintf("base64 decode err : %s", err))
		}
		var client string
		ucList := strings.Split(string(sDec), "|")
		if strings.Contains(string(sDec), "|") && len(ucList) == 2 {
			username = ucList[0]
			client = ucList[1]
		} else {
			err = errors.New(fmt.Sprintf("username not Contains |"))
			return
		}
		if g.GetConfig().IsDebug() {
			g.GetLog().Debug("/auth/token username:%s clientId:%s\n", username, client)
		}

		// 2 用 username and client 查询userid,   es or  mysql
		var u modules.User
		//modules.MysqlDb.Table("user").Select("*").Where(" name=? and  client=?", username, client).Find(&u)

		udb := modules.MysqlDb.Where(" name=? and  client=?", username, client).Find(&u)
		if udb.RowsAffected == 0 {
			err = errors.New(fmt.Sprintf("UserName: %s , ClientId:%s  DB Not Found", username, client))
			return
		}

		//确认密码eq
		if u.Password != utils.Md5V3(password) {
			err = errors.New(fmt.Sprintf("UserName: %s , ClientId:%s  password Not Match", username, client))
			return
		}

		userID = fmt.Sprintf("%s%s%s", strconv.FormatInt(u.Id, 10), g.SepStr, username)

		return
	})
}

//func DBLoadClient() {
//	var clientdetailList []modules.OauthClientDetails
//	modules.MysqlDb.Table("oauth_client_details").Select("*").Find(&clientdetailList)
//	for _, cd := range clientdetailList {
//		InsertClientStoreUser(&cd)
//	}
//}
//
//func InsertClientStoreUser(clientdetail *modules.OauthClientDetails) (err error) {
//
//	// clientStore insert
//	err = CStore.Create(clientdetail)
//	return err
//}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	fmt.Println(r)
	store, err := session.Start(r.Context(), w, r)
	fmt.Println(store)
	fmt.Println(err)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	fmt.Println(uid)
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}
