package funcs

import (
	"database/sql"
	"fmt"
	"gitee.com/taojun319/tjtools/db/oracle"
	"gitee.com/taojun319/tjtools/file"
	"github.com/PuerkitoBio/goquery"
	"github.com/ccpaging/nxlog4go"
	_ "github.com/mattn/go-oci8"
	"haifei/syncHtml/g"
	"haifei/syncHtml/modules"
	"haifei/syncHtml/utils"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//var ParsingSync *sync.WaitGroup = new(sync.WaitGroup)

//var MoveBakFileChan chan string = make(chan string, 0)
//var MoveFailFileChan chan string = make(chan string, 0)
var HtmlFilesChan chan []string = make(chan []string, 0)

//var FusionFileLock = make(chan struct{}, 0)
var Log *nxlog4go.Logger
var ParsingSync *sync.WaitGroup = new(sync.WaitGroup)

// xxx inithtml 临时变量
//var dom *goquery.Document
//var tmpSel1 *goquery.Selection
//var tmpSel2 *goquery.Selection
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

func FindHtml() { // xxx 这里必须要放在for 里，不能routine， 等待执行结束在重新扫描
	if err := file.CheckDirPerm(g.GetConfig().HtmlfileReg); err != nil {
		Log.Fatalf("%s not perm\n", g.GetConfig().HtmlfileReg)
	}

	fs := file.NewFilesStruct(g.GetConfig().HtmlfileReg)

	for {

		if ConnErr := oracle.CheckOracleConn(&g.GetConfig().OracleDsn); ConnErr != nil {
			_ = Log.Error("oracle 连接失败%s,wait %d 秒后,重试\n", strings.Replace(ConnErr.Error(), "\n", "", 1), g.GetConfig().Timeinter)

		} else {

			//aa:=NewFilesStruct("/home/go/GoDevEach/works/haifei/syncHtml_test8/htmlData/")

			if err := fs.ScanningFiles(); err != nil {
				_ = Log.Error("%s\n", err.Error())

			} else {
				Log.Info("%6s %d 个文件\n", "本次读取到", fs.Size())
				if fs.Size() > 0 {
					//var h = OnceHtmlFusionFiles.Get().(*modules.HtmlFusionFiles)
					loopStartTime := time.Now()
					//debug.SetGCPercent(100)
					ParsingSync.Add(fs.Size())
					HtmlFilesChan <- fs.GetFiles()
					ParsingSync.Wait()
					_ = Log.Warn("本次处理 %d个文件,共用时 %.4f 秒\n", fs.Size(), time.Now().Sub(loopStartTime).Seconds())

				}
				//debug.SetGCPercent(5)
			}

			//fmt.Printf("%d,%p,%p\n",fs.Size(),aa,fs.FileAbs)
			//err := file.GetFileList(&g.GetConfig().HtmlfileReg) // 遍历目录 包括子目录
			//
			//if err != nil {
			//	_ = Log.Error("%s\n", err.Error())
			//
			//} else {
			//	Log.Info("%6s %d 个文件\n", "本次读取到", len(file.GetFiles()))
			//	if len(file.GetFiles()) > 0 {
			//		//var h = OnceHtmlFusionFiles.Get().(*modules.HtmlFusionFiles)
			//		loopStartTime := time.Now()
			//		//debug.SetGCPercent(100)
			//		ParsingSync.Add(len(file.GetFiles()))
			//		HtmlFilesChan <- file.GetFiles()
			//		ParsingSync.Wait()
			//		_ = Log.Warn("本次处理 %d个文件,共用时 %.4f 秒\n", len(file.GetFiles()), time.Now().Sub(loopStartTime).Seconds())
			//		file.Clearfiles()
			//
			//	}
			//	//debug.SetGCPercent(1)
			//
			//}

		}
		time.Sleep(time.Duration(g.GetConfig().Timeinter) * time.Second)
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
//}

//var d map[string]interface{} = make(map[string]interface{},0)
func WaitHtmlFile() {
	for {
		select {
		case Files := <-HtmlFilesChan:
			for _, sfile := range Files {
				ParsingHtmlFile(sfile) // xxx  为了释放内存， 要控制一个一个执行
				//<-FusionFileLock          // xxx 不放在后台，不需要这个 通道
			}
			// xxx 提取sync.pool中 剩余的对象 ,并清空
			modules.HtmlFusionFree.Get()
			DomFree.Get()
			SelFree.Get()
			SelFree.Get()
			FiledHeaderFree.Get()
			//Log.Debug(fmt.Sprintf("释放HtmlFusion对象:%p,Dom:%p,Sele1:%p,Sele1:%p,FiledHeader:%p\n", h, dom, Sele1,Sele2, FiledHeader))
			//h = nil
			//dom = nil
			//Sele1 = nil
			//Sele2 = nil
			//FiledHeader = nil
			//fmt.Println(len(d))
		}

	}
}

func ParsingHtmlFile(file string) {

	//OnceHtmlFusion.Put(HtmlWheel)

	Log.Info("begin covert file:%s\n", file)

	//defer func() {
	//	FusionFileLock <- struct{}{}
	//}()
	var htmlWheel *modules.HtmlFusion = modules.HtmlFusionFree.Get().(*modules.HtmlFusion)
	// 临时页面变量

	var dom *goquery.Document = DomFree.Get().(*goquery.Document)
	var tmpSel1 *goquery.Selection = SelFree.Get().(*goquery.Selection)
	var tmpSel2 *goquery.Selection = SelFree.Get().(*goquery.Selection)
	var FiledHeader []string = FiledHeaderFree.Get().([]string) // 每个页面的列头 //行的数据

	//fmt.Printf("%p,%v\n", htmlwheel, len(htmlwheel.FusionSqlArrs))
	defer func() {
		htmlWheel.FusionSqlArrs = htmlWheel.FusionSqlArrs[0:0]
		utils.Clear(&htmlWheel.FusionSqlArrs)

		utils.Clear(htmlWheel)
		modules.HtmlFusionFree.Put(htmlWheel)
		// xxx 清除变量struct，释放memory

		utils.Clear(dom)
		utils.Clear(tmpSel2)
		utils.Clear(tmpSel1)

		FiledHeaderFree.Put(FiledHeader)
		DomFree.Put(dom)
		SelFree.Put(tmpSel1)
		SelFree.Put(tmpSel2)

		ParsingSync.Done()
	}()
	//defer func() {
	//	HtmlWheel = nil
	//}()

	//fmt.Println(string(file[len(file)-4:]))
	if strings.ToUpper(file[len(file)-4:]) != "HTML" { // 初始化文件 结构体
		_ = Log.Error("file :%s, ExtName not html\n", file)
		//time.Sleep(time.Duration(3 * time.Second))
		//MoveFailFileChan <- file
		htmlWheel.HtmlFileAbs = file
		htmlWheel.MoveFile(&g.GetConfig().HtmlFailDir)
		return
	}

	db, err := sql.Open("oci8", g.GetConfig().OracleDsn)

	if err != nil {
		_ = Log.Error("sql conn err\n")
		return
	}
	htmlWheel.DB = db

	defer func() {
		_ = db.Close()
	}()

	InitReadHtml(file, htmlWheel, dom, tmpSel1, tmpSel2, FiledHeader)

	if !htmlWheel.Fail {
		_ = os.Setenv("NLS_LANG", "")
		htmlWheel.StmtSql()
		// todo HTML中重复数据，可以通过 去除 更新SQL中的 HTML FILE，INSERT TIME 两个字段进行 重复插入
		// todo 如果只使用 INSERT TIME  updateSql 基本没有作用
		htmlWheel.ResultLogMoveFile()
	} else {
		htmlWheel.ResultLogMoveFile()
	}

	//fmt.Println(len(*HtmlWheel.SqlArrs))
	//fmt.Println(HtmlWheel.InsertSuccess,HtmlWheel.UpdateSuccess)
	//fmt.Println(len(*HtmlWheel.FusionSqlArrs))

}

var DomFree = sync.Pool{
	New: func() interface{} {
		return &goquery.Document{}
	},
}

var SelFree = sync.Pool{
	New: func() interface{} {
		return &goquery.Selection{}
	},
}

var FiledHeaderFree = sync.Pool{
	New: func() interface{} {
		return make([]string, 0)
	},
}

//var DataMapFree = sync.Pool{
//	New: func() interface{} {
//		return make(map[string]interface{}, 25)
//	},
//}
//var FiledsHeader []string = make([]string, 17)  // 每个页面的列头 //行的数据

func InitReadHtml(htmlfile string, htmlwheel *modules.HtmlFusion, dom *goquery.Document,
	tmpSel1 *goquery.Selection, tmpSel2 *goquery.Selection, FiledsHeader []string) {

	htmlwheel.Fail = false
	htmlwheel.HtmlFileAbs = htmlfile
	htmlwheel.Cv = 11
	htmlwheel.ChineseFile = false

	//Data := make([][]string, 0)
	//FiledsHeader := new([]string) // 列头，序号为KEY
	//cv := 0

	f, err := os.OpenFile(htmlfile, os.O_RDONLY, 0755)
	if err != nil {
		//Log.Error("file open err:%s\n", err)
		htmlwheel.ResultLog = fmt.Sprintf("file open err:%s\n", err)
		htmlwheel.Fail = true
		return
	}
	defer func() {
		_ = f.Close()
	}()

	//// 临时页面变量
	//
	//var dom *goquery.Document = DomFree.Get().(*goquery.Document)
	//var tmpSel1 *goquery.Selection = SeleFree.Get().(*goquery.Selection)
	//var tmpSel2 *goquery.Selection = SeleFree.Get().(*goquery.Selection)
	//var FiledsHeader []string = FiledHeaderFree.Get().([]string) // 每个页面的列头 //行的数据

	//var FiledsHeader []string = make([]string, 17)  // 每个页面的列头 //行的数据
	var DataMap = make(map[string]interface{}, 25)

	//var DataMap  = DataMapFree.Get().(map[string]interface{}) //行的数据
	Log.Debug(fmt.Sprintf("使用HtmlFusion对象:%p,Dom:%p,tmpSel1:%p,tmpSel2:%p,FiledHeader:%p\n",
		htmlwheel, dom, tmpSel1, tmpSel2, FiledsHeader))

	defer func() {
		//FiledsHeader = FiledsHeader[0:0]
		//FiledHeaderFree.Put(FiledsHeader)
		//DataMapFree.Put(DataMap)
		//dom.Empty()
		//tmpSel1.Empty()
		//tmpSel2.Empty()
		//utils.Clear(dom)
		//utils.Clear(tmpSel2)
		//utils.Clear(tmpSel1)
		//
		//DomFree.Put(dom)
		//SeleFree.Put(tmpSel1)
		//SeleFree.Put(tmpSel2)
		//tmpSel1 = nil
		//dom = nil
		//tmpSel2 = nil
	}()

	dom, err = goquery.NewDocumentFromReader(f)
	if err != nil {
		//log.Fatal(err)
		//Log.Error("dom err:%s\n", err)
		htmlwheel.ResultLog = fmt.Sprintf("dom err:%s", err)
		htmlwheel.Fail = true
		return
	}

	tmpSel1 = dom.Find("body:contains(报)")

	//fmt.Println(htmlfile, "zh:", zhStr.Length())
	if tmpSel1.Length() >= 1 {
		//Log.Info("Html:%s ; 进入中文处理线程\n", htmlfile)
		htmlwheel.ChineseFile = true
		// xxx 不处理中文文件
		htmlwheel.Fail = true
		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "文件是中文格式")
		return
	}

	tmpSel1 = dom.Find("body:contains(存储策略拷贝中的作业报告)") // 不正常
	tmpSel2 = dom.Find("body:contains(未找到符合选定条件的数据)") // 不正常
	if tmpSel1.Length() >= 1 || tmpSel2.Length() >= 1 {
		htmlwheel.Fail = true
		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "不正确的中文格式")
		return
	}

	tmpSel1 = dom.Find("body:contains(CommCell)")

	//tmpSel1.EachWithBreak(func(i int, selection *goquery.Selection) bool{
	//	if i == 0 {
	//		//fmt.Printf("CommCell :%+v\n", selection.Text())
	//		//htmlwheel.CommCell = strings.Split(strings.Split(selection.Text(), "CommCell:")[1], "--")[0]
	//		htmlwheel.CommCell = selection.Text()[0:250]
	//		//fmt.Printf("CommCell :%+v\n", ss)
	//		//CommCell = strings.TrimSpace(ss)
	//		htmlwheel.CommCell = strings.Split(htmlwheel.CommCell, "CommCell:")[1]
	//		htmlwheel.CommCell = strings.Split(htmlwheel.CommCell,"--")[0]
	//		selection.Empty()
	//		return false
	//	}
	//	return true
	//
	//})

	htmlwheel.CommCell = dom.Text()
	//fmt.Printf("CommCell :%+v\n", ss)
	//CommCell = strings.TrimSpace(ss)
	htmlwheel.CommCell = strings.Split(htmlwheel.CommCell, "CommCell:")[1]
	htmlwheel.CommCell = strings.Split(htmlwheel.CommCell, "--")[0]
	htmlwheel.CommCell = strings.Replace(htmlwheel.CommCell, " ", "", -1)

	tmpSel1 = dom.Find("body > table:nth-child(13) > tbody > tr:nth-child(1) > td")
	htmlwheel.TbodyDataNum = 13

	if len(tmpSel1.Nodes) == 0 {
		htmlwheel.Cv = 101 // 版本是10,HTML结构不同,列头一样
		// 特殊的CV10
		//t_tbodyHeader = dom.Find("body > table > tbody:nth-child(1) > tr:nth-child(1) > td[bgcolor='#cccccc']"  )
		tmpSel1 = dom.Find("body > table:nth-child(19) > tbody > tr:nth-child(1) > td")
		htmlwheel.TbodyDataNum = 19
	}

	//cv8 := false
	//有两种格式的HTML ,代表是CV8
	if len(tmpSel1.Nodes) == 0 {
		htmlwheel.Cv = 8
		tmpSel1 = dom.Find("body > table:nth-child(12) > tbody > tr:nth-child(1) > td")
		htmlwheel.TbodyDataNum = 12
	}

	// xxx 如果前面都没有列头，goquery前22个元素 找到列头的位置
	if tmpSel1.Length() == 0 {
		for i := 0; i <= 22; i++ {
			tmpSel1 = dom.Find("body > table:nth-child(" + strconv.Itoa(i) + ") > tbody > tr:nth-child(1) > td")
			//fmt.Printf("%d,%d\n", i, len(t_tbodyHeader.Nodes))\
			if len(tmpSel1.Nodes) > 9 { // xxx 找到一个至少有9列的行
				htmlwheel.TbodyDataNum = i
				htmlwheel.Cv = 100
				break
			}
		}
	}

	if len(tmpSel1.Nodes) == 0 {
		htmlwheel.Fail = true
		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "前22个元素，没有找到有效列头的表格")
		return

	}

	// 列头，序号为KEY
	// 循环找出列头
	var FindFiledHeader = func(i int, tmpSel2 *goquery.Selection) {
		defer func() {
			tmpSel2 = nil
		}()
		sa := tmpSel2.Text()
		modules.FormatFields(i, &sa, htmlwheel.Cv) // todo 格式化 列头
		//*FiledsHeader = append(*FiledsHeader, sa)
		//fmt.Println(i)
		FiledsHeader = append(FiledsHeader, sa)
		//FiledsHeader[i] = sa
	}

	tmpSel1.Each(FindFiledHeader)
	//fmt.Println(FiledsHeader)
	// 如果是CV8版本，手动修改列名, 删除 size of Backup cols 这列

	switch htmlwheel.Cv {
	case 8:
		//tmp := make([]string, 16)
		//copy(tmp, FiledsHeader)
		index := 9 // xxx cv8需要删除第9列
		FiledsHeader = append(FiledsHeader[0:index], FiledsHeader[index+1:]...)
		//FiledsHeader = tmp[0:9]
		//fmt.Println(*FiledsHeader)
		//*FiledsHeader = append(*FiledsHeader, []string{"Data Transferred", "Data Written"}...)
		//fmt.Println(*FiledsHeader)
		//FiledsHeader = append(*FiledsHeader, tmp[10:]...)
		//fmt.Println(*FiledsHeader)
		FiledsHeader[3] = "Job ID (CommCell)(Status)"

		break
	case 101:
		break
	default:
		break
	}

	switch true {
	case len(FiledsHeader) < g.FieldsLen-6: // col not enough,运维报告中不用判断

		htmlwheel.Fail = true

		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlfile, "col not enough")

		return
		//case (*FiledsHeader)[0] != "DATACLIENT": // col not English
		//
		//	htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlfile, "col not English")
		//	htmlwheel.Fail = true
		//	return
	}

	//var GenTimeSource *goquery.Selection
	// 生成的时间
	if htmlwheel.ChineseFile {
		tmpSel1 = dom.Find("body:contains(备份作业)") // 正常
	} else {
		tmpSel1 = dom.Find("body:contains(on)")
	}

	var SearchGenTimeCommCell = func(i int, tmpSel2 *goquery.Selection) bool {
		defer func() {
			tmpSel2 = nil
		}()
		if i == 0 {
			htmlwheel.GenTime = tmpSel2.Text()
			switch true {
			case htmlwheel.Cv == 8 && !htmlwheel.ChineseFile: //cv8 英文
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "generated on")[1]
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "CommCell ID:")[0]
				break
			case htmlwheel.ChineseFile && htmlwheel.Cv != 8: //除cv8 中文
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "备份作业摘要报告")[1]
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "上")[0]
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "\n")[1]
				break

			case htmlwheel.ChineseFile && htmlwheel.Cv == 8: //cv8 中文

				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "生成于")[1]
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "CommCell")[0]

				break
			default:
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "generated on")[1]
				htmlwheel.GenTime = strings.Split(htmlwheel.GenTime, "Version")[0]
			}

			htmlwheel.GenTime = strings.TrimSpace(htmlwheel.GenTime)
			htmlwheel.GenTime = strings.Replace(htmlwheel.GenTime, "\n", "", -1)
			return false
		}
		return true
	}

	tmpSel1.EachWithBreak(SearchGenTimeCommCell)

	//fmt.Println(htmlwheel.GenTime)
	//fmt.Println(htmlwheel.CommCell)
	//fmt.Println(*FiledsHeader)
	// 与 定义好的 9列进行比对， 这9列必须存在
	//fmt.Println(htmlwheel)
	if !htmlwheel.ChineseFile { // 中文文件不比对
		for i := 0; i < len(g.NeedFields); i++ {
			//fmt.Println((*FiledsHeader)[i])
			if g.NeedFields[i] != (FiledsHeader)[i] {
				if !strings.Contains((FiledsHeader)[i], g.NeedFields[i]) { // 包含关键字即可
					//fmt.Printf("%+v\n", NeedFields[i])
					htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlfile, "NeedField not enough")
					htmlwheel.Fail = true
					return
				}

			}
		}
	}

	//添加自定义列
	//*FiledsHeader = append(*FiledsHeader, []string{"HTMLFILE", "INSERTTIME", "COMMCELL", "APPLICATIONSIZE", "DATASUBCLIENT", "START TIME"}...)

	//var t_tbodyData *goquery.Selection
	var DetailFields map[int]string //这里与表格数据进行组合生成

	switch htmlwheel.Cv {
	case 8:
		//t_tbodyData = dom.Find("body > table:nth-child("+strconv.Itoa(htmlwheel.TbodyDataNum)+") > tbody > tr")
		DetailFields = g.DetailFieldsMapCv8
		break
	case 101:
		//t_tbodyData = dom.Find("body > table:nth-child("+strconv.Itoa(htmlwheel.TbodyDataNum)+") > tbody > tr")
		DetailFields = g.DetailFieldsMap
		break
	case 11:
		DetailFields = g.DetailFieldsMap
		break
	case 100: // xxx 自定义列头 ，先使用 默认列头
		DetailFields = g.DetailFieldsMap
		break
	default:
		//t_tbodyData = dom.Find("body > table:nth-child("+strconv.Itoa(htmlwheel.TbodyDataNum)+") > tbody > tr ")
		DetailFields = g.DetailFieldsMap
	}

	//var DetailSqlArr []map[string]string

	tmpSel1 = dom.Find("body > table:nth-child(" + strconv.Itoa(htmlwheel.TbodyDataNum) + ") > tbody > tr ")

	htmlwheel.TotalRow = tmpSel1.Length() - 1

	htmlwheel.FusionSqlArrs = make([]*modules.HtmlRowData, tmpSel1.Length()-1)
	//fmt.Println(len(htmlwheel.FusionSqlArrs))
	//fmt.Println(tmpSel1.Length())

	tmpSel1.Each(func(rowIndex int, rowsele *goquery.Selection) {
		if rowIndex == 0 { // 0行是列名
			return
		}

		//var rowReasonforfailure string //行 失败原因
		//var Color string = "0"         // 默认不显示颜色

		cols := rowsele.Find("td")

		if htmlwheel.Cv != 8 {
			if len(DetailFields) != len(cols.Nodes) {
				//Log.Error("htmfile:%s ,no eq rows:%d\n", htmlfile, i)
				htmlwheel.FailRow++
				return
			}
		} else {
			if len(DetailFields) != len(cols.Nodes)-1 { //
				//Log.Error("htmfile:%s ,no eq rows:%d\n", htmlfile, i)
				htmlwheel.FailRow++
				return
			}

		}

		RowWheel := modules.HtmlRowDataFree.Get().(*modules.HtmlRowData)

		// 先用上面判断 列数是否够
		// xxx 如果不在正常颜色中， 写入2  代表 需要颜色

		if rowColor, exists := rowsele.Attr("bgcolor"); exists {
			//fmt.Printf("%+v,%s,%v\n",RowWheel,rowColor,exists)
			if _, e := g.SuccessColors[rowColor]; !e {
				RowWheel.Color = "2"
			}
		} else {
			RowWheel.Color = "2" // xxx color找不到暂时是2
		}

		cols.Each(func(colnum int, colsele *goquery.Selection) {
			if htmlwheel.Cv == 8 && colnum == 9 {
				return
			}
			ss1 := colsele.Text()
			//fmt.Println(selection.Attr("bgcolor"))
			//switch true {
			//case strings.Contains(ss1,"N/A"):
			//	ss1=strings.Replace(ss1,"N/A","",-1)
			//}

			ss1 = *RowWheel.FormatValues(colnum, &ss1, htmlwheel) // todo 格式化 数据,重复使ssl

			if len(rowsele.Next().Find(" td").Nodes) == 1 {

				//fmt.Println(i, rowsele.Next().Find(" td").Text())
				RowWheel.RowReasonForFailure = rowsele.Next().Find(" td").Text()

				//fmt.Println(rowReasonforfailure)

				if strings.Contains(RowWheel.RowReasonForFailure, "Fail") ||
					strings.Contains(RowWheel.RowReasonForFailure, "ER") || len(strings.TrimSpace(RowWheel.RowReasonForFailure)) > 5 {
					RowWheel.Color = "2"
				}

			} else {
				RowWheel.RowReasonForFailure = "NULL"
			}

			//
			//fmt.Println(strings.Split())
			DataMap[DetailFields[colnum]] = ss1
			//tmpSubData = append(tmpSubData, ss1)
			//fmt.Println(i,ss1)
			colsele.Empty()
		})
		DataMap[g.DetailFieldsMapPlus[17]] = htmlwheel.CommCell
		DataMap[g.DetailFieldsMapPlus[18]] = RowWheel.StartTime
		DataMap[g.DetailFieldsMapPlus[19]] = RowWheel.SubClient
		DataMap[g.DetailFieldsMapPlus[20]] = RowWheel.ApplicationSize
		DataMap[g.DetailFieldsMapPlus[21]] = htmlfile
		DataMap[g.DetailFieldsMapPlus[22]] = strconv.FormatInt(time.Now().Unix(), 10) + "000"

		//fmt.Printf("%v,%v\n", &rowReasonforfailure, rowReasonforfailure)

		// xxx 截取最大字符长度
		utils.Rebuildstr(&RowWheel.RowReasonForFailure)
		utils.Split(&RowWheel.RowReasonForFailure, 4000)

		//fmt.Printf("%v,%v\n", &rowReasonforfailure, rowReasonforfailure)
		//fmt.Println(Color)
		//fmt.Println(rowReasonforfailure)
		DataMap[g.DetailFieldsMapPlus[23]] = strings.TrimSpace(RowWheel.RowReasonForFailure)
		DataMap[g.DetailFieldsMapPlus[24]] = RowWheel.Color
		DataMap[g.DetailFieldsMapPlus[25]] = RowWheel.RunType

		//DataMap = DataMap // 重复利用DataMap 每行必须要循环生成 SQL
		//htmlwheel.FusionSqlArrs = append(htmlwheel.FusionSqlArrs, RowWheel)
		htmlwheel.FusionSqlArrs[rowIndex-1] = RowWheel
		//rowsele.Empty()

		RowWheel.ToUpdateSql(htmlwheel, DataMap)
		RowWheel.ToInsertSql(htmlwheel, DataMap)

		// xxx 清空MAP
		utils.ClearMap(DataMap)
		rowsele.Empty()
		//d[fmt.Sprintf("%p",RowWheel)]=1
	})

	return
}
