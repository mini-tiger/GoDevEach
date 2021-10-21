package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// xxx 接口 go test -v .

/* xxx 覆盖 扫描子目录下 *_test.go 并执行,执行后的覆盖率，例子只写了两个接口的测试
1 go test -v -covermode=count -coverprofile=coverage.out -coverpkg ./... ./...
2 go tool cover -html=coverage.out -o coverage.html
*/

// PostJson 根据特定请求uri和参数param，以Json形式传递参数，发起post请求返回响应
func PostJson(router *gin.Engine, uri string, param map[string]interface{}) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(param)

	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w

}

// PostFormData 根据特定请求uri和参数param
func PostFormData(router *gin.Engine, uri string, param map[string]string) *httptest.ResponseRecorder {
	data := url.Values{}
	for k, v := range param {
		data.Set(k, v)
	}

	//reader := strings.NewReader("user=taojun&password=1")
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w

}

func Get(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var router *gin.Engine

func init() {
	SetupCfg()
	router = SetupServer()
	SetupPlugins()
}

func TestHelloWorld(t *testing.T) {
	body := gin.H{
		"statusText": "Hello, World,TaoJun First Deploy for DataCenter",
	}

	w := Get(router, "GET", "/")

	assert.Equal(t, http.StatusOK, w.Code) // 返回是否200

	var response map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	value, exists := response["statusText"]

	assert.Nil(t, err)
	assert.True(t, exists)                     // 返回是否存在 statusText
	assert.Equal(t, body["statusText"], value) // statusText 的内容是否相等
}

func TestPostJson(t *testing.T) {
	body := gin.H{
		"statusText": "ok",
	}

	w := PostJson(router, "/login1", map[string]interface{}{"user": "taojun", "password": 1111})

	assert.Equal(t, http.StatusOK, w.Code) // 返回是否200

	var response map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	value, exists := response["statusText"]

	assert.Nil(t, err)
	assert.True(t, exists)                     // 返回是否存在 statusText
	assert.Equal(t, body["statusText"], value) // statusText 的内容是否相等
}

func TestPostData(t *testing.T) {
	body := gin.H{
		"statusText": "ok",
	}

	w := PostFormData(router, "/login3", map[string]string{"user": "taojun", "password": "1111"})

	assert.Equal(t, http.StatusOK, w.Code) // 返回是否200

	var response map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	value, exists := response["statusText"]

	assert.Nil(t, err)
	assert.True(t, exists)                     // 返回是否存在 statusText
	assert.Equal(t, body["statusText"], value) // statusText 的内容是否相等
}
