package authfunc

import (
	"datacenter/modules"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-oauth2/mysql/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
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
	dumpvar     bool
	idvar       string
	secretvar   string
	domainvar   string
	portvar     int
	OAuthSrv    *server.Server
	Mstore      *mysql.Store
	ClientStore *store.ClientStore
)

func init() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store

	// use mysql token store
	Mstore = mysql.NewDefaultStore(
		mysql.NewConfig("root:+mysql2016@tcp(192.168.43.179:3306)/myapp_test?charset=utf8"),
	)

	manager.MapTokenStorage(Mstore)

	//manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	//manager.MapAccessGenerate(generates.NewAccessGenerate())

	ClientStore = store.NewClientStore()

	ClientStore.Set("client", &models.Client{
		ID:     "client",
		Secret: "123456",
		//Domain: domainvar,
	})

	manager.MapClientStorage(ClientStore)

	OAuthSrv = server.NewServer(server.NewConfig(), manager)

	//OAuthSrv.SetUserAuthorizationHandler(userAuthorizeHandler)
	OAuthSrv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {

		sDec, err := base64.StdEncoding.DecodeString(username)
		if err != nil {
			return "", err
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
		//fmt.Println(ucList)
		// 2 用 username and client 查询userid,   es or  mysql
		var u modules.User
		modules.MysqlDb.Table("user").Select("*").Where(" name=? and  client=?", username, client).Find(&u)
		userID = strconv.FormatInt(u.Id, 10)
		return
	})
}

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
