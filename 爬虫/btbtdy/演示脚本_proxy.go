package main

import (
	"fmt"
	"github.com/tebeka/selenium"
)

// eee 必须安装firefox 正式版 才能运行
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
		//selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		//selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}

	// 使用proxy前添加 xxx export https_proxy=http://127.0.0.1:1081
	// xxx export http_proxy=http://127.0.0.1:1081
	// xxx eport all _proxy=http://127.0.0.1:1081
	caps.AddProxy(selenium.Proxy{Type: selenium.System})

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("https://www.google.com"); err != nil {
		panic(err)
	}

	// Get a reference to the text box containing code.
	elem, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[2]/div[2]/form/div[2]/div[1]/div[3]/center/input[1]")
	if err != nil {
		panic(err)
	}
	fmt.Println(elem.GetAttribute("value"))
	// Remove the boilerplate code already in the text box.

}
