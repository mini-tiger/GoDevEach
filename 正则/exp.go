package main

import (
	"fmt"
	cregex "github.com/mingrammer/commonregex"
	"time"
)
// https://github.com/mingrammer/commonregex

func main() {
	text := `192.168.1.1 John, please get that article on www.linkedin.com to me by 5:00PM on Jan 9th 2012. 4:00 
would be ideal, actually. If you have any questions, You can reach me at (519)-236-2723x341 or 
get in touch with my associate at harold.smith@gmail.com`

	dateList := cregex.Date(text)
	// ['Jan 9th 2012']
	timeList := cregex.Time(text)
	// ['5:00PM', '4:00']
	linkList := cregex.Links(text)
	// ['www.linkedin.com', 'harold.smith@gmail.com']
	phoneList := cregex.PhonesWithExts(text)
	// ['(519)-236-2723x341']
	emailList := cregex.Emails(text)
	// ['harold.smith@gmail.com']
	ipList :=cregex.IPs(text)
	fmt.Println(ipList)
	fmt.Println(dateList)
	fmt.Println(timeList)
	fmt.Println(linkList)
	fmt.Println(phoneList)
	fmt.Println(emailList)
	fmt.Println("=======================================")

	text1:=text
	text1=text1+"\n 07/08/2020 3:05:08"
	dateList=cregex.Date(text1)
	timeList=cregex.Time(text1) // xxx 时间解析不到秒


	fmt.Println(dateList)
	fmt.Println(timeList)
	datestr:=fmt.Sprintf("%s %s",dateList[len(dateList)-1],timeList[len(timeList)-1])
	fmt.Println(datestr)
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("01/02/2006 15:04", datestr,loc )
	fmt.Println(theTime)
}