package g

import (
	"path/filepath"
	"runtime"
)

var (
	Basedir string
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	Basedir = filepath.Dir(filepath.Dir(file))

}

const (
	TimeLayoutChi = "2006/01/02 15:04:05"

	TimeLayout = "2006-01-02 15:04:05" // status 需要使用
	UA         = "Golang Downloader from tao.com"
)
