package main

import (
	"fmt"

	excelize "github.com/xuri/excelize/v2"
	"log"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

//https://github.com/360EntSecGroup-Skylar/excelize
// 中文 https://xuri.me/excelize/zh-hans/
// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#File.GetSheetName

const (
	ReadFile  = "开发计划.xlsx"
	WriteFile = "book1.xlsx"
	ChartFile = "book2.xlsx"
)

var CurrDir string

type Student struct {
	Name   string
	age    int
	Phone  string
	Gender string
	Mail   string
}

func GetStudents() []Student {
	students := make([]Student, 0)
	for i := 0; i < 10; i++ {
		stu := Student{}
		stu.Name = "name" + strconv.Itoa(i+1)
		stu.Mail = stu.Name + "@chairis.cn"
		stu.Phone = "1380013800" + strconv.Itoa(i)
		stu.age = 20
		stu.Gender = "男"
		students = append(students, stu)
	}
	return students
}

func main() {
	_, file, _, _ := runtime.Caller(0)
	CurrDir = path.Dir(file)
	Create()
	Read()
	AddChart()
	// 更改样式
	style()
}

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func formatAtom(v reflect.Value) string {
	fmt.Println(v.Kind())
	switch v.Kind() {

	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
		// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)

	case reflect.Slice:
		strArray := make([]string, v.Len())
		//fmt.Printf("%v\n",v.Slice(0,v.Len()))
		type a interface{}
		var b a
		if v.CanInterface() { // 是否可以转换为接口
			b = v.Interface()
		} else {
			return "111111" // 否则返回错误
		}
		for i, vv := range b.([]interface{}) { //todo 再将 b 转换为[]interface
			strArray[i] = vv.(string) // 元素转换为字符串
		}

		return strings.Join(strArray, ",")
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 5, 32)

	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

func Create() {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet2")
	// Set value of a cell.

	stus := GetStudents()
	// todo  坐标转单元格
	// 坐标转单元格 等工具函数  见 中文文档  工具函数
	cells := GetFieldName(stus[0]) // 反射结构体 所有的元素 为 数组
	//fmt.Println(cells)

	for r := 0; r < len(stus); r++ {
		for c := 0; c < len(cells); c++ {
			cellname, _ := excelize.CoordinatesToCellName(c+1, r+1) // 坐标 转 单元格
			cell := reflect.ValueOf(&stus[r]).Elem()                // 反射 结构体数据
			//fmt.Printf("%T\n",cell.FieldByName(cells[c]))
			_ = f.SetCellValue("Sheet2", cellname, formatAtom(cell.FieldByName(cells[c]))) // 根据元素名 提取 结构体的数据，并断言出类型数据
		}
		//cellname,_:=excelize.CoordinatesToCellName(1,r+1)
		//f.SetCellValue("Sheet2", cellname, stus[r].Name)
		//
		//cellname,_=excelize.CoordinatesToCellName(2,r+1)
		//f.SetCellValue("Sheet2", cellname, stus[r].age)
		//
		//cellname,_=excelize.CoordinatesToCellName(3,r+1)
		//f.SetCellValue("Sheet2", cellname, stus[r].Gender)
		//
		//cellname,_=excelize.CoordinatesToCellName(4,r+1)
		//f.SetCellValue("Sheet2", cellname, stus[r].Mail)

	}
	//f.SetCellValue("Sheet2", "A2", "Hello world.")
	_ = f.SetCellValue("Sheet1", "A1", 100)
	f.SetSheetName("Sheet1", "更改style")

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save xlsx file by the given path.
	err := f.SaveAs(path.Join(CurrDir, WriteFile))
	if err != nil {
		fmt.Println(err)
	}
}

func Read() {
	f, err := excelize.OpenFile(path.Join(CurrDir, ReadFile))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("sheet1 B2:", cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func AddChart() {
	categories := map[string]string{"A2": "Small", "A3": "Normal", "A4": "Large", "B1": "Apple", "C1": "Orange", "D1": "Pear"}
	values := map[string]int{"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}
	f := excelize.NewFile()
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}
	err := f.AddChart("Sheet1", "E1", `{"type":"col3DClustered","series":[{"name":"Sheet1!$A$2","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$2:$D$2"},{"name":"Sheet1!$A$3","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$3:$D$3"},{"name":"Sheet1!$A$4","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$4:$D$4"}],"title":{"name":"Fruit 3D Clustered Column Chart"}}`)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Save xlsx file by the given path.
	err = f.SaveAs(path.Join(CurrDir, ChartFile))
	if err != nil {
		fmt.Println(err)
	}
}

func style() {
	xlsx, err := excelize.OpenFile(path.Join(CurrDir, WriteFile))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取工作表，打印出数据
	rows, err := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		for _, cell := range row {
			fmt.Println(cell)
		}
	}

	sheetName := xlsx.GetSheetName(1)

	cell, err := xlsx.GetCellValue(sheetName, "A1")
	fmt.Println(cell)

	// 合并 A1-E1
	xlsx.MergeCell(sheetName, "A1", "E1")

	// A1-E1 的内容居中
	style, err := xlsx.NewStyle(`{"alignment":{"horizontal":"center"}}`)
	if err != nil {
		fmt.Println("创建样式失败", err)
		return
	}
	xlsx.SetCellStyle(sheetName, "A1", "E1", style)

	xlsx.Save()
}
