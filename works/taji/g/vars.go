package g

import (
	"bytes"
	"github.com/robfig/cron"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	Runtype    = map[string]string{"key": "!=", "value": "已完成"}
	Solvetype  = map[string]string{"key": "in", "value": "(0,1)"}
	Basedir    string
	C          *cron.Cron
	JsonBytes  []byte
	Fp         string
	HtmlBuffer *bytes.Buffer
	HttpClient *http.Client

	DownLoadUrl        string
	RequestCommCellUrl string
	RequestTotalUrl    string
	StartRunTime       int64
	SendMailSync       *sync.WaitGroup = new(sync.WaitGroup)
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	Basedir = filepath.Dir(filepath.Dir(file))
	StartRunTime = time.Now().Unix()
}

const (
	TimeLayoutChi = "2006/01/02 15:04:05"

	TimeLayout = "2006-01-02 15:04:05" // status 需要使用
	UA         = "Golang Downloader from tao.com"
)
