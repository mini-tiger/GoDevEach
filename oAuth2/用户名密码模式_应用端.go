package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	authServerURL = "http://192.168.43.28:9096"
)

var (
	config = oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222", // xxx 这里要和 “用户名密码模式_认证服务器.go ” line 34 一样
		//Scopes:       []string{"all"},
		//RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			//AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func main() {

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		globalToken.Expiry = time.Now()
		token, err := config.TokenSource(context.Background(), globalToken).Token()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/pwd", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
		//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
		username := r.FormValue("username")
		//fmt.Println(username)
		password := r.FormValue("password")

		token, err := config.PasswordCredentialsToken(context.Background(), username, password)
		fmt.Println(config.Endpoint.AuthStyle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	log.Println("Client is running at 9094 port.Please open http://localhost:9094")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
