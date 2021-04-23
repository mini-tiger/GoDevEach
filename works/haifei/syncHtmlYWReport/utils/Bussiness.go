package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FormatData(s *string)  {
	var fl float64
	var err error
	if *s == "0" || strings.Contains(*s, "不适用") || *s == "" {
		*s="0"
	}

	if strings.Contains(*s, "(") {
		*s = strings.Split(*s, "(")[0]
	}

	switch true {
	case strings.Contains(*s, "TB"):
		*s = strings.Split(*s, " TB")[0]
		fl, err = strconv.ParseFloat(*s, 64)
		if err != nil {
			break
		}

		//am:=strconv.FormatFloat(fl*1024,'E',-1,64)
		*s = fmt.Sprintf("%.2f", fl*1024*1024)
	case strings.Contains(*s, "GB"):
		*s = strings.Split(*s, " GB")[0]
		fl, err = strconv.ParseFloat(*s, 64)
		if err != nil {
			//_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
			break
		}

		//am:=strconv.FormatFloat(fl*1024,'E',-1,64)
		*s = fmt.Sprintf("%.2f", fl*1024)

	case strings.Contains(*s, "MB"):
		*s = strings.Split(*s, " MB")[0]
		fl, err = strconv.ParseFloat(*s, 64)
		if err != nil {
			//_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
			break
		}
		*s = fmt.Sprintf("%.2f", fl)
	case strings.Contains(*s, "KB"):
		*s = strings.Split(*s, " KB")[0]
		fl, err = strconv.ParseFloat(*s, 64)
		if err != nil {
			//_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
			break
		}
		*s = fmt.Sprintf("%.2f", fl/1024)
	case strings.Contains(*s, "Bytes"):
		*s = strings.Split(*s, " Bytes")[0]
		fl, err = strconv.ParseFloat(*s, 64)
		if err != nil {
			//_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
			break
		}
		*s = fmt.Sprintf("%.2f", fl/1024/1024)
	case strings.Contains(*s, "字节"):
		*s = strings.Split(*s, " 字节")[0]
		fl, err = strconv.ParseFloat(*s, 64)
		if err != nil {
			//_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
			break
		}
		*s = fmt.Sprintf("%.2f", fl/1024/1024)
	case strings.Contains(*s, "Not Run because of another job running for the same subclient"):
		break
	default:

		floatNum, _ := regexp.MatchString(`(\d+)\.(\d+)`, *s)
		//fmt.Println(i,err)

		intNum, _ := regexp.MatchString(`(\d+)`, *s)
		if floatNum || intNum {
			break
		}

		application, e := strconv.Atoi(strings.TrimSpace(*s))
		if e != nil {
			//_ = Log.Error("change type string:%s,err:%s\n", *s, e)
			return
		}
		if application == 0 {
			*s = "0"
		}

	}

}
