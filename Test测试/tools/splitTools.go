package tools

import (
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: tools
 * @File:  splitTools
 * @Version: 1.0.0
 * @Date: 2021/4/25 下午1:25
 */

func Split(s string, sep string) (result []string) {
	i := strings.Index(s, sep)

	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):] // 这里使用len(sep)获取sep的长度
		i = strings.Index(s, sep)
	}
	result = append(result, s)
	return
}
