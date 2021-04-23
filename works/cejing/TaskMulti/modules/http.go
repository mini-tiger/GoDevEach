package modules

import (
	"bytes"
	"cejing/TaskMulti/g"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//type TotalRow struct {
//	BIANMA   string
//	COMMCELL string
//	CLIENT   string
//}

//var PostResultFree = sync.Pool{
//	New: func() interface{} {
//		return &PostResult{}
//	},
//}


func GetHtmlDataHandler(this HtmlGetData,req string) interface{}{

	var Res interface{}

	switch req {
	case "crud":
		Res = PostResFree.Get().(*PostResult)
		break
	case "api":
		Res = &PostApiResult{}
		break
	case "savetask":
		Res = &PostSaveTaskResult{}
	}


	g.Logge.Debugf("开始 请求Url:%s", this.GetUrl())
	//var Res =&PostResult{}
	//p := PostCommCellJson{StartTime: this.RequestArgs.StartTime, Page: this.RequestArgs.Page, PageSize: this.RequestArgs.PageSize}

	//jsonBytes, _ := json.Marshal(this.BaseReqArgs)
	//g.Logge.Debug("请求总COMMCELL参数:%+v\n", string(jsonBytes))

	request, _ := http.NewRequest("POST", this.GetUrl(), bytes.NewBuffer(this.JsonConvert()))


	//request.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	//request.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
	//request.Header.Set("Accept-Encoding","gzip,deflate,sdch")
	//request.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
	//request.Header.Set("Cache-Control","max-age=0")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("content-type", "application/json;charset=UTF-8")

	response, err := g.HttpClient.Do(request)

	if err != nil {
		//g.Logge.Errorf("请求失败:  %s", err)
		panic(err)
		//return Res
	}

	if response.StatusCode != 200 {
		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		g.Logge.Errorf("status code error: %d %s", response.StatusCode, response.Status)

		return Res
	}
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Printf("body:%+v\n", string(body))
	//var res *PostResult = PostResultFree.Get().(*PostResult)
	//defer PostResultFree.Put(res)

	_ = json.Unmarshal(body, Res)
	//defer utils.Clear(Res)
	//fmt.Printf("body res:%+v\n", Res)

	//HtmlData.CommCell = Res.Data.Datalist
	//HtmlData.CommCellTotal = Res.Data.Total
	g.Logge.Debugf("请求成功")
	return Res
}

//
//func (this *RequestArgs) GetTotalData() bool {
//	g.Logge.Debug("开始 请求TotalData\n")
//
//	//client := &http.Client{}
//	//
//	//client = &http.Client{
//	//	Timeout: 60 * time.Second,
//	//}
//	//
//
//	request, _ := http.NewRequest("POST", g.RequestTotalUrl, bytes.NewBuffer(g.JsonBytes))
//
//	//request.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
//	//request.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
//	//request.Header.Set("Accept-Encoding","gzip,deflate,sdch")
//	//request.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
//	//request.Header.Set("Cache-Control","max-age=0")
//	request.Header.Set("Connection", "keep-alive")
//	request.Header.Set("content-type", "application/json;charset=UTF-8")
//
//	response, err := g.HttpClient.Do(request)
//	if err != nil {
//		_ = g.Logge.Error("请求totaldata失败:  %s\n", err)
//
//		return false
//	}
//
//	if response.StatusCode != 200 {
//		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
//		_ = g.Logge.Error("status code error: %d %s\n", response.StatusCode, response.Status)
//
//		return false
//	}
//	body, _ := ioutil.ReadAll(response.Body)
//	//fmt.Printf("body:%+v\n", string(body))
//	//var res *PostResult = PostResultFree.Get().(*PostResult)
//	//defer PostResultFree.Put(res)
//
//	_ = json.Unmarshal(body, Res)
//	defer utils.Clear(Res)
//	//fmt.Printf("body res:%+v\n",res)
//	HtmlData.TotalData = Res.Data.Datalist
//	g.Logge.Debug("请求totaldata 成功\n")
//	return true
//}
//
//func (this *RequestArgs) DownFile() bool {
//
//	g.Logge.Debug("开始 download 明细数据 %s,\n请求参数:%s\n", g.Fp, string(g.JsonBytes))
//
//	//client := &http.Client{}
//	//
//	//client = &http.Client{
//	//	Timeout: 60 * time.Second,
//	//}
//
//	request, _ := http.NewRequest("POST", g.DownLoadUrl, bytes.NewBuffer(g.JsonBytes))
//
//	//request.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
//	//request.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
//	//request.Header.Set("Accept-Encoding","gzip,deflate,sdch")
//	//request.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
//	//request.Header.Set("Cache-Control","max-age=0")
//	request.Header.Set("Connection", "keep-alive")
//	request.Header.Set("content-type", "application/json;charset=UTF-8")
//	//request.Header.Set("User-Agent", userAgentSlice[rand.Intn(len(userAgentSlice))])
//
//	response, err := g.HttpClient.Do(request)
//	if err != nil {
//		_ = g.Logge.Error("下载失败:  %s\n", err)
//		//this.DownloadSuccess = false
//		return false
//	}
//
//	if response.StatusCode != 200 {
//		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
//		_ = g.Logge.Error("status code error: %d %s\n", response.StatusCode, response.Status)
//		//c <- struct{}{}
//		//this.DownloadSuccess = false
//		return false
//	}
//
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		f := filepath.Dir(g.Fp)
//		_ = g.Logge.Error("body err: %s,dir: %s, url:%s\n", err.Error(), f, g.DownLoadUrl)
//		//c <- struct{}{}
//		//this.DownloadSuccess = false
//		return false
//	}
//
//	//fp := string(filepath.Join("c:\\", "1"))
//
//	err = ioutil.WriteFile(g.Fp, body, 0777)
//	if err != nil {
//		_ = g.Logge.Error("%v fp:[%v]\n", err.Error(), g.Fp)
//
//		//this.DownloadSuccess = false
//		return false
//	}
//	g.Logge.Debug("Download 成功: %+v\n", g.Fp)
//	//this.DownloadSuccess = true
//	return true
//}
