package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/szyhf/go-excel"
)

// defined a struct
type Standard1 struct {
	// xxx column 要与title 名字一样
	// use field name as default column name
	Phone int `xlsx:"column(phone)"`
	// column means to map the column name
	Name string `xlsx:"column(name)"`
	// you can map a column into more than one field
	NamePtr *string `xlsx:"column(name)"` // xxx 重复读name列
	// omit `column` if only want to map to column name, it's equal to `column(AgeOf)`
	Age  int    `xlsx:"age" validate:"gte=0,lte=130"` // xxx 与xlsx:"column(age)"效果一样
	Addr string `xlsx:"column(addr)"`
	Mail string `xlsx:"column(mail);default(abc@mail.com)" validate:"required,email"`
	IP   string `xlsx:"column(IP)" validate:"required,ipv4"`

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

	conn := excel.NewConnecter()
	err := conn.Open("K:\\work\\project-dev\\GoDevEach\\读写excel\\stuTojson2.xlsx")
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
		Sheet:         "student_list",
		TitleRowIndex: 1, // xxx title 行的位置
		Skip:          0,
		Prefix:        "",
		Suffix:        "",
	}
	rd, err := conn.NewReaderByConfig(cfg)
	validate := validator.New()
	for rd.Next() {
		//fmt.Printf("rd:%v\n",rd.GetTitles())
		var s Standard1
		// Read a row into a struct.
		err := rd.Read(&s)
		if err != nil {
			panic(err)
		}
		//xxx 检验数据
		err = validate.Struct(&s)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Printf("%#v\n", err)
				fmt.Println("Namespace:", err.Namespace())
				fmt.Println("Field:", err.Field())
				fmt.Println("StructNamespace:", err.StructNamespace())
				fmt.Println("StructField:", err.StructField())
				fmt.Println("Tag:", err.Tag(), err)
				fmt.Println("ActualTag:", err.ActualTag())
				fmt.Println("Kind:", err.Kind())
				fmt.Println("Type:", err.Type())
				fmt.Println("Value:", err.Value())
				fmt.Println("Param:", err.Param())
				fmt.Println()
			}
		}
		fmt.Printf("%+v\n", s)

	}
}
