package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: utils
 * @File:  tools
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午5:45
 */

func Md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)

	return strings.ToUpper(fmt.Sprintf("%x", w.Sum(nil)))
}
