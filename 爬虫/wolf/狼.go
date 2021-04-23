package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gitee.com/taojun319/tjtools/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/toolkits/file"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// todo https://godoc.org/github.com/PuerkitoBio/goquery

type Config struct {
	MasterUrl        string `json:"master_url"`
	MasterDir        string `json:"master_dir"`
	Pages            int    `json:"pages"`
	MaxOld           int64  `json:"max_old"`
	ExistCover       bool   `json:"exist_cover"`
	UseProxy         bool   `json:"use_proxy"`
	ProxyUrl         string `json:"proxy_url"`
	DownloadInterval int64  `json:"downloadInterval"`
}

var MasterUrl, MasterDir, proxyUrl string
var PAGES int
var MaxOld, DownloadInterval int64
var ExistCover, useProxy bool

var tmpChanWeb chan struct{} = make(chan struct{}, PAGES) //主页退出 通道
var tmpChan chan struct{} = make(chan struct{}, 1)        //最后下载图片种子后 退出，通道
var now = time.Now()
var w *sync.WaitGroup = new(sync.WaitGroup)
var masterChan chan *goquery.Document = make(chan *goquery.Document, PAGES)
var userAgentSlice = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:2.0b13pre) Gecko/20110307 Firefox/4.0b13pre"}

//var	dir_url map[string]string =make(map[string]string)// 下载图片目录  第二层链接地址
var dirImageUrls map[string][]string = make(map[string][]string)   // 下载图片目录  链接地址
var dirTorrentUrls map[string][]string = make(map[string][]string) // 下载种子  链接地址

func DecodeToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(dst[:nDst]), nil
}

func checkTime(ts, baseFormat string) bool {
	//base_format := "2006-01-02 15:04"
	parseStrTime, _ := time.Parse(baseFormat, ts) // todo 字符串转时间
	todayStrTime, _ := time.Parse("2006-1-2", fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()))
	if todayStrTime.Unix() < MaxOld*86400+parseStrTime.Unix() {
		return true
	} else {
		return false
	}
}

func split(any string, old string) string {
	return strings.Replace(any, old, "", -1)
}

func imagesUrl(url string) (tmpSlice []string, tmpSlice1 []string) {
	dom := UrlDomGet(fmt.Sprintf("%s%s", MasterUrl, url))

	// 查找图片
	dom.Find("#read_tpc > img").Each(func(i int, s *goquery.Selection) {
		img, ok := s.Attr("src")
		//log.Println("11111111112222222222",img)
		if ok {
			tmpSlice = append(tmpSlice, strings.TrimSpace(img))
		}
	})
	dom.Find("#read_tpc > a > img").Each(func(i int, s *goquery.Selection) {
		img, ok := s.Attr("src")
		//log.Println("11111111113333333",img)
		if ok {
			tmpSlice = append(tmpSlice, strings.TrimSpace(img))
		}
	})
	dom.Find("#read_tpc > font > a > img").Each(func(i int, s *goquery.Selection) {
		img, ok := s.Attr("src")
		//log.Println("11111111113333333",img)
		if ok {
			tmpSlice = append(tmpSlice, strings.TrimSpace(img))
		}
	})
	dom.Find("#read_tpc > font > img").Each(func(i int, s *goquery.Selection) {
		img, ok := s.Attr("src")
		//log.Println("11111111113333333",img)
		if ok {
			tmpSlice = append(tmpSlice, strings.TrimSpace(img))
		}
	})
	//查找种子链接
	torrent, e := dom.Find("#main > form > div > table > tbody > tr.r_one > td > div[id] > a").Attr("href")
	if e {
		tmpSlice1 = append(tmpSlice1, fmt.Sprintf("%s%s", MasterUrl, strings.TrimSpace(torrent)))
		//log.Println("111111111144444444444",torrent)
	}

	w.Done()
	return
}
func ParsMasterWeb(dom *goquery.Document) {
	defer func() { //url  不能打开的 恢复机制
		if err := recover(); err != nil {
			//log.Printf("跳过err:%s \n",err)
			//panic(fmt.Sprintf("err:%s\n",url))
			w.Done()
		}
	}()
	//tmpImageUrl := make(map[string]string, 0)
	t_tbody_rows := dom.Find("#ajaxtable").Find("tbody:nth-child(2)").Find("tr[align]")
	//log.Printf("request url:%s, tbody math:%d 个",dom.Url,t_tbody_rows.Length())
	//dd,e:=t_tbody_rows.Html()
	//time.Sleep(time.Duration(2*time.Second))
	//fmt.Println("111111111111111",dd,e)

	t_tbody_rows.Each(func(i int, s *goquery.Selection) {
		childTime := s.Find("td:nth-child(3)").Find("div").Text()
		//v := sa.Find(" td:last-child>span").Text() //需要两种判断日期的DOM结构，是否span 下是否有title
		if checkTime(childTime, "2006-1-2") { //时间是否在范围内
			sa := s.Find("td[id]").Find(".subject")
			childHref, exist := sa.Attr("href")
			if exist {
				w.Add(1)
				dir_string, e := DecodeToGBK(sa.Text())

				if e != nil {
					log.Printf("目录中文名转换失败%s\n", e)
				}

				dir_string = strings.Split(childTime, " ")[0] + "_" + strings.Replace(dir_string, "/", "", -1)

				dir_string = split(strings.TrimSpace(dir_string), ":")
				dir_string = split(strings.TrimSpace(dir_string), "?")

				url_string := childHref
				log.Printf("开始解析url:%s 的图片和种子", url_string)
				dirImageUrls[dir_string], dirTorrentUrls[dir_string] = imagesUrl(url_string) //不加go 并发太大可能会503拒绝连接
			}
		}
	})
	w.Done()
}

func UnLinks() {
	for {
		select {
		case dom, ok := <-masterChan:
			if ok {
				w.Add(1)                 //
				tmpChanWeb <- struct{}{} // 阻塞 向下运行，
				go ParsMasterWeb(dom)    //解析每行 的地址，并通过地址 解析出 图片和种子地址
			} else {
				break
			}
		}
	}
}
func ForumGet() {

	for i := 1; i <= PAGES; i++ {
		url := fmt.Sprintf("%sthread-htm-fid-4-page-%d.html", MasterUrl, i) // 亚洲小格式
		log.Printf("亚洲小格式 request : %s\n", url)
		go func() {
			masterChan <- UrlDomGet(url)
		}()
	}
	for i := 1; i <= PAGES; i++ {
		url := fmt.Sprintf("%sthread-htm-fid-99-page-%d.html", MasterUrl, i) // 亚洲原创
		log.Printf("亚洲无码原创 request : %s\n", url)
		go func() {
			masterChan <- UrlDomGet(url)
		}()
	}

	for i := 1; i <= PAGES; i++ {
		url := fmt.Sprintf("%sthread-htm-fid-21-page-%d.html", MasterUrl, i) // 欧美
		log.Printf("欧美 request : %s\n", url)
		go func() {
			masterChan <- UrlDomGet(url)
		}()
	}
	for i := 1; i <= PAGES; i++ {
		url := fmt.Sprintf("%sthread-htm-fid-5-page-%d.html", MasterUrl, i) // 亚洲无码
		log.Printf("亚洲无码 request : %s\n", url)
		go func() {
			masterChan <- UrlDomGet(url)
		}()
	}
	for i := 1; i <= PAGES*4; i++ {
		<-tmpChanWeb
	}
	//time.Sleep(10 * time.Second)
	close(masterChan)
}

func UrlDomGet(url string) *goquery.Document {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, _ := http.NewRequest("GET", url, nil)

	//request.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	//request.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
	//request.Header.Set("Accept-Encoding","gzip,deflate,sdch")
	//request.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
	//request.Header.Set("Cache-Control","max-age=0")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", userAgentSlice[rand.Intn(len(userAgentSlice))])

	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Printf("跳过url:%s,err:%s \n", err, url)
	//		//panic(fmt.Sprintf("err:%s\n",url))
	//	}
	//}()

	response, err := client.Do(request)

	if err != nil {
		log.Printf("[Error]:request:%s, url:%s", err, url)
		return nil
	}

	//if response.StatusCode == 200 {
	//	body, _ := ioutil.ReadAll(response.Body)
	//	fmt.Println(string(body))
	//}

	if response.StatusCode != 200 {
		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		log.Printf("url :%s,status code error: %d %s", url, response.StatusCode, response.Status)
		return nil
	}

	// Load the HTML document
	dom, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		//log.Fatal(err)
		log.Printf("[Error]:response:%s, url:%s", err, url)
	}
	//fmt.Println(dom)
	return dom
}

func DownFile(url, fp string, wdownload *sync.WaitGroup, tmpUseProxy bool) {

	defer func() {
		//c <- struct{}{}
		wdownload.Done()
	}()

	client := &http.Client{}
	if useProxy || tmpUseProxy {
		proxy := func(_ *http.Request) (*neturl.URL, error) {
			return neturl.Parse(proxyUrl)
		}

		transport := &http.Transport{Proxy: proxy, TLSHandshakeTimeout: 30 * time.Second}

		client = &http.Client{Transport: transport, Timeout: 60 * time.Second}
	} else {
		client = &http.Client{
			Timeout: 60 * time.Second,
		}
	}

	request, _ := http.NewRequest("GET", url, nil)

	//request.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	//request.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
	//request.Header.Set("Accept-Encoding","gzip,deflate,sdch")
	//request.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
	request.Header.Set("Cache-Control", "max-age=0")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", userAgentSlice[rand.Intn(len(userAgentSlice))])

	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Printf("跳过url:%s,err:%s \n", err, url)
	//		c <- struct{}{}
	//		//panic(fmt.Sprintf("err:%s\n",url))
	//
	//	}
	//}()

	response, err := client.Do(request)
	if err != nil {
		if strings.Contains(fp, "torrent") {
			log.Printf("[Error]:种子请求失败%s, url:%s", err, url)
			//DownFile(url, fp, wdownload, true)
		}
		if strings.Contains(fp, "jpg") {
			log.Printf("[Error]:图片请求失败%s, url:%s", err, url)
			//DownFile(url, fp, wdownload, true)
		}

		//c <- struct{}{}
		//wdownload.Done()
		return
	}

	if response.StatusCode != 200 {
		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		log.Printf("status code error: %d %s", response.StatusCode, response.Status)
		//c <- struct{}{}
		//wdownload.Done()
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//f := filepath.Dir(fp)
		log.Printf("Download 失败 body err: %s,file: %s, url:%s\n", err.Error(), fp, url)
		//c <- struct{}{}
		//wdownload.Done()
		if tmpUseProxy || useProxy { // 如果没有使用过 代码下载， 重新使用代码下载一次
			return
		} else {
			wdownload.Add(1)
			DownFile(url, fp, wdownload, true)
		}
		return
	}

	//fp := string(filepath.Join("c:\\", "1"))

	err = ioutil.WriteFile(fp, body, 0777)
	if err != nil {
		log.Printf("Download 失败 %v fp:[%v]\n", err.Error(), fp)
		//c <- struct{}{}
		//wdownload.Done()
		return
	}
	log.Printf("Download 成功: %+v\n", fp)
	//c <- struct{}{}

	//wdownload.Done()
}

func makedir() {
	for k, _ := range dirImageUrls {
		err := os.MkdirAll(filepath.Join(MasterDir, k), 0777)
		if err != nil {
			log.Printf("mkdir Error DIR: %s, err:%s", k, err)
			//time.Sleep(100 * time.Second)
		}
	}
}

//func exist(file string) bool {
//
//	if _, err := os.Stat(file); err != nil {
//		if os.IsNotExist(err) {
//			//fmt.Printf("文件: %s 不存在\n", file)
//			return false
//		}
//	} else {
//		//fmt.Printf("文件: %s 存在\n", file)
//		return true
//	}
//
//	return false
//}
func jiange() {
	time.Sleep(time.Duration(DownloadInterval) * time.Millisecond)
}
func downloadall() {
	makedir()
	imgIndex := 1
	var wDownload *sync.WaitGroup = new(sync.WaitGroup)
	for k, v := range dirImageUrls { //一级目录
		log.Printf("正在下载第 %d 个页面的图片，共有%d个页面\n", imgIndex, len(dirImageUrls))
		//tmpC := make(chan struct{}, len(v)) //控制下载 并发，每个一级子目录 下的图片为一次并发
		mDir := filepath.Join(filepath.Join(MasterDir, k))
		for i := 0; i < len(v); i++ {
			//u:=fmt.Sprintf("%s%s", MasterUrl, v[i])
			tmpFile := fmt.Sprintf("%d.jpg", i)
			wDownload.Add(1)
			if utils.Exist(filepath.Join(mDir, tmpFile)) {
				if ExistCover { //存在文件 且常量定义为覆盖，则覆盖
					go DownFile(v[i], filepath.Join(mDir, tmpFile), wDownload, false)
					jiange() // todo 下载间隔
				} else {
					//log.Printf("file:%s 跳过", filepath.Join(mDir, tmpFile))
					//tmpC <- struct{}{}
					wDownload.Done()
					continue
				}
			} else {
				go DownFile(v[i], filepath.Join(mDir, tmpFile), wDownload, false)
				jiange() // todo 下载间隔
			}
		}
		//for i := 0; i < len(v); i++ { //控制并发
		//	<-tmpC
		//}
		imgIndex = imgIndex + 1

	}
	torrentIndex := 1
	for k, v := range dirTorrentUrls {
		log.Printf("正在下载第 %d 个页面的种子，共有%d个页面\n", torrentIndex, len(dirTorrentUrls))
		//tmpC := make(chan struct{}, len(v))
		mDir := filepath.Join(filepath.Join(MasterDir, k))
		for i := 0; i < len(v); i++ {
			//u:=fmt.Sprintf("%s%s", MasterUrl, v[i])
			wDownload.Add(1)
			tmp_file := fmt.Sprintf("%d.torrent", i)
			if utils.Exist(filepath.Join(mDir, tmp_file)) {
				if ExistCover { //存在且常量定义为覆盖，覆盖
					go DownFile(v[i], filepath.Join(mDir, tmp_file), wDownload, false)
					jiange() // todo 下载间隔
				} else {
					//log.Printf("file:%s 跳过", filepath.Join(mDir, tmp_file))
					//tmpC <- struct{}{}
					wDownload.Done()
					continue
				}
			} else {
				go DownFile(v[i], filepath.Join(mDir, tmp_file), wDownload, false)
				jiange() // todo 下载间隔
			}
		}
		//for i := 0; i < len(v); i++ {
		//	<-tmpC
		//}
		torrentIndex = torrentIndex + 1
	}
	wDownload.Wait()
	tmpChan <- struct{}{}
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {

		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)

	}

	var c Config
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}
	MasterUrl = c.MasterUrl
	MasterDir = c.MasterDir
	PAGES = c.Pages
	proxyUrl = c.ProxyUrl
	ExistCover = c.ExistCover
	MaxOld = c.MaxOld
	useProxy = c.UseProxy
	DownloadInterval = c.DownloadInterval

	log.Println("read config file:", cfg, "successfully")
	//WLog(fmt.Sprintf("read config file: %s successfully",cfg))
}

func SetupCfg() {
	_, filename, _, _ := runtime.Caller(0)
	devJson := filepath.Join(filepath.Dir(filename), "cfg.json")

	//ParseConfig("cfg.json") //
	ParseConfig(devJson)
	MasterUrlCustom := flag.String("url", "", "url")
	//UseProxyCustom := flag.Bool("proxy", false, "proxy") // 只要在命令行 写入 proxy 就是true
	maxold := flag.Int64("maxold", 0, "MaxOld")

	flag.Parse() // todo 优先使用命令行参数
	//log.Println(*MasterUrlCustom, *maxold, *UseProxyCustom)
	if *maxold != 0 {
		MaxOld = *maxold
	}
	//if !*UseProxyCustom {
	//	useProxy = false
	//}
	if *MasterUrlCustom != "" {
		MasterUrl = *MasterUrlCustom
	}
	log.Printf("========配置参数")
	log.Println("最大天数:", MaxOld, "页面最多:", PAGES, "useproxy:", useProxy, "MasteruRL:", MasterUrl)
}

func main() {

	// todo 读取并配置参数
	SetupCfg()

	time.Sleep(time.Duration(1) * time.Second)

	_now := time.Now().Unix()
	runtime.GOMAXPROCS(1)
	go UnLinks()
	ForumGet()
	log.Printf("当前进程PID:%d\n", os.Getpid())
	w.Wait()

	downloadall()

	<-tmpChan
	fmt.Printf("总共用时%d秒", time.Now().Unix()-_now)
}
