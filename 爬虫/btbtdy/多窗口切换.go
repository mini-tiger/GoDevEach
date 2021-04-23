package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	//设置常量 分别设置chromedriver.exe的地址和本地调用端口
	Port = 9515

	seleniumPath = "/home/go/GoDevEach/爬虫/btbtdy/selenium-server-standalone-3.141.59.jar"
)

var geckoDriverPath = "/home/go/GoDevEach/爬虫/btbtdy/geckodriver"
var Webservice *selenium.Service
var Caps selenium.Capabilities

func init() {

	currentPath, _ := os.Getwd()
	os.Chdir(currentPath)
	if runtime.GOOS == "windows" {
		geckoDriverPath = geckoDriverPath + ".exe"
	}

	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		//selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(false)
	var err error
	Webservice, err = selenium.NewSeleniumService(seleniumPath, Port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
		//return err
	}
	//defer Webservice.Stop()
	// Connect to the WebDriver instance running locally.
	//Caps = selenium.Capabilities{"browserName": "firefox"}

	if runtime.GOOS == "linux" {
		Caps = selenium.Capabilities{}
		//Caps := selenium.Capabilities{"browserName": "firefox"}
		firefox_caps := firefox.Capabilities{
			Binary: "/root/firefox/firefox", // xxx 必须使用 开发版
			//Args:    []string{"--devtools"},
			Args:    nil,
			Profile: "",
			Log:     &firefox.Log{Level: firefox.Info},
			Prefs:   nil,
		}
		err := firefox_caps.SetProfile("/root/.mozilla/firefox/g5ccft6j.firefoxDev85.0b9_linux")
		if err != nil {
			log.Fatalln(err)
		}
		Caps.AddFirefox(firefox_caps)

	} else {
		Caps = selenium.Capabilities{"browserName": "firefox"} // xxx firefox 正常版此配置即可

		// xxx firefox-dev
		firefox_caps := firefox.Capabilities{
			Binary:  "C:\\Program Files\\Firefox Developer Edition\\firefox.exe",
			Args:    []string{"--devtools"},
			Profile: "",
			Log:     nil,
			Prefs:   nil,
		}
		//err:=firefox_caps.SetProfile("C:\\work\\go-dev\\GoDevEach\\爬虫\\3cdjytzk.firefox_win_profile")
		//if err!=nil{
		//	log.Fatalln(err)
		//}
		Caps.AddFirefox(firefox_caps)
	}

	//Caps.SetLogLevel(log.Type("driver"), log.Level("Info"))

	// eee 代理使用方法 未确认 ,linux 通过profile 配置
	//Caps.AddProxy(selenium.Proxy{Type: selenium.System})

}

func main() {

	//延迟关闭服务
	defer Webservice.Stop()

	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(Caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", Port))
	if err != nil {
		panic(err)
	}
	//延迟退出chrome
	defer wd.Quit()

	//3.打开多页面chrome实例
	//目前就想到两种方式可以打开，
	//第一种就是页面中有url连接，通过click（）方式打开,不推荐

	//第二种方式就是通过脚本方式打开。wd.ExecuteScript,不推荐
	if err := wd.Get("http://cdn1.python3.vip/files/selenium/sample3.html"); err != nil {
		panic(err)
	}

	//第一种方式，找到页面中的url地址，进行页面跳转
	we, err := wd.FindElement(selenium.ByTagName, "a")
	if err != nil {
		panic(err)
	}
	we.Click()

	//第二种方式，通过运行通用的js脚本打开新窗口，因为我们暂时不需要操作获取的结果，所有不获取返回值。
	wd.ExecuteScript(`window.open("https://www.qq.com", "_blank");`, nil)
	wd.ExecuteScript(`window.open("https://www.runoob.com/jsref/obj-window.html", "_blank");`, nil)

	//这一行是发送警报信息，写这一行的目的，主要是看当前主窗口是哪一个
	//wd.ExecuteScript(`window.alert(location.href);`, nil) xxx 执行此行 不能获取当前页面url

	//查看当前窗口的handle值
	//handle, err := wd.CurrentWindowHandle()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(handle)
	fmt.Println("-------遍历所有Tab-------------------")

	//查看所有网页的handle值
	handles, err := wd.WindowHandles()
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	for _, handle := range handles {
		wd.SwitchWindow(handle)
		url, _ := wd.CurrentURL()
		title, _ := wd.Title()
		fmt.Printf("当前页面handle:%s,url:%s,title:%s\n", handle, url, title)
	}

	//4.跳转到指定的网页
	//我们虽然打开了多个页面，但是我们当前的handle值，还是第一个页面的，我们要想办法搞定它。
	//记得保存当前主页面的handle值
	//mainhandle := handle

	//通过判断条件进行相应的网页
	//获取所有handle值
	handles, err = wd.WindowHandles()
	if err != nil {
		panic(err)
	}
	fmt.Println("--------定位qq.com元素------------------")
	//xxx 遍历所有handle值，通过url找到目标页面，判断相等时，break出来，就是停到相应的页面了。
	for _, handle := range handles {
		wd.SwitchWindow(handle)
		url, _ := wd.CurrentURL()
		if strings.Contains(url, "qq.com") {
			break
		}
	}
	//查看此页面的handle
	handle, err := wd.CurrentWindowHandle()
	if err != nil {
		panic(err)
	}
	fmt.Println(handle)

	elem, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[1]/div[2]/h1/a/img")
	if err != nil {
		panic(err)
	}
	fmt.Println(elem.GetAttribute("alt"))
	//这一行是发送警报信息，写这一行的目的，主要是看当前主窗口是哪一个
	wd.ExecuteScript(`window.alert(location.href);`, nil)
	//切换回第一个页面
	//wd.SwitchWindow(mainhandle)

	//睡眠20秒后退出
	time.Sleep(10 * time.Second)
}
