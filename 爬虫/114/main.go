package main

import (
	"btbtdy/g"
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"strings"
	"sync"
	"time"
)

var subkeshi_fuke selenium.WebElement
var beginurl = "https://www.114yygh.com/hospital/122/home"
var date = "02月01日"
var doctor = "郭永红"
var piaochan chan struct{} = make(chan struct{}, 0)
var wg sync.WaitGroup

func getxuanhao() {
	// 1. 打开页面
	g.Wd.Get(beginurl)

	// 2. 可以通过存储cookies 提取 免登录，一般30分钟过期
	//cookiesArr,err:=g.Wd.GetCookies()
	//if err!=nil{
	//	log.Panicf("cookies fail %s\n",err)
	//}
	//
	//for _,c := range cookiesArr{
	//	fmt.Printf("%+v\n",c)
	//}

	// 3. 点击登录按钮，需要手机微信登录
	time.Sleep(10 * time.Second)

	loginelem, err := g.Wd.FindElement(selenium.ByXPATH, "//*[@id=\"main\"]/div[1]/div/div[2]/span[2]")
	if err != nil {
		log.Panicf("login btn err:%s\n", err)
	}

	loginelem.Click()
	time.Sleep(10 * time.Second)

	// 4. 找到妇科并点击 跳转
	keshi_elem, err := g.Wd.FindElement(selenium.ByCSSSelector, "div.main-container div.el-scrollbar div.main-scroll.el-scrollbar__wrap:nth-child(1) div.el-scrollbar__view div.nav-container.page-component div.page-container div.hospital-home div.select-dept-wrapper > div.sub-dept-container")
	if err != nil {
		log.Panicf("科室class Fail err:%s\n", err)
	}
	subkeshi, err := keshi_elem.FindElements(selenium.ByCSSSelector, ".sub-dept-wrapper")
	if err != nil {
		log.Panicf("子科室class Fail err:%s\n", err)
	}

	for _, subkeshi := range subkeshi {
		//fmt.Printf("序号:%d 科室:%s\n",index,subkeshi)
		tmpelem, err := subkeshi.FindElement(selenium.ByCSSSelector, "div.sub-title")
		if err != nil {
			continue
		}
		keshiname, _ := tmpelem.Text()
		if strings.Contains(keshiname, "妇科") {
			subkeshi_fuke = subkeshi
		}
	}

	e, err := subkeshi_fuke.FindElement(selenium.ByCSSSelector, "div.sub-item-wrapper div.sub-item > span.v-link.clickable")
	if err != nil {
		log.Panicf("妇科 link Fail:%s\n", err)
	}
	e.Click()
	time.Sleep(5 * time.Second)
	fmt.Println(g.Wd.CurrentURL())

}
func main() {

	//go func() {
	//	<-g.ExitChan
	//	fmt.Println(1)
	//	<-g.ExitChan
	//	fmt.Println(2)
	//	// xxx 两个错误 停止操作
	//	//PrintResult()
	//	log.Fatalf("两次错误 停止!\n")
	//	//os.Exit(1)
	//}()

	go func() {
		WaitChan()
	}()

	defer func() {
		g.Wd.Quit()
		_ = g.Webservice.Stop()
	}()

	// 进入选号链接
	//getxuanhao()

	// 5. 直接选号链接
	beginurl = "https://www.114yygh.com/hospital/122/4eed349fd65ba377f936b177a56ebdcb/200001013/source"
	g.Wd.Get(beginurl)
	time.Sleep(20 * time.Second)

	//fmt.Println(e,err)
	//fmt.Println(e.Text())

	// 6.找到对应日期 点击
	for {

		//if dataSub,b:=LoopGetDate();b{   // 找到元素
		//	dataSub.Click()
		//	fmt.Println("click")
		//	break

		if e, err := g.Wd.FindElement(selenium.ByXPATH, fmt.Sprintf("//span[contains(text(),'%s')]", date)); err == nil {
			e.Click()
			fmt.Printf("日期: %s Click\n", date)
			time.Sleep(5 * time.Second)
			wg.Add(1)
			piaochan <- struct{}{}
		}

		wg.Wait()

		//找不到指定日期 和  刷新页面
		log.Println("刷新页面")
		g.Wd.ExecuteScript(`location.reload()`, nil)
		time.Sleep(5 * time.Second)
	}

}

func WaitChan() {
	for {
		select {
		case <-piaochan:
			// 7. 找到对应医生 查看能否预约

			sw, err := g.Wd.FindElement(selenium.ByXPATH, "//*[@id=\"main\"]/div[2]/div/div[1]/div/div/div[2]/div/div[4]")
			if err != nil {
				fmt.Printf("上午没找到err%s\n", err)
			}

			swlist, err := sw.FindElements(selenium.ByCSSSelector, ".list-item")
			//fmt.Println(len(swlist),err)

			for _, swsub := range swlist {
				doctorElem, _ := swsub.FindElement(selenium.ByCSSSelector, "div.item-wrapper div.title-wrapper > div.name")
				//fmt.Println(doctorElem.Text())
				if t, _ := doctorElem.Text(); strings.Contains(t, doctor) { // 找到对应医生

					haoElem, err := swsub.FindElement(selenium.ByCSSSelector, "div.right-wrapper div.button-wrapper")
					if err != nil {
						log.Panicln(err)
					}

					classhao, _ := haoElem.GetAttribute("class")
					//fmt.Println(classhao)
					if strings.Contains(classhao, "disabled") { //是否有号
						log.Printf("Doctor:%s, 没号了\n", doctor)
					} else {
						haoElem1, _ := haoElem.FindElement(selenium.ByTagName, "span")
						fmt.Println(haoElem1.Text())
						haoElem1.Click()
						log.Println("点击剩余")

						time.Sleep(2 * time.Second)
						peopleElem, err := g.Wd.FindElement(selenium.ByXPATH, "/html/body/div/div[1]/div[2]/div/div[1]/div/div/div[2]/div/div[2]/div[2]/div[1]/div")
						if err != nil {
							log.Printf("就诊人获取失败:%s\n", err)
						}
						//time.Sleep(100*time.Second)
						peopleElem.Click()
						time.Sleep(500 * time.Millisecond)
						yanzhengElem, err := g.Wd.FindElement(selenium.ByXPATH, "//*[@id=\"main\"]/div[2]/div/div[1]/div/div/div[2]/div/div[2]/div[7]/div[2]/form/div[2]/div/div/span/span/span")
						yanzhengElem.Click()
						log.Println("等待验证码10分钟")
						log.Println(g.Wd.CurrentURL())
						time.Sleep(10 * time.Minute)
						break
					}
				}
			}
			wg.Done()
		}
	}
}

func LoopGetDate() (dateSub selenium.WebElement, Success bool) {
	datelists, err := g.Wd.FindElements(selenium.ByCSSSelector, "div.calendar-list-wrapper > div.calendar-item")
	if err != nil {
		log.Panicf("预期日期错误 Fail:%s\n", err)
	}
	for _, dateSub := range datelists {
		dateelem, err := dateSub.FindElement(selenium.ByCSSSelector, "div.date-wrapper")
		if err != nil {
			continue
		}
		tmpElem, _ := dateelem.Text()
		if strings.Contains(tmpElem, date) {
			Success = true
			return dateSub, Success
		}
	}
	return nil, false
}
