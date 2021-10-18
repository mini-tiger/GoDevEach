package main

import (
	"os/exec"

	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(sb interface{}, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		switch sb.(type) {
		case string:
			sb1 := sb.(string)
			str, _ = simplifiedchinese.GB18030.NewDecoder().String(sb1)
		case []byte:
			sb1 := sb.([]byte)
			var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(sb1)
			str = string(decodeBytes)
		}

		// todo simplifiedchinese.GB18030.NewDecoder().Reader()
	case UTF8:
		fallthrough
	default:
		str = ""
	}

	return str
}

func main() {
	cmd := exec.Command("ipconfig")
	//ss:=make([]byte,0)
	//buf:=bytes.NewBuffer(ss)
	//buf:=new(bytes.Buffer)
	//cmd.Stdout=buf
	//cmd.Run()
	//
	//r,_:=buf.ReadBytes('\n')
	//for _,v:=range r{
	//	fmt.Println(string(v))
	//}
	buf, _ := cmd.CombinedOutput()
	fmt.Println(ConvertByte2String(buf, GB18030)) // todo 统一返回

	// todo 逐行返回
	//aa:=bytes.NewBuffer(buf)
	//for   {
	//	line,err:=aa.ReadString('\n')
	//	//fmt.Println(line)
	//	if ok,_:=regexp.MatchString("^\r\n$",line);ok{  // todo 跳过只有换行符的 空行
	//		continue
	//	}
	//	line=strings.Split(line,"\r")[0] // 命令返回每行 结尾是\r\n,去掉其中一个
	//	//line=ss[0]
	//	if err!=nil || err==io.EOF{
	//		break
	//	}
	//	fmt.Println(ConvertByte2String(line,GB18030))
	//	}

}
