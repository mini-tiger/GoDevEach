package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)



func Clear(v interface{}) { // 必须传入指针
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}
//var db *sql.DB
//var rows *sql.Rows
//func CheckOracleConn(dsn *string) (err error) {
//	_ = os.Setenv("NLS_LANG", "")
//	//if len(os.Args) != 2 {
//	//	log.Fatalln(os.Args[0] + " user/password@host:port/sid")
//	//}
//
//	db, err = sql.Open("oci8", *dsn)
//	//fmt.Printf("%+v\n",db)
//	if err != nil {
//		return err
//	}
//
//	defer func() {
//		_ = db.Close()
//	}()
//
//	rows, err = db.Query("select 3.14 from dual")
//	if err != nil {
//		return err
//	}
//	//fmt.Println(rows.Next())
//	defer func() {
//		_ = rows.Close()
//	}()
//	return nil
//}
func main() {
	a := "D:\\work\\project-dev\\GoDevEach\\works\\haifei\\syncHtml\\HisData\\10.135.13.164_5day_Chinese.html"
	fmt.Println(filepath.Base(a))
	fmt.Println(filepath.Split(a))
	dir, filename := filepath.Split(a)
	fmt.Println(dir)
	fmt.Println(filepath.Dir(dir))
	fmt.Println(filepath.Dir(filename))
	fmt.Println(filepath.Base(dir))
	fmt.Println(strings.Split(a, string(os.PathSeparator)))
	fmt.Println(strings.Split(a, "D:\\work\\project-dev\\GoDevEach\\works\\haifei\\syncHtml\\HisData\\"))
}

//var filesPool = sync.Pool{
//	New: func() interface{} {
//		b := make([]string, 0)
//		return &b
//	},
//}

func Split(s *string, maxLen int) (*string) {
	if len(*s) > maxLen && (*s) != "NULL" {
		(*s) = (*s)[0 : maxLen-153]
		return s
	} else {
		return s
	}

}

func Rebuildstr(s *string) {
	//rstr := ""
	if strings.Contains(*s, "N/A") {
		*s = strings.Replace(*s, "N/A", "", -1)

	}


	if strings.Contains(*s, "%") {
		*s = strings.Replace(*s, "%", "%%", -1)
		//s = &rstr
	}

	*s = strings.ReplaceAll(*s, "'", "")
	//if len(rstr) > 4000 {
	//	rstr = rstr[0:3947]
	//}

	//return &rstr
}

func FormatTime(toBeCharge *string, ChineseFile bool) {
	//toBeCharge := "09/04/2019 11:03:16"
	timeLayout := "01/02/2006 15:04:05" // 	月/日/年
	if ChineseFile {
		timeLayout = "2006/01/02 15:04:05"
	}
	//转化所需模板
	loc, _ := time.LoadLocation("Local")                             //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, *toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	//sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	// xxx 如果时间不对，可能是时间排列顺序的问题
	if theTime.Unix() < 0 {
		timeLayout = "01/02/2006 15:04:05"                              // 	月/日/年 xxx 此与线上版本不一样
		theTime, _ = time.ParseInLocation(timeLayout, *toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	}
	*toBeCharge = strconv.FormatInt(theTime.Unix(), 10)
	//return strconv.FormatInt(sr, 10)
}

func NewDateFileName(oldfileAbs string, bakdir string) string {

	filename := filepath.Base(oldfileAbs)
	ext := filepath.Ext(oldfileAbs)

	filenameOnly := strings.TrimSuffix(filename, ext)
	//fmt.Println(oldfileAbs)
	//fmt.Println(filepath.Join(bakdir,fmt.Sprintf("%s_%s%s",filenameOnly,time.Now().Format("2006-01-02"),ext)))

	return filepath.Join(bakdir, fmt.Sprintf("%s_%d%s", filenameOnly, time.Now().UnixNano(), ext))
	// 同一天的目录中有 相同文件名， 文件名中加入 日期
	// /home/go/GoDevEach/works/haifei/syncHtml/htmlData/1monthReport/cv10/cv10-bj1/14day_4copy.html
	// /home/go/GoDevEach/works/haifei/syncHtml/htmlData/HisData/bak/2019-11-08/14day_4copy_2019-11-08.html

}

func ClearMap(m map[string]interface{}) {
	for k, _ := range m {
		delete(m, k)
	}
	m=nil
}

//
//var filesPool = sync.Pool{
//	New: func() interface{} {
//		b := make([]string, 0)
//		return b
//	},
//}
type HtmlFiles struct {
	Files   []string
	Success bool
	Err     error
}

//var files []string = make([]string, 0, 1<<13) // 不超过8192,可以重复利用
//var syncFile *sync.Mutex = new(sync.Mutex)
//
//func GetFiles() []string {
//	syncFile.Lock()
//	defer func() {
//		syncFile.Unlock()
//	}()
//	return files
//}
//func Clearfiles() {
//	syncFile.Lock()
//	defer func() {
//		syncFile.Unlock()
//	}()
//	files = files[0:0]
//}
//
//func GetFileList(path *string) (err error) {
//	//files := make([]string, 0)
//	//tfiles := filesPool.Get().([]string)
//
//	f, err := os.Stat(*path)
//	if err != nil {
//		err = errors.New(fmt.Sprintf("path: %s ,Err:%s", *path, err))
//		return err
//	}
//	if !f.IsDir() {
//		err = errors.New(fmt.Sprintf("path: %s ,Not Dir", *path))
//		return err
//	}
//	syncFile.Lock()
//	defer func() {
//		syncFile.Unlock()
//	}()
//
//	err = filepath.Walk(*path, func(path string, f os.FileInfo, err error) error {
//		if f == nil {
//			return err
//		}
//		if f.IsDir() {
//			return nil
//		}
//		files = append(files, path)
//		return nil
//	})
//	if err != nil {
//		err = errors.New(fmt.Sprintf("path: %s ,err:%s", *path, err))
//		return err
//	}
//
//	//fmt.Printf("%+v,%d\n",tfiles,len(tfiles))
//	//if len(files) > 0 {
//	//	return err
//	//}
//	//fmt.Printf("%p\n",tfiles)
//	return nil
//}
//var Day, interval, num ,h,m,s uint64
//
//
//func GetCurrentTime() (*uint64, *uint64, *uint64, *uint64) {
//	return &Day, &h, &m, &s
//}
//
//func GetTime(sinter int64) {
//	getBasic(uint64(sinter), 86400)
//
//	oneDay()
//
//	s = uint64(sinter) - (Day * 86400) - (h * 3600) - (m * 60)
//
//}
//
//func oneDay()  { // 一天内的小时分钟
//	h = interval / 3600
//	if h == 0 {
//		m = interval / 60
//	} else {
//		interval = interval - (3600 * h)
//		m = interval / 60
//	}
//
//}
//
//func getBasic(sinter uint64, tt uint64) { // 是否不足1天, 有几天, 减去天数后的时间差
//	if num = sinter / tt; num > 0 {
//		interval = sinter - (num * tt)
//	} else {
//		num = 0
//		interval = sinter
//	}
//
//}
