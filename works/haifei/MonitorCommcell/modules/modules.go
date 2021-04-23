package modules

import (
	"encoding/json"
	"haifei/MonitorCommcell/g"
	"sync"
)

type BaseReqArgs struct { //请求commcell参数
	Page     int `json:"page"`
	PageSize int `json:"pagesize"`
}

type ReqCCArrArgs struct {
	*BaseReqArgs
}

type ReqCCArgs struct {
	*BaseReqArgs
	ReportTime interface{} `json:"REPORTTIME"`
	CommCell   string      `json:"COMMCELL"`
}

var ReqCCArgsFree = sync.Pool{
	New: func() interface{} {
		return &ReqCCArgs{}
	},
}

type ReqInsHisCC struct {
	CommCell   string      `json:"COMMCELL"`
	ReportTime interface{} `json:"REPORTTIME"`
	EXIST      string      `json:"EXIST"`
	NormTime   int64       `json:"NORMTIME"`
}

var ReqInsHisCCFree = sync.Pool{
	New: func() interface{} {
		return &ReqInsHisCC{}
	},
}

type HtmlGetData interface {
	GetUrl() string
	JsonConvert() []byte
}

//type EmailGenSend struct {
//	HtmlGetData
//	SendMailOk bool
//}

type PostResult struct {
	Data       PostResData `json:"data"`
	Status     int            `json:"status"`
	StatusText string         `json:"statusText"`
}
type PostResData struct {
	Total    int                      `json:"total"`
	Datalist []map[string]interface{} `json:"dataList"`
}

var PostResFree = sync.Pool{
	New: func() interface{} {
		return &PostResult{}
	},
}

func GenReqTime() {

	g.ReqNormTime = g.CurrWorkTime - (int64(g.GetConfig().SpaceTime*60) * 1000)

	g.ReqReportTime = map[string]interface{}{"key": "between", "value":
	map[string]interface{}{
		"st": g.ReqNormTime - int64(g.GetConfig().ValidTime*60*1000),
		"et": g.ReqNormTime + int64(g.GetConfig().ValidTime*60*1000)},
	}
}

func (this *ReqCCArrArgs) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqCCArrArgs) GetUrl() string {
	return g.ReqCCUrl
}

func (this *ReqCCArgs) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqCCArgs) GetUrl() string {
	return g.ReqCCArrUrl
}

func (this *ReqInsHisCC) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}

func (this *ReqInsHisCC) GetUrl() string {
	return g.ReqInsHisUrl
}
