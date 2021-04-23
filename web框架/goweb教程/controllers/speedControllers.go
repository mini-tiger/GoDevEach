package controllers

import (
	"github.com/tidwall/gjson"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  mysqlSearch
 * @Version: 1.0.0
 * @Date: 2021/4/21 下午2:06
 */

func registerSpeedControllers() {
	http.Handle("/speed", http.HandlerFunc(speedView))

}

func speedView(w http.ResponseWriter, r *http.Request) {
	//dec := json.NewDecoder(r.Body) // xxx 相比 json.Marshal 适合 stream 解析
	//var h *Home = &Home{}
	//err := dec.Decode(h)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte(err.Error()))
	//} else {
	//
	//	if err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		w.Write([]byte(err.Error()))
	//	}
	//}
	bytes := make([]byte, 1024)
	r.Body.Read(bytes)

	if !gjson.GetBytes(bytes, "user").Exists() || !gjson.GetBytes(bytes, "password").Exists() {
		ResponseJsonError(w, r, 400, "params err")
		return
	}

	u := gjson.GetBytes(bytes, "user")

	if u.String() != "taojun" {
		ResponseJsonError(w, r, 400, "user not taojun")
		return
	}

	ResponseJsonDataSuccess(w, r, "success")

}
