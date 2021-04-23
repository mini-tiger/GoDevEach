package main

import (
	"fmt"
	"github.com/tebeka/selenium/firefox"
	"log"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

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

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("http://play.golang.org/?simple=1"); err != nil {
		panic(err)
	}

	// Get a reference to the text box containing code.
	elem, err := wd.FindElement(selenium.ByCSSSelector, "#code")
	if err != nil {
		panic(err)
	}
	// Remove the boilerplate code already in the text box.
	if err := elem.Clear(); err != nil {
		panic(err)
	}

	// Enter some new code in text box.
	err = elem.SendKeys(`
		package main
		import "fmt"
		func main() {
			fmt.Println("Hello WebDriver!\n")
		}
	`)
	if err != nil {
		panic(err)
	}

	// Click the run button.
	btn, err := wd.FindElement(selenium.ByCSSSelector, "#run")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}

	// Wait for the program to finish running and get the output.
	outputDiv, err := wd.FindElement(selenium.ByCSSSelector, "#output")
	if err != nil {
		panic(err)
	}

	var output string
	for {
		output, err = outputDiv.Text()
		if err != nil {
			panic(err)
		}
		if output != "Waiting for remote server..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))

	err = wd.Get(" http://btbtdy3.com/down/34845-0-0.html")
	fmt.Println(err)

	// Example Output:
	// Hello WebDriver!
	//
	// Program exited.
}
