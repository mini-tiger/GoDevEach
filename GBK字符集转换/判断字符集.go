package main

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"path/filepath"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  判断字符集
 * @Version: 1.0.0
 * @Date: 2021/10/19 下午2:26
 */
func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}
func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//8bit中首个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

const (
	GBK     string = "GBK"
	UTF8    string = "UTF8"
	UNKNOWN string = "UNKNOWN"
)

//需要说明的是，isGBK()是通过双字节是否落在gbk的编码范围内实现的，
//而utf-8编码格式的每个字节都是落在gbk的编码范围内，
//所以只有先调用isUtf8()先判断不是utf-8编码，再调用isGBK()才有意义
func GetStrCoding(data []byte) string {
	if isUtf8(data) == true {
		return UTF8
	} else if isGBK(data) == true {
		return GBK
	} else {
		return UNKNOWN
	}
}

func main() {
	str := "月色真美，风也温柔，233333333，~！@#"   //go字符串编码为utf-8
	fmt.Println("before convert:", str) //打印转换前的字符串

	fmt.Println("coding:", GetStrCoding([]byte(str))) //判断是否是utf-8

	gbkData, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str)) //使用官方库将utf-8转换为gbk
	fmt.Println("gbk直接打印会出现乱码:", string(gbkData))                       //乱码字符串
	fmt.Println("coding:", GetStrCoding(gbkData))                       //判断是否是gbk

	utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(gbkData) //将gbk再转换为utf-8
	fmt.Println("coding:", GetStrCoding(utf8Data))                   //判断是否是utf-8
	fmt.Println("after convert:", string(utf8Data))                  //打印转换后的字符串

	fmt.Println("======================================")
	cwd, _ := os.Getwd()
	file, err := os.Open(filepath.Join(cwd, "gbk.csv"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytefile, _ := io.ReadAll(file)
	fmt.Printf("csv文件 codeing is:%#v\n", GetStrCoding(bytefile)) //

}
