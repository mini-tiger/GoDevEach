package main

import (
	"encoding/json"
	"fmt"
	excel "github.com/szyhf/go-excel"
)

type Standard struct {
	// use field name as default column name
	ID int
	// column means to map the column name

	// you can map a column into more than one field
	NamePtr *string `xlsx:"column(NameOf)"`
	// omit `column` if only want to map to column name, it's equal to `column(AgeOf)`
	Age int `xlsx:"AgeOf"`
	// split means to split the string into slice by the `|`
	Slice []int `xlsx:"split(|)"`
	// *Temp implement the `encoding.BinaryUnmarshaler`
	Temp *Temp  `xlsx:"column(UnmarshalString)"`
	Name string `xlsx:"column(NameOf)"`
	// support default encoding of json
	//TempEncoding *TempEncoding `xlsx:"column(UnmarshalString);encoding(json)"`
	// use '-' to ignore.
	Ignored string `xlsx:"-"`
}

// func (this Standard) GetXLSXSheetName() string {
// 	return "Some other sheet name if need"
// }

type Temp struct {
	Foo string
}

// Test decode tag
type TempEncoding Temp

// self define a unmarshal interface to unmarshal string.
func (this *Temp) UnmarshalBinary(d []byte) error {
	return json.Unmarshal(d, this)
}
func main() {
	println("defaultUsage")
	defaultUsage()
	println("simpleUsage")
	simpleUsage()

}

func simpleUsage() {
	// will assume the sheet name as "Standard" from the struct name.
	var stdList []Standard
	err := excel.UnmarshalXLSX("K:\\work\\project-dev\\GoDevEach\\读写excel\\stuTojson.xlsx", &stdList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", stdList)
}

func defaultUsage() {
	conn := excel.NewConnecter()
	err := conn.Open("K:\\work\\project-dev\\GoDevEach\\读写excel\\stuTojson.xlsx")
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
	rd, err := conn.NewReader("Standard")
	if err != nil {
		panic(err)
	}
	defer rd.Close()

	for rd.Next() {
		var s Standard
		// Read a row into a struct.
		err := rd.Read(&s)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v", s)
		fmt.Printf("   %+v", *s.NamePtr)
		fmt.Printf("   %+v\n", *s.Temp)
	}

	// Read all is also supported.
	// var stdList []Standard
	// err = rd.ReadAll(&stdList)
	// if err != nil {
	//   panic(err)
	//	 return
	// }
	// fmt.Printf("%+v",stdList)

	// map with string key is support, too.
	// if value is not string
	// will try to unmarshal to target type
	// but will skip if unmarshal failed.

	// var stdSliceList [][]string
	// err = rd.ReadAll(&stdSliceList)
	// if err != nil {
	//   panic(err)
	//	 return
	// }
	// fmt.Printf("%+v",stdSliceList)

	// var stdMapList []map[string]string
	// err = rd.ReadAll(&stdMapList)
	// if err != nil {
	//   panic(err)
	//	 return
	// }
	// fmt.Printf("%+v",stdMapList)

	// Using binary instead of file.
	// xlsxData, err := ioutil.ReadFile(filePath)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// conn := excel.NewConnecter()
	// err = conn.OpenBinary(xlsxData)
}
