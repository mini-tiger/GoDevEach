package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"runtime"
	"path"
	"strconv"
)
//https://godoc.org/github.com/tealeg/xlsx


const (
	ReadFile  = "开发计划.xlsx"
	WriteFile = "stu.xlsx"
)

var CurrDir string

func main() {
	_, file, _, _ := runtime.Caller(0)
	CurrDir = path.Dir(file)
	Read()
	Write()
}

func Read() {

	excelFileName := path.Join(CurrDir, ReadFile)
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Printf("current sheetName:%s\n", sheet.Name)
		for rowi, row := range sheet.Rows {
			fmt.Printf("第%d行:", rowi)
			for celli, cell := range row.Cells {
				text := cell.String()
				fmt.Printf("第%d列,内容:%s", celli, text)
			}
			fmt.Println()
		}
	}
}

type Student struct {
	Name   string
	age    int
	Phone  string
	Gender string
	Mail   string
}

func Write() {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("student_list")
	if err != nil {
		fmt.Printf(err.Error())
	}
	stus := GetStudents()
	// add header

	head := sheet.AddRow()
	for _, headName := range []string{"名字", "年龄", "电话", "性别", "邮箱"} {
		headCell := head.AddCell()
		headCell.Value = headName
	}
	//add data
	for _, stu := range stus {
		row := sheet.AddRow()
		nameCell := row.AddCell()
		nameCell.Value = stu.Name

		ageCell := row.AddCell()
		ageCell.Value = strconv.Itoa(stu.age)

		phoneCell := row.AddCell()
		phoneCell.Value = stu.Phone

		genderCell := row.AddCell()
		genderCell.Value = stu.Gender

		mailCell := row.AddCell()
		mailCell.Value = stu.Mail
	}
	err = file.Save(path.Join(CurrDir, WriteFile))
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Println("\n\nexport success")
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
