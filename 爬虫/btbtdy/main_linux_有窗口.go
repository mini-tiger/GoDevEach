package main

import (
	"btbtdy/g"
	"fmt"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"log"
	"sort"
)

// xxx 被占用的端口 无窗口，反之 有 窗口
var MaxGoroutine int = 2
var finlishGoroutine chan struct{} = make(chan struct{}, 0)
var wdslice []selenium.WebDriver = make([]selenium.WebDriver, 2)

func GetNums(min int, max int, step int) (nums [][]int) {
	nums = make([][]int, 0)
	if min >= max {
		return
	}
	allnums := make([]int, 0)
	for i := min; i <= max; i++ {
		allnums = append(allnums, i)
	}
	for i := 0; i < len(allnums); i = i + step {
		var tmpnum []int
		if i+step >= len(allnums) {
			tmpnum = allnums[i:len(allnums)]
		} else {
			tmpnum = allnums[i : i+step]
		}
		nums = append(nums, tmpnum)
	}

	return nums
}

func PrintResult() {
	var keys []int
	for k := range g.Result {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, value := range keys {
		fmt.Println(value, g.Result[value])
	}

}
func main() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	const (
		// These paths will be different on your system.
		seleniumPath    = "/home/go/GoDevEach/爬虫/btbtdy/selenium-server-standalone-3.141.59.jar"
		geckoDriverPath = "/home/go/GoDevEach/爬虫/btbtdy/geckodriver"
		port            = 8081
	)
	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(),           // xxx 开启则必须启动xvfb Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		//selenium.Output(os.Stderr),            // eee 输出日志 Output debug information to STDERR.
	}

	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}

	firefox_caps := firefox.Capabilities{
		Binary: "/root/firefox/firefox", // xxx 必须使用 开发版
		//Args:    []string{"--devtools"},
		Args:    nil,
		Profile: "",
		Log:     &firefox.Log{Level: firefox.Info},
		Prefs:   nil,
	}
	err = firefox_caps.SetProfile("/root/.mozilla/firefox/g5ccft6j.firefoxDev85.0b9_linux")
	if err != nil {
		log.Fatalln(err)
	}
	caps.AddFirefox(firefox_caps)

	for i := 0; i < len(wdslice); i++ {
		wdslice[i], _ = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	}

	defer func() {
		for _, wd := range wdslice {
			wd.Quit()
		}
	}()

	//if err != nil {
	//	panic(err)
	//}

	var Title, DownUrl string

	for _, subslice := range GetNums(0, 23, MaxGoroutine) {

		for index, value := range subslice {
			go func(index, value int) {
				defer func() {
					finlishGoroutine <- struct{}{}
				}()
				if err = wdslice[index].Get(fmt.Sprintf("http://btbtdy3.com/down/34845-0-%d.html", value)); err != nil {
					log.Fatalln(err)
				}

				elem, err := wdslice[index].FindElement(selenium.ByID, "video-down")
				listelem, err := elem.FindElements(selenium.ByTagName, "p")

				if err != nil {
					log.Printf("video-down -> p 失败:%v,跳过\n", err)
				}

				if len(listelem) == 2 {
					Title, _ = listelem[0].Text()
					DownUrl, _ = listelem[1].Text()

				}
				fmt.Printf("title:%s,downurl:%s\n", Title, DownUrl)
				time.Sleep(5 * time.Second)
			}(index, value)

		}
		fmt.Println(1)
		for _ = range subslice {
			<-finlishGoroutine
		}

		fmt.Println(2)
	}

}
