package g

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"log"
	"os"
	"runtime"
	"sync"
)

const (
	// These paths will be different on your system.
	seleniumPath = "selenium-server-standalone-3.141.59.jar"

	Port     = 8081
	Urlbasic = "http://btbtdy3.com/down/34845-0-%d.html"
)

var geckoDriverPath = "geckodriver"

var Webservice *selenium.Service
var Caps selenium.Capabilities
var Wg sync.WaitGroup

//var Sema *nsema.Semaphore = nsema.NewSemaphore(1)
var Result map[int]interface{} = make(map[int]interface{})
var ExitChan chan struct{} = make(chan struct{}, 0)
var Wd selenium.WebDriver

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
	Wd, _ = selenium.NewRemote(Caps, fmt.Sprintf("http://localhost:%d/wd/hub", Port))
	// eee 代理使用方法 未确认 ,linux 通过profile 配置
	//Caps.AddProxy(selenium.Proxy{Type: selenium.System})

	//fmt.Println(Webservice)
}
