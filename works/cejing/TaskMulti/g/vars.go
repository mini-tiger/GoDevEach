package g

import (
	"github.com/robfig/cron"
	"net/http"
	"path/filepath"
	"runtime"
	"time"
)

var (
	Basedir string
	C       *cron.Cron
	//JsonBytes  []byte
	//Fp         string
	//HtmlBuffer *bytes.Buffer
	HttpClient *http.Client

	//TotalCommCellUrl        string
	ReqModelUrl       string
	ReqWellUrl string
	StartRunTime             int64
	ReqWellQXWJMUrl string
	ReqQXWJMListUrl string
	ReqSaveTaskUrl string
	//SendMailSync *sync.WaitGroup = new(sync.WaitGroup)
	//Loc *time.Location
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	Basedir = filepath.Dir(filepath.Dir(file))
	StartRunTime = time.Now().Unix()
	//Loc, _ = time.LoadLocation("Local")
}

const (
	TimeLayoutChi = "2006/01/02 15:04:05"

	TimeLayout = "2006-01-02 15:04:05" // status 需要使用
	//UA         = "Golang Downloader from tao.com"
)
