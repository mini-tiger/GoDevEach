package controllers

import (
	"encoding/json"
	"net/http"
)

type httpStatus struct {
	Statuscode int
	Msg        string
}

type httpDataResp struct {
	*httpStatus
	Data interface{}
}

func ResponseJsonSuccess(writer http.ResponseWriter, request *http.Request) {
	// Set response header
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(&httpStatus{200, "resp success"})
}

func ResponseJsonDataSuccess(writer http.ResponseWriter, request *http.Request, data interface{}) {
	// Set response header
	writer.Header().Set("Content-Type", "application/json")
	//fmt.Println(data)

	js, err := json.Marshal(&httpDataResp{&httpStatus{200, "resp success"}, data})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Write(js)
}

func ResponseJsonError(writer http.ResponseWriter, request *http.Request, statucode int, Msg string) {

	// Set response header
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statucode)
	_ = json.NewEncoder(writer).Encode(&httpStatus{statucode, Msg})
}
