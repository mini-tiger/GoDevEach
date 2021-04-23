package 桃花族

import (
	"encoding/json"
	"flag"
	"fmt"
	"gitee.com/taojun319/tjtools/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/toolkits/file"
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
	MasterUrl  string `json:"master_url"`
	MasterDir  string `json:"master_dir"`
	Pages      int    `json:"pages"`
	MaxOld     int64  `json:"max_old"`
	ExistCover bool   `json:"exist_cover"`
	UseProxy   bool   `json:"use_proxy"`
	ProxyUrl   string `json:"proxy_url"`
}

//const (
//	MasterUrl  = "http://thzbt.co/"
//	MasterDir  = "g:\\image\\"
//	PAGES      = 3     //最多看3页的数据，3
//	MaxOld     = 7     //最大几天前
//	ExistCover = false //存在是否覆盖
//	useProxy   = true  // 使用ssr翻墙，本地1080端口
//	proxyUrl   = "http://192.168.1.100:1080"
//)

var MasterUrl, MasterDir, proxyUrl string
var PAGES int
var MaxOld int64
var ExistCover, useProxy bool

var tmpChanWeb chan struct{} = make(chan struct{}, PAGES) //主页退出 通道
var tmpChan chan struct{} = make(chan struct{}, 1)        //最后下载图片种子后 退出，通道
var now = time.Now()
var w *sync.WaitGroup = new(sync.WaitGroup)
var masterChan chan *goquery.Document = make(chan *goquery.Document, PAGES)
var userAgentSlice = []string{
	"chrome 67",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:2.0b13pre) Gecko/20110307 Firefox/4.0b13pre"}

//var	dir_url map[string]string =make(map[string]string)// 下载图片目录  第二层链接地址
var dirImageUrls map[string][]string = make(map[string][]string)   // 下载图片目录  链接地址
var dirTorrentUrls map[string][]string = make(map[string][]string) // 下载种子  链接地址

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

func imagesUrl(url string) (tmpSlice []string, tmpSlice1 []string) {
	//tmpSlice := make([]string, 0)
	//tmpSlice1 := make([]string, 0)

	dom := UrlDomGet(fmt.Sprintf("%s%s", MasterUrl, url))
	td := dom.Find("table>tbody td.t_f") //todo 子元素选择器 不是直接上下级关系 的 中间有空格
	//td:=dom.Find("table>tbody>tr>td.t_f")//todo 子元素选择器 是直接上下级关系 的 > 号

	// 查找图片
	td.Find("img ").Each(func(i int, s *goquery.Selection) {
		img, ok := s.Attr("file")
		if ok {
			tmpSlice = append(tmpSlice, img)
		}
	})

	//查找种子链接
	td = dom.Find("p.attnm > a")
	td.Each(func(i int, s *goquery.Selection) {
		torrentParentUrl, ok := s.Attr("href")
		if ok {
			tmpDom := UrlDomGet(fmt.Sprintf("%s%s", MasterUrl, torrentParentUrl)) //再次请求种子页面
			initHref := tmpDom.Find("div.f_c div[style^=padding-left] a")         //在种子下载页面查找
			//fmt.Println(init_href.Html())
			torrent, ok := initHref.Attr("href")
			if ok {
				tmpSlice1 = append(tmpSlice1, torrent)
			}
		}
	})

	w.Done()
	return
}
func ParsMasterWeb(dom *goquery.Document) { //解析第一层主页
	defer func() { //url  不能打开的 恢复机制
		if err := recover(); err != nil {
			//log.Printf("跳过err:%s \n",err)
			//panic(fmt.Sprintf("err:%s\n",url))
			w.Done()
		}
	}()
	//tmpImageUrl := make(map[string]string, 0)
	t_tbody := dom.Find("table[summary]").Find("tbody")
	//log.Printf("request url:%s, tbody math:%d 个",dom.Url,t_tbody.Length())
	t_tbody.Each(func(i int, s *goquery.Selection) {
		sa := s.Find("tr>td.by").First().Find("em").Children()
		/*
			一种是
			<em><span><span title="2018-9-20">7&nbsp;天前</span></span></em>
			一种是
			<em><span>2018-9-13</span></em>
		*/
		if v, b := sa.Find("span").Attr("title"); b { //需要两种判断日期的DOM结构，是否span 下是否有title
			if checkTime(v, "2006-1-2") { //时间是否在范围内
				dir_string := s.Find("tr>th>a").Text()
				//fmt.Println(dir_string)
				url_string, _ := s.Find("tr>th>a.s.xst").Attr("href")
				//fmt.Println(url_string)
				w.Add(1)
				log.Printf("开始解析url:%s 的图片和种子", url_string)
				//tmpImageUrl[dir_string]=url_string
				//fmt.Println("------------",strings.Contains(dom.Url.String(),"forum-181"))
				dirImageUrls[dir_string], dirTorrentUrls[dir_string] = imagesUrl(url_string) //不加go 并发太大可能会503拒绝连接

			}
		} else {
			aa, _ := sa.Html()
			if checkTime(aa, "2006-1-2") {
				dir_string := s.Find("tr>th>a").Text()
				//fmt.Println(dir_string)
				url_string, _ := s.Find("tr>th>a.s.xst").Attr("href")
				//fmt.Println(url_string)
				w.Add(1)
				log.Printf("开始解析url:%s 的图片和种子", url_string)
				//tmpImageUrl[dir_string]=url_string
				dirImageUrls[dir_string], dirTorrentUrls[dir_string] = imagesUrl(url_string) //不加go 否则可能会503拒绝连接
			}
		}
	})
	//for dirString,urlString:=range tmpImageUrl{
	//	w.Add(1)
	//	go func() {
	//		dirImageUrls[dirString], dirTorrentUrls[dirString] = imagesUrl(urlString)
	//	}()
	//}
	w.Done()
}

func UnLinks() {
	for {
		select {
		case dom, ok := <-masterChan:
			if ok {
				w.Add(1)                 // todo 阻塞 二级以下页面解析
				tmpChanWeb <- struct{}{} //todo 阻塞一级页面解析
				go ParsMasterWeb(dom)
			} else {
				break
			}
		}
	}
}
func ForumGet() {

	for i := 1; i <= PAGES; i++ {
		url := fmt.Sprintf("%sforum-181-%d.html", MasterUrl, i)
		log.Printf("亚洲 request : %s\n", url)
		go func() {
			masterChan <- UrlDomGet(url)
		}()
	}
	for i := 1; i <= PAGES; i++ {
		url := fmt.Sprintf("%sforum-182-%d.html", MasterUrl, i)
		log.Printf("欧美 request : %s\n", url)
		go func() {
			masterChan <- UrlDomGet(url)
		}()
	}

	for i := 1; i <= PAGES*2; i++ {
		<-tmpChanWeb
	}
	//time.Sleep(10 * time.Second)
	close(masterChan)
}

func UrlDomGet(url string) *goquery.Document {

	client := &http.Client{
		Timeout: 15 * time.Second,
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
		log.Printf("[Error]:%s, url:%s", err, url)
	}

	//if response.StatusCode == 200 {
	//	body, _ := ioutil.ReadAll(response.Body)
	//	fmt.Println(string(body))
	//}

	if response.StatusCode != 200 {
		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		log.Printf("url :%s,status code error: %d %s", url, response.StatusCode, response.Status)
	}

	// Load the HTML document
	dom, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		//log.Fatal(err)
		log.Printf("[Error]:%s, url:%s", err, url)
	}
	//fmt.Println(dom)
	return dom

	//// Find the review items // 同级元素连接写， 父子级中间有空格
	//doc.Find("ul.BookList.PubIndex-recommends li").Each(func(i int, s *goquery.Selection) {
	//	// For each item found, get the band and title
	//	link, _ := s.Find("a").Attr("href") //获取属性值
	//	sonel := s.Find("a .Image")
	//	imagelink, _ := sonel.Attr("src")
	//	imagelinks = append(imagelinks, imagelink)
	//	//title, _ := sonel.Attr("alt")
	//	title := s.Find(".BookItem-title").Text() //获取文本
	//	fmt.Printf("NO.%d: link:%s - title:%-16s - imagelink:%-50s\n", i+1, link, title, imagelink)
	//})
	//for _, url := range imagelinks {
	//	downfile(url)
	//}
}

func DownFile(url, fp string, wdownload *sync.WaitGroup) {

	defer func() {
		//c <- struct{}{}
		wdownload.Done()
	}()

	log.Printf("开始 download %s,url:%s", fp, url)
	//resp, err := http.Get(url)
	//if err != nil {
	//	fmt.Printf("%v\n", err.Error())
	//	return
	//}
	//defer resp.Body.Close()
	client := &http.Client{}
	if useProxy {
		proxy := func(_ *http.Request) (*neturl.URL, error) {
			return neturl.Parse(proxyUrl)
		}

		transport := &http.Transport{Proxy: proxy}

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
	//request.Header.Set("Cache-Control","max-age=0")
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

		}
		if strings.Contains(fp, "jpg") {
			log.Printf("[Error]:图片请求失败%s, url:%s", err, url)

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
		f := filepath.Dir(fp)
		fmt.Printf("body err: %s,dir: %s, url:%s\n", err.Error(), f, url)
		//c <- struct{}{}
		//wdownload.Done()
		return
	}

	//fp := string(filepath.Join("c:\\", "1"))

	err = ioutil.WriteFile(fp, body, 0777)
	if err != nil {
		fmt.Printf("%v fp:[%v]\n", err.Error(), fp)
		//c <- struct{}{}
		//wdownload.Done()
		return
	}
	fmt.Printf("Download 成功: %+v\n", fp)
	//c <- struct{}{}
	//wdownload.Done()
}

func makedir() {
	for k, _ := range dirImageUrls {
		err := os.MkdirAll(filepath.Join(MasterDir, k), 0777)
		if err != nil {
			log.Println("mkdirerr %s err:%s", k, err)
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
					go DownFile(v[i], filepath.Join(mDir, tmpFile), wDownload)
				} else {
					log.Printf("file:%s 跳过", filepath.Join(mDir, tmpFile))
					//tmpC <- struct{}{}
					wDownload.Done()
					continue
				}
			} else {
				go DownFile(v[i], filepath.Join(mDir, tmpFile), wDownload)
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
					go DownFile(v[i], filepath.Join(mDir, tmp_file), wDownload)
				} else {
					log.Printf("file:%s 跳过", filepath.Join(mDir, tmp_file))
					//tmpC <- struct{}{}
					wDownload.Done()
					continue
				}
			} else {
				go DownFile(v[i], filepath.Join(mDir, tmp_file), wDownload)
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

	log.Println("read config file:", cfg, "successfully")
	//WLog(fmt.Sprintf("read config file: %s successfully",cfg))
}

func SetupCfg() {
	_, filename, _, _ := runtime.Caller(0)
	devJson := filepath.Join(filepath.Dir(filename), "cfg.json")

	//ParseConfig("cfg.json") //
	ParseConfig(devJson)
	MasterUrlCustom := flag.String("url", "", "url")
	UseProxyCustom := flag.Bool("proxy", false, "proxy") // 只要在命令行 写入 proxy 就是true
	maxold := flag.Int64("maxold", 0, "MaxOld")

	flag.Parse() // todo 优先使用命令行参数
	//log.Println(*MasterUrlCustom, *maxold, *UseProxyCustom)
	if *maxold != 0 {
		MaxOld = *maxold
	}
	if !*UseProxyCustom {
		useProxy = false
	}
	if *MasterUrlCustom != "" {
		MasterUrl = *MasterUrlCustom
	}
	log.Printf("========配置参数")
	log.Println("最大天数:", MaxOld, "页面最多:", PAGES, "useproxy:", useProxy, "MasteruRL:", MasterUrl)
}

func main() {

	// todo 读取并配置参数
	SetupCfg()

	time.Sleep(time.Duration(2) * time.Second)

	_now := time.Now().Unix()
	runtime.GOMAXPROCS(1)
	go UnLinks()
	ForumGet()
	log.Printf("当前进程PID:%d\n", os.Getpid())
	w.Wait()

	downloadall()

	<-tmpChan
	//for k, v := range dirImageUrls {
	//	fmt.Println(k, len(v))
	//}
	//fmt.Println(len(dirImageUrls))
	//for k, v := range dirTorrentUrls {
	//	fmt.Println(k, len(v))
	//}
	//fmt.Println(len(dirTorrentUrls))
	fmt.Printf("总共用时%d秒", time.Now().Unix()-_now)
}
