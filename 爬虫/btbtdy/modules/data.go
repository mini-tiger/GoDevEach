package modules

import (
	"btbtdy/g"
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"regexp"
	"strconv"
)

var Title, DownUrl string

func Parse(url string) {
	elem, err := g.Wd.FindElement(selenium.ByID, "video-down")
	if err != nil {
		log.Printf("video-down  失败:%v,跳过\n", err)
		g.ExitChan <- struct{}{}

	}
	listelem, err := elem.FindElements(selenium.ByTagName, "p")

	if err != nil {
		log.Printf("video-down -> p 失败:%v,跳过\n", err)
		g.ExitChan <- struct{}{}

	}

	if len(listelem) == 2 {
		Title, _ = listelem[0].Text()
		DownUrl, _ = listelem[1].Text()

		reg1 := regexp.MustCompile(`(.*)(\d{2,3}).html`)
		if reg1 == nil { //解释失败，返回nil
			fmt.Println("regexp err")
			return
		}
		//根据规则提取关键信息
		result1 := reg1.FindAllStringSubmatch(url, -1)
		//fmt.Println("result1 = ", result1[0][2])
		//value:=result1[0][2]
		value, _ := strconv.Atoi(result1[0][2])
		g.Result[value] = DownUrl
	}

	fmt.Printf("title:%s,DownUrl:%s\n", Title, DownUrl)
}
