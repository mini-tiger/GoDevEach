package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
)

type TotalRow struct {
	BIANMA   string
	COMMCELL string
	CLIENT   string
}

type htmlData struct {
	CommCell        []string
	TotalData       []TotalRow
	CommCellVisible bool
	M               map[string]string
}

func main() {
	// 	t, err := template.ParseFiles("./ex/tpl.html")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	var ta []TotalRow = make([]TotalRow, 0)
	for i := 0; i < 5; i++ {
		ta = append(ta, TotalRow{BIANMA: strconv.Itoa(i), COMMCELL: strconv.Itoa(i), CLIENT: strconv.Itoa(i)})
	}

	var h htmlData
	h.TotalData = ta
	h.CommCell = []string{"10", "中文"}
	h.CommCellVisible = true
	h.M = map[string]string{"COMMCELL": "1222222", "ABC": "2333333333"}

	//funcMap := template.FuncMap{"totitle": upper}
	fmt.Printf("%+v\n", h)

	t := template.New("mail")
	// t = template.Must(t.ParseFiles("templates/layout.html", "templates/index.html"))
	t = template.Must(t.ParseFiles("/home/go/GoDevEach/html模板转换/ex/1.HTML"))
	// if err != nil {
	// 	panic(err)
	// }

	err := t.ExecuteTemplate(os.Stdout, "mail", h) //指定模板文件{{ define "layout" }}……{end}
	if err != nil {
		fmt.Println(err)
	}

}
