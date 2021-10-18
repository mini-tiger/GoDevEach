package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func main() {

	s := "GBK 与 UTF-8 编码转换测试"
	gbk, err := Utf8ToGbk([]byte(s))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(gbk))
	}

	utf8, err := GbkToUtf8(gbk)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(utf8))
	}

	// csv
	cwd, _ := os.Getwd()
	file, err := os.Open(filepath.Join(cwd, "gbk.csv"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// csv内容转换为UTF8
	bytefile, _ := io.ReadAll(file)
	utf8str, err := GbkToUtf8(bytefile)
	fmt.Println(string(utf8str))

	// 另存为utf8
	csvFile, err := os.Create("GBKToUTF8.csv")

	// xxx csv文件的开头写入 UTF-8 BOM
	csvFile.WriteString("\xEF\xBB\xBF")

	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvFile.Write(utf8str)
}
