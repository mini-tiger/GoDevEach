package funcs

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"远程执行命令/modules"
)

/**
 * @Author: Tao Jun
 * @Description: funcs
 * @File:  utils
 * @Version: 1.0.0
 * @Date: 2021/2/22 上午10:03
 */

func DiskFormat(s string) map[string]string {
	m := make(map[string]string, 0)
	for _, row := range strings.Split(s, "\n") {
		row = strings.TrimSpace(row)
		d := strings.Split(row, " ")
		if len(d) == 2 {
			m[d[0]] = d[1]
		}
	}
	return m
}
func gte(usage string, rate float64) bool {

	var i float64
	s := strings.Split(usage, "%")
	i, err := strconv.ParseFloat(s[0], 10)
	if err != nil {
		return false
	}

	if i >= rate {
		return true
	}
	return false
}
func MailHtml(ResultHosts []*modules.HostMonitor) {
	funcMap := template.FuncMap{"gte": gte}
	t := template.New("mail").Funcs(funcMap)
	// t = template.Must(t.ParseFiles("templates/layout.html", "templates/index.html"))
	t = template.Must(t.ParseFiles("/home/go/GoDevEach/运维/远程执行命令/mail.htm"))
	// if err != nil {
	// 	panic(err)
	// }
	f, _ := os.OpenFile("/home/go/GoDevEach/运维/远程执行命令/mailResult.htm", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer f.Close()

	err := t.ExecuteTemplate(f, "mail", ResultHosts) //指定模板文件{{ define "layout" }}……{end}
	if err != nil {
		fmt.Println(err)
	}
}
