package main

import (
	"btbtdy/g"
	"btbtdy/modules"
	"fmt"
	"log"
	"sort"
	"time"
)

func GetNums(min int, max int) (nums []int) {
	nums = make([]int, 0)
	if min >= max {
		return
	}

	for i := min; i <= max; i++ {
		nums = append(nums, i)
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
	print("==================================================\n")
	for _, value := range keys {
		fmt.Println(g.Result[value])
	}

}

func main() {

	go func() {
		<-g.ExitChan
		fmt.Println(1)
		<-g.ExitChan
		fmt.Println(2)
		// xxx 两个错误 停止操作
		PrintResult()
		log.Fatalf("两次错误 停止!\n")
		//os.Exit(1)
	}()

	for _, value := range GetNums(17, 23) {
		//g.Wg.Add(1)
		g.Wd.ExecuteScript(fmt.Sprintf(`window.open("http://btbtdy3.com/down/34845-0-%d.html", "_blank");`, value), nil)
		//modules.ParseHtml(value)
	}
	fmt.Println("begin")
	time.Sleep(30 * time.Second)

	handles, err := g.Wd.WindowHandles()
	if err != nil {
		panic(err)
	}
	//time.Sleep(10 * time.Second)

	for _, handle := range handles {

		g.Wd.SwitchWindow(handle)
		url, _ := g.Wd.CurrentURL()
		//title, _ := wd.Title()
		if url == "about:blank" {
			g.Wd.ExecuteScript(fmt.Sprintf(`window.close()`), nil)
			continue
		}

		fmt.Printf("当前页面handle:%s,url:%s\n", handle, url)

		modules.Parse(url)

	}

	//g.Wg.Wait()
	g.Wd.Quit()
	_ = g.Webservice.Stop()
	PrintResult()
}
