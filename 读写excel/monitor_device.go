package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/szyhf/go-excel"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// defined a struct
type Row struct {
	// xxx column 要与title 名字一样
	// use field name as default column name
	//Phone int `xlsx:"column(phone)"`
	// column means to map the column name
	Name string `xlsx:"column(Name)"`
	// you can map a column into more than one field
	Os string `xlsx:"column(OS)" validate:"required,OsValidation"`
	// omit `column` if only want to map to column name, it's equal to `column(AgeOf)`
	Port string `xlsx:"column(port)" validate:"required,PortValidation"`
	//Mail string `xlsx:"column(mail);default(abc@mail.com)" validate:"required,email"`
	IP string `xlsx:"column(IP)" validate:"required,ipv4"`

	// split means to split the string into slice by the `|`

	//Slice   []int `xlsx:"split(|);req(Slice);"` // xxx req： 没有title 返回错误

	// *Temp implement the `encoding.BinaryUnmarshaler`
	//Temp    *Temp `xlsx:"column(UnmarshalString)"`
	// support default encoding of json
	//TempEncoding *TempEncoding `xlsx:"column(UnmarshalString);encoding(json)"`
	// use '-' to ignore.
	//Ignored string `xlsx:"-"`
}

func main() {

	_, file, _, _ := runtime.Caller(0)
	CurrDir := filepath.Dir(file)

	conn := excel.NewConnecter()
	err := conn.Open(path.Join(CurrDir, "monitor_device.xlsx"))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Generate an new reader of a sheet
	// sheetNamer: if sheetNamer is string, will use sheet as sheet name.
	//             if sheetNamer is int, will i'th sheet in the workbook, be careful the hidden sheet is counted. i ∈ [1,+inf]
	//             if sheetNamer is a object implements `GetXLSXSheetName()string`, the return value will be used.
	//             otherwise, will use sheetNamer as struct and reflect for it's name.
	// 			   if sheetNamer is a slice, the type of element will be used to infer like before.

	cfg := &excel.Config{
		Sheet:         "monitor_devices",
		TitleRowIndex: 1, // xxx title 行的位置
		Skip:          0,
		Prefix:        "",
		Suffix:        "",
	}
	rd, err := conn.NewReaderByConfig(cfg)
	validate := validator.New()
	validate.RegisterValidation("OsValidation", OsValidationFunc)
	validate.RegisterValidation("PortValidation", PortValidationFunc)

	var index = 3
	for rd.Next() {
		//fmt.Printf("rd:%v\n",rd.GetTitles())
		var r Row
		// Read a row into a struct.
		err := rd.Read(&r)
		if err != nil {
			continue
		}
		fmt.Println("============", index)

		//xxx 检验数据
		err = validate.Struct(&r)
		if err != nil {
			errslice := err.(validator.ValidationErrors)
			if len(errslice) > 0 {
				fmt.Printf("行: %d 列: %s 值: %s 不符合:%s 规范\n", index, (errslice[0]).Field(), errslice[0].Value(), errslice[0].Field())
			}

			//for _, err := range err.(validator.ValidationErrors) {
			//
			//	fmt.Printf("%#v\n", err)
			//	fmt.Println("Namespace:", err.Namespace())
			//	fmt.Println("Field:", err.Field())
			//	fmt.Println("StructNamespace:", err.StructNamespace())
			//	fmt.Println("StructField:", err.StructField())
			//	fmt.Println("Tag:", err.Tag(), err)
			//	fmt.Println("ActualTag:", err.ActualTag())
			//	fmt.Println("Kind:", err.Kind())
			//	fmt.Println("Type:", err.Type())
			//	fmt.Println("Value:", err.Value())
			//	fmt.Println("Param:", err.Param())
			//	fmt.Println()
			//}
		}
		fmt.Printf("%+v\n", r)
		index++
	}
}

func OsValidationFunc(f1 validator.FieldLevel) bool {
	// f1 包含了字段相关信息
	// f1.Field() 获取当前字段信息
	// f1.Param() 获取tag对应的参数
	// f1.FieldName() 获取字段名称
	//fmt.Printf("%+v\n",f1)
	value := strings.ToLower(f1.Field().String())
	return value == "linux" || value == "windows"
}

func PortValidationFunc(f1 validator.FieldLevel) bool {
	// f1 包含了字段相关信息
	// f1.Field() 获取当前字段信息
	// f1.Param() 获取tag对应的参数
	// f1.FieldName() 获取字段名称
	//fmt.Printf("%+v\n",f1)
	value, err := strconv.Atoi(f1.Field().String())
	if err != nil {
		return false
	}
	return value > 1 && value < 65535
}
