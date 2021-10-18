package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  createRead官方依赖包
 * @Version: 1.0.0
 * @Date: 2021/10/18 上午9:44
 */

func main() {
	cwd, _ := os.Getwd()
	file, err := os.Open(filepath.Join(cwd, "gbk.csv"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytefile, _ := io.ReadAll(file)
	fmt.Println(string(bytefile)) //乱码

	var b *bytes.Buffer = new(bytes.Buffer)
	b.Write(bytefile)

	reader := simplifiedchinese.GB18030.NewDecoder().Reader(b)
	body, _ := ioutil.ReadAll(reader)
	fmt.Println(string(body))
}
