package modules

import (
	"encoding/json"
)

var (
	Res      *PostResult = &PostResult{}
	HtmlData             = &Htmldata{}
)

type Htmldata struct {
	CommCell        []map[string]interface{}
	TotalData       []map[string]interface{}
	CommCellTotal   int
	CommCellVisible bool
}

type BaseReqArgs struct { //请求commcell参数
	StartTime interface{} `json:"START TIME"`
	Page      int         `json:"page"`
	PageSize  int         `json:"pagesize"`
}

type RequestArgs struct {
	*BaseReqArgs
	RunType   interface{} `json:"RUNTYPE"`
	SolveType interface{} `json:"SOLVETYPE"`
}




type HtmlGetData interface {
	GetCommCells() bool
	GetTotalData() bool
	DownFile() bool
	SaveStartTime(currentUnix int64)
	JsonConvert() []byte
}

type EmailGenSend struct {
	HtmlGetData
	SendMailOk bool
}

type PostResult struct {
	Data PostResultData `json:"data"`
}
type PostResultData struct {
	Total    int                      `json:"total"`
	Datalist []map[string]interface{} `json:"dataList"`
}

func (this *RequestArgs) SaveStartTime(currentUnix int64) {

	//this.StartTime = map[string]string{"key": "gt", "value": strconv.FormatInt((currentUnix-int64(g.GetConfig().ValidHours*3600))*1000, 10)}
	tmp:=map[string]interface{}{"st":1592236800000,"et":1592928000000}
	this.StartTime = map[string]interface{}{"key": "between", "value": tmp}
}

func (this *RequestArgs) JsonConvert() (b []byte) {
	b, _ = json.Marshal(this)
	return
}
