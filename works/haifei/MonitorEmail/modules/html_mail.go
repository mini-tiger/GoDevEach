package modules

import (
	"errors"
	"fmt"
	"haifei/MonitorEmail/g"
	"html/template"
)

func getFailNum(m map[string]interface{}, s string) interface{} {
	//fmt.Println(11111111,m,s)
	return m[s]
}

func Color(num float64) string {
	if num > 0 {
		return "FF9999"
	} else {
		return "FFFFFF"
	}
}

func  GenHtml() {
	// 	t, err := template.ParseFiles("./ex/tpl.html")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//var ta []map[string]interface{} = make([]map[string]interface{}, 0)
	//for i := 0; i < 5; i++ {
	//	dd:=map[string]interface{}{"BIANMA": strconv.Itoa(i), "COMMCELL": strconv.Itoa(i), "CLIENT": strconv.Itoa(i)}
	//	ta = append(ta, dd)
	//}
	//
	////var h htmlData
	//this.HtmlData.TotalData = ta
	//this.HtmlData.CommCell = []string{"10", "中文"}
	HtmlData.CommCellVisible = true
	//funcMap := template.FuncMap{"totitle": upper}
	//fmt.Printf("%+v\n", h)

	funcMap := template.FuncMap{"FailtoNum": getFailNum, "Color": Color}

	t := template.New("mail").Funcs(funcMap)
	// t = template.Must(t.ParseFiles("templates/layout.html", "templates/index.html"))
	t = template.Must(t.ParseFiles(g.GetConfig().HtmlTpl))
	// if err != nil {
	// 	panic(err)
	// }

	err := t.ExecuteTemplate(g.HtmlBuffer, "mail", HtmlData) //指定模板文件{{ define "layout" }}……{end}
	if err != nil {
		//_ = g.GetLog().Error("生成HTML文件失败:%s\n", err)
		panic(errors.New(fmt.Sprintf("生成HTML文件失败:%s\n", err)))
	}

}
