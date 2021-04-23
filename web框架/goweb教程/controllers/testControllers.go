package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  mysqlSearch
 * @Version: 1.0.0
 * @Date: 2021/4/21 下午2:06
 */

func registerStreamControllers() {
	http.Handle("/stream", http.HandlerFunc(streamView))
	http.Handle("/auth", http.HandlerFunc(authView))
	http.Handle("/timeout", http.HandlerFunc(timeoutView))
}

type Home struct {
	Username string
	Password string
}

func streamView(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body) // xxx 相比 json.Marshal 适合 stream 解析
	var h *Home = &Home{}
	err := dec.Decode(h)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
	ResponseJsonDataSuccess(w, r, h)

}

func authView(writer http.ResponseWriter, request *http.Request) {

	ResponseJsonSuccess(writer, request)
}

func timeoutView(writer http.ResponseWriter, request *http.Request) {
	time.Sleep(3 * time.Second)
	ResponseJsonSuccess(writer, request)
}
