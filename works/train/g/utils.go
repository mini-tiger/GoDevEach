package g

import (
	"strconv"
	"time"
)

func GetNow() (now int) {
	return int(time.Now().Unix())
}

func Getint(s string) (i64 int, b bool) {

	i64, err := strconv.Atoi(s)
	if err != nil {
		return i64, false
	} else {
		return i64, true
	}

}
func GetDateStr(unixStamp int) (dateTimeStr string) {
	timeLayout:="2006-01-02 15:04:05"
	dateTimeStr = time.Unix(int64(unixStamp), 0).Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	return
}