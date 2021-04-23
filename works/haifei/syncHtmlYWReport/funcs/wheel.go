package funcs

import (
	nsema "gitee.com/taojun319/tjtools/control"
	"gitee.com/taojun319/tjtools/db/oracle"
	tu "gitee.com/taojun319/tjtools/file"
	"github.com/PuerkitoBio/goquery"
	"github.com/ccpaging/nxlog4go"
	"haifei/syncHtmlYWReport/g"
	"haifei/syncHtmlYWReport/modules"
	"haifei/syncHtmlYWReport/utils"
	"os"
	"strings"
	"sync"
	"time"
)

//var ParsingSync *sync.WaitGroup = new(sync.WaitGroup)

//var MoveBakFileChan chan string = make(chan string, 0)
//var MoveFailFileChan chan string = make(chan string, 0)
var HtmlFilesChan chan *tu.SelectFiles = make(chan *tu.SelectFiles, 0)

//var FusionFileLock = make(chan struct{}, 0)
var Log *nxlog4go.Logger

var ParsingSync *sync.WaitGroup = new(sync.WaitGroup)

// xxx inithtml 临时变量
//var dom *goquery.Document
//var tmpSele1 *goquery.Selection
//var tmpSele2 *goquery.Selection
//var DataMap = make(map[string]interface{},25)

//var c *g.Config

//var DataMap = make(map[string]interface{}, 25)

func LoadLogAndCfg() {
	Log = g.GetLog()
	//c = g.GetConfig()
}

//var OnceHtmlFusionFiles = sync.Pool{
//	New: func() interface{} {
//		return &modules.HtmlFusionFiles{ParsingSync: new(sync.WaitGroup)}
//	},
//}

//var files []string = make([]string, 0, 10000) // 不超过10000,重复利用
//var files = make([]string, 0)

var sema *nsema.Semaphore = nsema.NewSemaphore(2)

func FindHtml() { // xxx 这里必须要放在for 里，不能routine， 等待执行结束在重新扫描
	if err := tu.CheckDirPerm(g.GetConfig().HtmlFileReg); err != nil {
		Log.Fatalf("%s not perm\n", g.GetConfig().HtmlFileReg)
	}

	for {
		//
		if ConnErr := oracle.CheckOracleConn(&g.GetConfig().OracleDsn); ConnErr != nil {
			_ = Log.Error("oracle 连接失败%s,wait %d 秒后,重试\n", strings.Replace(ConnErr.Error(), "\n", "", 1), g.GetConfig().TimeInter)

		} else {

			//var f = &tu.SelectFiles{}

			//defer tu.SelectFilesFree.Put(f)
			var f = &tu.SelectFiles{}
			f.Path = g.GetConfig().HtmlFileReg
			//f.Files = files
			err := f.GetFileList()
			//err := file.GetFileList(&g.GetConfig().HtmlfileReg) // 遍历目录 包括子目录

			if err != nil {
				_ = Log.Error("%s\n", err.Error())

			} else {
				//Log.Info("本次读取到 %d 个文件,结构体内存地址: %p,文件切片内存地址: %p,Len: %d, Cap: %d\n", f.Len(),
				//	&f,
				//	f.GetFiles(),
				//	f.Len(),
				//	cap(f.GetFiles()))
				if len(f.GetFiles()) > 0 {
					//var h = OnceHtmlFusionFiles.Get().(*modules.HtmlFusionFiles)
					loopStartTime := time.Now()
					//debug.SetGCPercent(100)
					modules.GetDBConn()

					ParsingSync.Add(int(f.Len()))
					HtmlFilesChan <- f
					ParsingSync.Wait()
					modules.CloseConn()

					_ = Log.Warn("本次处理 %d个文件,共用时 %.4f 秒\n", len(f.GetFiles()), time.Now().Sub(loopStartTime).Seconds())
					f.Cleanfiles()
				} else {
					//debug.SetGCPercent(50)
				}

			}

			//err := file.GetFileList(&g.GetConfig().HtmlfileReg) // 遍历目录 包括子目录
			//
			//if err != nil {
			//	_ = Log.Error("%s\n", err.Error())
			//
			//} else {
			//	Log.Info("本次读取到 %d 个文件,内存地址: %p\n", len(file.GetFiles()),file.GetFiles())
			//	if len(file.GetFiles()) > 0 {
			//		//var h = OnceHtmlFusionFiles.Get().(*modules.HtmlFusionFiles)
			//		loopStartTime := time.Now()
			//		debug.SetGCPercent(100)
			//		ParsingSync.Add(len(file.GetFiles()))
			//		HtmlFilesChan <- file.GetFiles()
			//		ParsingSync.Wait()
			//		_ = Log.Warn("本次处理 %d个文件,共用时 %.4f 秒,地址: %p\n", len(file.GetFiles()), time.Now().Sub(loopStartTime).Seconds())
			//		file.Cleanfiles()
			//	}else{
			//		debug.SetGCPercent(100)
			//	}
			//
			//}
		}
		time.Sleep(time.Duration(g.GetConfig().TimeInter) * time.Second)

	}
}

//type HtmlFusionFiles struct {
//	hs          []string
//	ParsingSync *sync.WaitGroup
//}
//var filesPool = sync.Pool{
//	New: func() interface{} {
//		return make([]string, 0, 0)
//	},
//}

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
//	return nil
//}

//func WaitMoveFile() {
//
//	for {
//		select {
//		case f, ok := <-MoveBakFileChan:
//			if ok {
//				go MoveFile(&f, &c.HtmlBakDir)
//			} else {
//				break
//			}
//		case f, ok := <-MoveFailFileChan:
//			if ok {
//				go MoveFile(&f, &c.HtmlFailDir)
//			} else {
//				break
//			}
//		}
//	}
func runDo(sfile string) {
	defer func() {
		sema.Release()
		ParsingSync.Done()
	}()
	ParsingHtmlFile(sfile) // xxx  为了释放内存， 要控制执行
	if modules.DB.Stats().OpenConnections > 10 {
		_ = Log.Warn("DB Status:%+v\n", modules.DB.Stats())
	}
}
func WaitHtmlFile() {

	for {
		select {
		case Files := <-HtmlFilesChan:
			//modules.GetDBConn()

			for _, sfile := range Files.GetFiles() {
				sema.Acquire()

				go runDo(sfile)

			}
			//modules.HtmlFusionFree.Get()
			//modules.DomFree.Get()
			//modules.SeleFree.Get()
			//modules.SeleFree.Get()
			//modules.SeleFree.Get()
			//modules.SeleFree.Get()
			//modules.FiledHeaderFree.Get()

			//modules.CloseConn()
		}
	}
}

func ParsingHtmlFile(file string) {

	//OnceHtmlFusion.Put(HtmlWheel)

	Log.Info("begin covert file:%s\n", file)

	//defer func() {
	//	FusionFileLock <- struct{}{}
	//}()
	var htmlwheel *modules.HtmlFusion = modules.HtmlFusionFree.Get().(*modules.HtmlFusion)

	//var dom *goquery.Document = modules.DomFree.Get().(*goquery.Document)
	//var tmpSele1 *goquery.Selection = modules.SeleFree.Get().(*goquery.Selection)
	//var tmpSele2 *goquery.Selection = modules.SeleFree.Get().(*goquery.Selection)
	//var tmpSele3 *goquery.Selection = modules.SeleFree.Get().(*goquery.Selection)
	//var FiledsHeader []string = modules.FiledHeaderFree.Get().([]string)
	var Summary_tbodyHeader = modules.SeleFree.Get().(*goquery.Selection)

	defer func() {
		//utils.Clear(htmlwheel)
		utils.Clear(htmlwheel)
		modules.HtmlFusionFree.Put(htmlwheel)
		// xxx 清除变量struct，释放memory

		//FiledsHeader = FiledsHeader[0:0]
		//modules.FiledHeaderFree.Put(FiledsHeader)

		//utils.Clear(dom)
		//utils.Clear(tmpSele2)
		//utils.Clear(tmpSele1)
		//utils.Clear(tmpSele3)

		utils.Clear(Summary_tbodyHeader)

		//modules.DomFree.Put(dom)
		//modules.SeleFree.Put(tmpSele1)
		//modules.SeleFree.Put(tmpSele2)
		//modules.SeleFree.Put(tmpSele3)
		modules.SeleFree.Put(Summary_tbodyHeader)

	}()

	if strings.ToUpper(file[len(file)-4:]) != "HTML" { // 初始化文件 结构体
		_ = Log.Error("file :%s, ExtName not html\n", file)
		//time.Sleep(time.Duration(3 * time.Second))
		//MoveFailFileChan <- file
		htmlwheel.HtmlFileAbs = file
		htmlwheel.MoveFile(&g.GetConfig().HtmlFailDir)
		return
	}

	//db, err := sql.Open("oci8", g.GetConfig().OracleDsn)
	//
	//if err != nil {
	//	_ = Log.Error("sql conn err\n")
	//	return
	//}
	//htmlwheel.DB = db
	//
	//defer func() {
	//	_ = db.Close()
	//}()
	htmlwheel.HtmlFileAbs = file

	htmlwheel.InitReadHtml(Summary_tbodyHeader)

	if !htmlwheel.Fail {
		_ = os.Setenv("NLS_LANG", "")

		// 插入详细数据
		htmlwheel.StmtSql(false)

		//htmlwheel.SortKeys = htmlwheel.SortKeys[0:0]

		// 生成摘要数据 xxx 不记录 summary数据
		//htmlwheel.InitSummary()
		//if !htmlwheel.Fail {
		//	htmlwheel.StmtSql(true)
		//}

		// todo HTML中重复数据，可以通过 去除 更新SQL中的 HTMLFILE，INSERTTIME 两个字段进行 重复插入
		// todo 如果只使用 INERRTTIME  updateSql 基本没有作用
		htmlwheel.ResultLogMoveFile()
	} else {

		htmlwheel.ResultLogMoveFile()
	}

	//fmt.Println(len(*HtmlWheel.SqlArrs))
	//fmt.Println(HtmlWheel.InsertSuccess,HtmlWheel.UpdateSuccess)
	//fmt.Println(len(*HtmlWheel.FusionSqlArrs))

}

//var DataMapFree = sync.Pool{
//	New: func() interface{} {
//		return make(map[string]interface{}, 25)
//	},
//}
//var FiledsHeader []string = make([]string, 17)  // 每个页面的列头 //行的数据
