package modules

import (
	"cejing/TaskMulti/g"
	"encoding/json"
	"sync"
)

type BaseReqArgs struct { //请求commcell参数
	Page     int `json:"page"`
	PageSize int `json:"pagesize"`
}

type ReqModelArgs struct {
	*BaseReqArgs
}

type ReqWellArgs struct {

}

type ReqWellQXWJMArgs struct {
	JH  string `json:"JH"`
}

type ReqQXWJMListArgs struct {
	QXWJM  string `json:"QXWJM"`
}


type ReqSaveTask1Args struct {
	Tasktitle string `json:"tasktitle"`
	Status string `json:"status"`
}

type ReqSaveTask2Args struct {

	Checkedlist []string `json:"checkedlist"`
	Model_config map[string]interface{} `json:"model_config"`
	Model_id interface{} `json:"model_id"`
	Qxwjm_list string `json:"qxwjm_list"`
	Qxwjmlist []string `json:"qxwjmlist"`
	Share string `json:"share"`
	Status string `json:"status"`
	Taskdec string `json:"taskdec"`
	Taskrole string `json:"taskrole"`
	Tasktitle string `json:"tasktitle"`
	Taskuser string `json:"taskuser"`
	Wellname string `json:"wellname"`
	Yqtbm string `json:"yqtbm"`
	Id interface{} `json:"_id"`
}

var ReqReqSaveTask1ArgsFree = sync.Pool{
	New: func() interface{} {
		return &ReqSaveTask1Args{}
	},
}


type PostSaveTaskResult struct {
	Data       map[string]interface{} `json:"data"`
	Status     int            `json:"status"`
	StatusText string         `json:"statusText"`
}



var ReqWellQXWJMArgsFree = sync.Pool{
	New: func() interface{} {
		return &ReqWellQXWJMArgs{}
	},
}

var ReqQXWJMListArgsFree = sync.Pool{
	New: func() interface{} {
		return &ReqQXWJMListArgs{}
	},
}

type HtmlGetData interface {
	GetUrl() string
	JsonConvert() []byte
}


type PostResult struct {
	Data       PostResData `json:"data"`
	Status     int            `json:"status"`
	StatusText string         `json:"statusText"`
}
type PostResData struct {
	Total    int                      `json:"total"`
	Datalist []map[string]interface{} `json:"dataList"`
}


type PostApiResult struct {
	Data       []string `json:"data"`
	Status     int            `json:"status"`
	StatusText string         `json:"statusText"`
}


var PostResFree = sync.Pool{
	New: func() interface{} {
		return &PostResult{}
	},
}



func (this *ReqModelArgs) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqModelArgs) GetUrl() string {
	return g.ReqModelUrl
}

func (this *ReqWellArgs) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqWellArgs) GetUrl() string {
	return g.ReqWellUrl
}

func (this *ReqWellQXWJMArgs) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqWellQXWJMArgs) GetUrl() string {
	return g.ReqWellQXWJMUrl
}


func (this *ReqQXWJMListArgs) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqQXWJMListArgs) GetUrl() string {
	return g.ReqQXWJMListUrl
}

func (this *ReqSaveTask1Args) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqSaveTask1Args) GetUrl() string {
	return g.ReqSaveTaskUrl
}

func (this *ReqSaveTask2Args) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqSaveTask2Args) GetUrl() string {
	return g.ReqSaveTaskUrl
}

