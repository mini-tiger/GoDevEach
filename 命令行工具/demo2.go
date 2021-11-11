package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  demo2
 * @Version: 1.0.0
 * @Date: 2021/11/11 下午6:27
 */

func main() {

}

// https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/

//--flag    // 布尔flags, 或者非必须选项默认值
//--flag x  // 只对于没有默认值的flags
//--flag=x

//xxx  ./demo2 -f 1
func init() {
	var ip = flag.IntP("flagname", "f", 1234, "help message")

	// 设置非必须选项的默认值
	//flag.Lookup("flagname").NoOptDefVal = "4321"
	flag.Parse()
	fmt.Println(*ip)
}
