package main

import (
	"html/template" //模板类似 jinja  if range for 等查看徇私文档
	"os"
	"strings"
	//template.ParseFiles 文件模板方法
)

func upper(s string) string {
	return strings.ToTitle(s)
}

func main() {
	// 	t, err := template.ParseFiles("./ex/tpl.html")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	data := struct {
		Title string
		Days  []string
		Arg   bool
	}{
		Title: "golang html template demo",
		Arg:   true,
		Days:  []string{"Mon", "Tue", "Wed", "Ths", "Fri", "Sat", "Sun"},
	}
	// 	// tmpl, err := template.New("tt").Funcs(template.FuncMap{"totitle": upper}).Parse(t)
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	// tmpl.Execute(os.Stdout, data)
	// 	funcMap := template.FuncMap{"totitle": upper}
	// 	tp := template.New("tt").Funcs(funcMap)
	// 	// t = template.Must(t.ParseFiles("templates/layout.html", "templates/index.html"))
	// 	tp = template.Must(t, err)
	// 	tp.Execute(os.Stdout, data)

	funcMap := template.FuncMap{"totitle": upper}
	t := template.New("layout").Funcs(funcMap)
	// t = template.Must(t.ParseFiles("templates/layout.html", "templates/index.html"))
	t = template.Must(t.ParseFiles("C:\\work\\go-dev\\GoDevEach\\html模板转换\\ex\\tpl.html"))
	// if err != nil {
	// 	panic(err)
	// }
	t.ExecuteTemplate(os.Stdout, "layout", data) //指定模板文件{{ define "layout" }}……{end}

}
