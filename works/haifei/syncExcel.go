package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"path"
)

//https://github.com/360EntSecGroup-Skylar/excelize
// 中文 https://xuri.me/excelize/zh-hans/
// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#File.GetSheetName

const (
	ReadFile = "测试台账.xlsx"
)

var CurrDir string

func main() {
	//_, file, _, _ := runtime.Caller(0)
	CurrDir = "F:\\BaiduNetdiskDownload\\海尔"
	//Create()
	Read()
	//AddChart()
	// 更改样式
	//style()
}

//func getCellValue(f *excelize.File, colnum int) (cell string) {
//	cell, err := f.GetCellValue("Sheet1", )
//	if err != nil {
//		return
//	}
//	return
//}

func Read() {
	var ColMap map[int]string = make(map[int]string, 0)
	var excelData []map[string]string = make([]map[string]string, 0)

	f, err := excelize.OpenFile(path.Join(CurrDir, ReadFile))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	//cell, err := f.GetCellValue("Sheet1", "B2")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("sheet1 B2:",cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")

	// 生成字段 顺序对照,sql 的字段

	for index, value := range rows[0] {
		ColMap[index] = value

	}
	//fmt.Println(ColMap)

	baseSql := ""                     // 多行SQL使用
	baseSqlSlice := make([]string, 0) // 单选Sql使用
	// 生成每行 是 map， 总共是 []map
	for _, row := range rows[2:] {
		tmp := make(map[string]string, 0)
		tmpsql := "into ABC("
		tmpsqlvalue := "values("
		for cindex, colCell := range row {
			if v, ok := ColMap[cindex]; ok {
				tmp[v] = colCell
				tmpsql = tmpsql + v + ","
				tmpsqlvalue = tmpsqlvalue + colCell + ","
			}
			excelData = append(excelData, tmp)
			//fmt.Print(colCell, "\t")
		}
		//fmt.Println(tmp)
		tmpsql = tmpsql[:len(tmpsql)-1]
		tmpsql = tmpsql + ")"

		tmpsqlvalue = tmpsqlvalue[:len(tmpsqlvalue)-1]
		tmpsqlvalue = tmpsqlvalue + ")"

		baseSqlSlice = append(baseSqlSlice, tmpsql+" "+tmpsqlvalue)
		//fmt.Println(tmpsql)
		//fmt.Println(tmpsqlvalue)
		baseSql = baseSql + " " + tmpsql + " " + tmpsqlvalue
	}
	//genMuchSql(baseSql) //生成一次插入多条
	//genSql(baseSqlSlice)
	//fmt.Println(excelData)
}

func genMuchSql(baseSql string) {
	//baseSql:="Insert All"+
	//"into ABC('bianma','hostName',"currentDatasize","starttime","lxdh") values('20','1','2','3','4')'+
	//into ABC("bianma","hostName","starttime","lxdh","currentDatasize") values('20','1','2','3','4')
	//Select 1 from dual'
	Sql := "Insert All " + baseSql + " " + "Select 1 from dual"
	fmt.Println(Sql)
}
func genSql(baseSqlslice []string) {
	for index, value := range baseSqlslice {
		fmt.Println(index+1, fmt.Sprintf("Insert %s", value))
	}

}
