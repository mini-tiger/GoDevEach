package modules

import (
	"fmt"
	"gitee.com/taojun319/tjtools/file"
	"github.com/PuerkitoBio/goquery"
	"github.com/ccpaging/nxlog4go"
	"haifei/syncHtmlYWReport/g"
	"haifei/syncHtmlYWReport/utils"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//type HtmlFusionFiles struct {
//	FileNames   []string
//	ParsingSync *sync.WaitGroup
//}
//
//func (h *HtmlFusionFiles) GetFileName() []string {
//	return h.FileNames
//}

type HtmlFusion struct {
	HtmlFileAbs    string         // 文件名绝对路径
	NewHtmlFileAbs string         //备份文件名绝对路径
	Cv             int            //文件版本
	ResultLog      string         //错误提示字符串
	Fail           bool           // 是否成功，流程是否往下进行
	FusionSqlArrs  []*HtmlRowData //字段和数据的MAP
	//SqlArrs        []map[string]string//SQL语句MAP，分插入和更新
	//SqlArrs        []map[string]string //SQL语句MAP，分插入和更新
	//SortKeys    []string
	ChineseFile bool
	//DetailSqlArr  [](*map[string]string)
	FiledsHeader []string

	CommCell            string
	GenTime             string
	DetailInsertSuccess uint32 // 插入和更新 成功多少条
	DetailUpdateSuccess uint32
	DetailInsertFail    uint32
	DetailUpdateFail    uint32
	DetailFailRow       uint32 // 行数据不正常的个数
	DetailTotalRow      int
	DetailTbodyDataNum  int // 数据在goquery 节点的位置

	SummaryInsertSuccess uint32 // 插入和更新 成功多少条
	SummaryUpdateSuccess uint32
	SummaryInsertFail    uint32
	SummaryUpdateFail    uint32
	SummaryFailRow       uint32 // 行数据不正常的个数
	SummaryTotalRow      int
	SummaryTbodyDataNum  int // 数据在goquery 节点的位置
	//DB *sql.DB

	Dom *goquery.Document
}

var HtmlFusionFree = sync.Pool{
	New: func() interface{} {
		return &HtmlFusion{}
	},
}

//var DataMap = make(map[string]interface{}, 28)

var Log *nxlog4go.Logger
var c *g.Config

//var RowSortKey = make([]string, 25)
var logstr string

func LoadLogAndCfg() {
	Log = g.GetLog()
	c = g.GetConfig()
}

func (h *HtmlFusion) InitSummary() {
	//var SummaryFields map[int]string //这里与表格数据进行组合生成
	breakFields := make([]int, 0)
	switch h.Cv {
	case 8:

		h.GenSummaryData(g.SummaryFieldsMapCv8, breakFields) // 生成 详细数据,version 11 使用默认DetailFieldsMap

		break
	case 10:

		h.GenSummaryData(g.SummaryFieldsMapCv10, breakFields) // 生成 详细数据,version 11 使用默认DetailFieldsMap
		break
	case 11:

		//SummaryFields = g.SummaryFieldsMap
		h.GenSummaryData(g.SummaryFieldsMap, breakFields) // 生成 详细数据,version 11 使用默认DetailFieldsMap

		break

	default:
		rn, breakFields, tmpfileds := h.GenSummaryFields() // 生成 摘要数据,version 1000 使用默认DetailFieldsMap
		if rn > 0 {

			//Log.Error(logstr)
			h.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "摘要列不能自动生成")
			h.Fail = true
			return
		}
		h.GenSummaryData(tmpfileds, breakFields) // 生成 详细数据,version 11 使用默认DetailFieldsMap
	}
}

func (h *HtmlFusion) GenSummaryFields() (int, []int, map[int]string) { // 生成摘要表的列名

	//var t_tbodyData *goquery.Selection
	var returnnum int = 0
	//t_tbodyData = htmlwheel.Dom.Find(domstr)
	//Data = make([][]string, 0)
	//Log.Printf("HtmlFile:%s,Summary rows:%d \n", HtmlFile, len(t_tbodyData.Nodes))
	tmpFields := make(map[int]string, 20)
	breakFields := make([]int, 20) // 需要跳过的列

	f, err := os.OpenFile(h.HtmlFileAbs, os.O_RDONLY, 0755)
	if err != nil { // 不应该主会有错误

		return 1, breakFields, tmpFields
	}
	defer func() {
		_ = f.Close()
	}()

	h.Dom, err = goquery.NewDocumentFromReader(f)
	if err != nil {
		return 1, breakFields, tmpFields
	}

	SummaryDomFindStr := "body > table:nth-child(" + strconv.Itoa(h.SummaryTbodyDataNum) + ") > tbody > tr:nth-child(2) > td"

	tmpSele := h.Dom.Find(SummaryDomFindStr)

	tmpSele.Each(func(colnum int, colsele *goquery.Selection) {
		ss1 := colsele.Text()
		//fmt.Println(selection.Attr("bgcolor"))
		//switch true {
		//case strings.Contains(ss1,"N/A"):
		//	ss1=strings.Replace(ss1,"N/A","",-1)
		//}

		if v, b := g.AllChEngSummaryMap.Get(ss1); b {
			tmpFields[colnum] = v.(string)
		} else {
			breakFields = append(breakFields, colnum)
		}

		//ss1 = formatDetailValues(colnum, ss1) // todo 格式化 数据
		//fmt.Println(ss1)
	})

	if len(tmpFields) == 0 {
		returnnum = 6
	}

	return returnnum, breakFields, tmpFields

}

func (h *HtmlFusion) GenSummaryData(FiledsHeader map[int]string, breakFields []int) {
	f, err := os.OpenFile(h.HtmlFileAbs, os.O_RDONLY, 0755)
	if err != nil { // 不应该主会有错误

		return
	}
	defer func() {
		_ = f.Close()
	}()

	h.Dom, err = goquery.NewDocumentFromReader(f)
	if err != nil {
		return
	}

	SummaryDomFindStr := "body > table:nth-child(" + strconv.Itoa(h.SummaryTbodyDataNum) + ") > tbody > tr"

	tmpSele := h.Dom.Find(SummaryDomFindStr)

	h.FusionSqlArrs = make([]*HtmlRowData, tmpSele.Length()-2)
	h.SummaryTotalRow = tmpSele.Length() - 2

	tmpSele.Each(func(rowIndex int, rowsele *goquery.Selection) {
		if rowIndex < 2 { // 0行是 摘要 1行是列名
			return
		}
		RowWheel := HtmlRowDataFree.Get().(*HtmlRowData)

		tmpSubData := make(map[string]interface{}, 24)
		cols := rowsele.Find("td")
		if len(FiledsHeader) != len(cols.Nodes) { //
			_ = Log.Error("htmfile:%s ,summary no eq rows:%d,version:%d\n", h.HtmlFileAbs, rowIndex, h.Cv)
			return
		}
		cols.Each(func(colnum int, colsele *goquery.Selection) {
			ss1 := colsele.Text()
			if b, _ := utils.Contain(colnum, breakFields); b { // 如果有需要跳过的列
				return
			}
			//fmt.Println(selection.Attr("bgcolor"))
			//switch true {
			//case strings.Contains(ss1,"N/A"):
			//	ss1=strings.Replace(ss1,"N/A","",-1)
			//}

			//ss1 = formatDetailValues(colnum, ss1) // todo 格式化 数据
			ss1 = strings.TrimSpace(ss1)
			if strings.Contains(FiledsHeader[colnum], "STARTTIME") {
				//ss1 = g.FormatTime(ss1)

				utils.FormatTime(&ss1, true)
				ss1 = fmt.Sprintf("%s000", ss1)

				if _, err := strconv.Atoi(ss1); err != nil {
					_ = Log.Warn("HTMLFILE:%s,rowline:%d,colline:%d,headr:%s, 转换数字失败,用0代替 实际字符:%s\n", h.HtmlFileAbs, rowIndex, colnum, "START TIME", ss1)
					ss1 = "0"
					return
				}
			}

			if strings.Contains(FiledsHeader[colnum], "FAILEOBJECT") || strings.Contains(FiledsHeader[colnum], "FAILEDFOLDERS") ||
				strings.Contains(FiledsHeader[colnum], "TOTALJOB") || strings.Contains(FiledsHeader[colnum], "COMPLETED") ||
				strings.Contains(FiledsHeader[colnum], "COMPLETEDWITHERRORS") ||
				strings.Contains(FiledsHeader[colnum], "COMPLETEDWITHWARNINGS") ||
				strings.Contains(FiledsHeader[colnum], "KILLED") ||
				strings.Contains(FiledsHeader[colnum], "UNSUCCESSFUL") ||
				strings.Contains(FiledsHeader[colnum], "RUNNING") ||
				strings.Contains(FiledsHeader[colnum], "DELAYED") ||
				strings.Contains(FiledsHeader[colnum], "NORUN") ||
				strings.Contains(FiledsHeader[colnum], "NOSCHEDULE") ||
				strings.Contains(FiledsHeader[colnum], "COMMITTED") ||
				strings.Contains(FiledsHeader[colnum], "COMMITTED") ||
				strings.Contains(FiledsHeader[colnum], "PROTECTEDOBJECTS") {
				ss1 = strings.Join(strings.Split(ss1, ","), "")
				if _, err := strconv.Atoi(ss1); err != nil {
					_ = Log.Error("HTMLFILE:%s,rowline:%d,rowline:%d,headr:%s, 转换数字失败", h.HtmlFileAbs, rowIndex, colnum, FiledsHeader[colnum])
					ss1 = "NULL"
					return
				}

			}

			//fmt.Println(i,colnum,FiledsHeader[colnum],strings.Contains(FiledsHeader[colnum],"Failed Folder"))
			//
			//fmt.Println(strings.Split())
			tmpSubData[FiledsHeader[colnum]] = ss1
			//tmpSubData = append(tmpSubData, ss1)
			//fmt.Println(i,ss1)
			//colsele.Empty()
		})
		tmpSubData[g.SummaryFieldsMapPlus[20]] = h.CommCell
		tmpSubData[g.SummaryFieldsMapPlus[21]] = h.GenTime
		tmpSubData[g.SummaryFieldsMapPlus[22]] = h.HtmlFileAbs
		tmpSubData[g.SummaryFieldsMapPlus[23]] = strconv.FormatInt(time.Now().Unix(), 10) + "000"
		//fmt.Println(tmpSubData)
		h.FusionSqlArrs[rowIndex-2] = RowWheel

		RowWheel.DataMap = tmpSubData
		RowWheel.ToUpdateSql(true)
		RowWheel.ToInsertSql(true)

		//SummarySqlArr = append(SummarySqlArr, &tmpSubData)
		// xxx 清空MAP
		utils.ClearMap(tmpSubData)
		//rowsele.Empty()
	})

	return

}

func (h *HtmlFusion) MoveFile(dir *string) {

	bakDir := filepath.Join(*dir, time.Now().Format("2006-01-02"))
	err := os.MkdirAll(bakDir, 0777)
	if err != nil {
		_ = Log.Error("mkdir Error DIR: %s, err:%s", bakDir, err)
	}

	//fmt.Printf("fileBak:%s\n", filepath.Join(bakDir, path.Base(oldfile)))
	//fmt.Println("fileBak:", filepath.Join(bakDir, filepath.Base(oldfile)), file.IsExist(filepath.Join(bakDir, filepath.Base(oldfile))))

	if file.IsExist(filepath.Join(bakDir, filepath.Base(h.HtmlFileAbs))) { // 备份文件夹中是否存在相同名字要 备份的文件

		NewName := utils.NewDateFileName(h.HtmlFileAbs, bakDir)

		//err := os.Rename(h.HtmlFileAbs, NewName)

		cmd := exec.Command("/bin/mv", h.HtmlFileAbs, NewName)
		_, err = cmd.Output()

		if err != nil {
			_ = Log.Error("osRename html File Error DIR: %s,cmd :%s, err:%s\n", bakDir, cmd.String(), err)
		} else {
			Log.Info("move file old:%s,new:%s\n", h.HtmlFileAbs, NewName)
		}
	} else {
		_, fileName := filepath.Split(h.HtmlFileAbs)
		Log.Info("move file old:%s,new:%s\n", h.HtmlFileAbs, filepath.Join(bakDir, fileName))
		//err = os.Rename(h.HtmlFileAbs, filepath.Join(bakDir, fileName))   // xxx os.Rename 不支持 跨硬盘

		cmd := exec.Command("/bin/mv", h.HtmlFileAbs, filepath.Join(bakDir, fileName))
		_, err = cmd.Output()

		if err != nil {
			_ = Log.Error("osRename html File Error DIR: %s,cmd :%s, err:%s\n", bakDir, cmd.String(), err)
		}
	}

}

func (h HtmlFusion) String() string {
	if len(h.ResultLog) > 5 {
		return fmt.Sprintf("%s.HtmFile:%s, version:%d,DetailTotalRows:%d,DetailFailRow:%d,DetailInsertSuccess:%d,"+
			"DetailUpdateSuccess:%d,DetailUpdateFail:%d,DetailInsertFail:%d,"+
			"SummaryTotalRow:%d,SummaryFailRow:%d,SummaryInsertSuccess:%d,SummaryInsertFail:%d,SummaryUpdateSuccess:%d,SummaryUpdateFail:%d,"+
			"CommCell:%s,GenTime:%s",
			h.ResultLog,
			filepath.Base(h.HtmlFileAbs),
			h.Cv,
			h.DetailTotalRow,
			h.DetailFailRow,
			h.DetailInsertSuccess,
			h.DetailUpdateSuccess,
			h.DetailUpdateFail,
			h.DetailInsertFail,
			h.SummaryTotalRow,
			h.SummaryFailRow,
			h.SummaryInsertSuccess,
			h.SummaryInsertFail,
			h.SummaryUpdateSuccess,
			h.SummaryUpdateFail,
			h.CommCell, h.GenTime)
	} else {
		return fmt.Sprintf("HtmFile:%s, version:%d,DetailTotalRows:%d,DetailFailRow:%d,DetailInsertSuccess:%d,"+
			"DetailUpdateSuccess:%d,DetailUpdateFail:%d,DetailInsertFail:%d,"+
			"SummaryTotalRow:%d,SummaryFailRow:%d,SummaryInsertSuccess:%d,SummaryInsertFail:%d,SummaryUpdateSuccess:%d,SummaryUpdateFail:%d,"+
			"CommCell:%s,GenTime:%s",
			filepath.Base(h.HtmlFileAbs),
			h.Cv,
			h.DetailTotalRow,
			h.DetailFailRow,
			h.DetailInsertSuccess,
			h.DetailUpdateSuccess,
			h.DetailUpdateFail,
			h.DetailInsertFail,
			h.SummaryTotalRow,
			h.SummaryFailRow,
			h.SummaryInsertSuccess,
			h.SummaryInsertFail,
			h.SummaryUpdateSuccess,
			h.SummaryUpdateFail,
			h.CommCell, h.GenTime)
	}
}

func (h *HtmlFusion) ResultLogMoveFile() {
	//file := h.HtmlFileAbs

	if h.Fail {
		//MoveFailFileChan <- h.HtmlFileAbs

		//Log.Error(h.ResultLog)
		_ = Log.Error("%s\n", h)

		h.MoveFile(&c.HtmlFailDir)
	} else {
		//Log.Info("finish covert file:%s\n", h.HtmlFileAbs)
		//MoveBakFileChan <- h.HtmlFileAbs

		Log.Info("%s\n", h)

		if h.DetailUpdateSuccess > 0 {
			_ = Log.Warn("!!!!!!!HtmlFile:%s, update gt 0\n", h.HtmlFileAbs)
		}
		h.MoveFile(&c.HtmlBakDir)

	}
}

func updatesliceToString(sl map[string]interface{}) (returnstring string) {
	var keys []string = make([]string, len(sl))
	var index int = 0
	for k, _ := range sl {
		//keys = append(keys, k)
		keys[index] = k
		index++

	}
	sort.Strings(keys)
	wherestring := ""
	colsStr := ""
	for i, v := range keys {
		if (sl)[v] == "NULL" || (sl)[v] == "" || v == "INSERTTIME" || v == "HTMLFILE" {
			continue
		}
		switch true {
		case i == len(sl)-1:
			//returnstring = returnstring
			//returnstring = returnstring + "\"" + sl1[i] + "\"" + ")"
			colsStr = colsStr + "\"" + v + "\"='" + (sl)[v].(string) + "'"
			wherestring = wherestring + "\"" + v + "\"='" + (sl)[v].(string) + "'"

			break
		case i < len(sl)-1:
			//returnstring = returnstring
			colsStr = colsStr + "\"" + v + "\"='" + (sl)[v].(string) + "',"
			wherestring = wherestring + "\"" + v + "\"='" + (sl)[v].(string) + "' and "
			break
		}

	}

	// fixme 详细表, wheresql只保留 COMMCELL  Jobid, 替换上面的 where

	returnstring = colsStr + " where " + wherestring
	return
}

func (h *HtmlFusion) StmtSql(IsSummary bool) {
	if h.Fail {
		return
	}
	// 判断数据库是否正常
	//if ConnErr := oracle.CheckOracleLive(DB); ConnErr != nil {
	//	GetDBConn()
	//}
	//if len(os.Args) != 2 {
	//	log.Fatalln(os.Args[0] + " user/password@host:port/sid")
	//}

	//db, err := sql.Open("oci8", c.OracleDsn)

	//if err != nil {
	//	_ = Log.Error("sql conn err\n")
	//	return
	//}

	//defer func() {
	//	_ = db.Close()
	//}()

	//updatenum := 0
	//insertnum := 0

	if len(h.FusionSqlArrs) == 0 {
		logstr = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "sql Gen err,sqlarr len is 0")
		//Log.Fatalf(logstr)
		h.ResultLog = logstr
		h.Fail = true
		return
	}

	//ss := []string{"to_date('2019-08-25 22:11:11','yyyy-mm-dd hh24:mi:ss')"}

	for _, RowStruct := range h.FusionSqlArrs {

		//Log.Error("%s\n", (sqlArr)[i]["update"])
		//fmt.Println(fmt.Sprintf(sqlArr[i]["insert"]+",%s)", StartTimeArr[i]))
		//fmt.Println(sqlArr[i]["insert"]+",:1)",StartTimeArr[i])
		//r, err := db.Exec(fmt.Sprintf("insert into HF_BACKUPDETAIL(\"START TIME\") values(%s)", "to_date('2019-08-25 22:11:11','yyyy-mm-dd hh24:mi:ss')"))
		//fmt.Println(r,err)
		//ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second)
		//r,err:=db.ExecContext(ctx,sqlArr[i]["insert"]+",:1)",StartTimeArr[i])
		//cancel()

		if RowStruct == nil {
			continue
		}
		//Log.Error("%s\n", *RowStruct.UpdateSql)
		Result, err := DB.Exec(RowStruct.UpdateSql.String())
		//fmt.Println(Result.LastInsertId())
		//fmt.Println(Result.RowsAffected())
		//fmt.Println("===========", r, err)
		if err != nil {
			_ = Log.Error("HtmlFile:%s ,Sql Exec Err:%s,,,Updatesql:%s\n",
				h.HtmlFileAbs, strings.Replace(err.Error(), "\n", "", 1), RowStruct.UpdateSql.String())
			h.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "insert or update have err")
			h.Fail = true
			continue
		}

		r, err := Result.RowsAffected()
		if err != nil {
			_ = Log.Error("HtmlFile:%s ,Sql Exec Err:%s\n", h.HtmlFileAbs, strings.Replace(err.Error(), "\n", "", 1))
			h.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "insert or update have err")
			h.Fail = true
			if IsSummary {
				//h.SummaryUpdateFail++
				atomic.AddUint32(&h.SummaryUpdateFail, 1)
			} else {
				//h.DetailUpdateFail++
				atomic.AddUint32(&h.DetailUpdateFail, 1)
			}
			continue
		}

		if r == 0 {
			RowStruct.ToInsertSql(false)
			_, err := DB.Exec(RowStruct.InsertSql.String())
			if err != nil {
				_ = Log.Error("HtmlFile: %s ,\nSql Exec Err: %s,\nSql: %s\n", h.HtmlFileAbs,
					strings.Replace(err.Error(), "\n", "", 1), RowStruct.InsertSql.String())
				h.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "insert or update have err")
				h.Fail = true

				if IsSummary {
					//h.SummaryInsertFail++
					atomic.AddUint32(&h.SummaryInsertFail, 1)
				} else {
					//h.DetailInsertFail++
					atomic.AddUint32(&h.DetailInsertFail, 1)
				}

				continue
			}
			if IsSummary {
				//h.SummaryInsertSuccess++
				atomic.AddUint32(&h.SummaryInsertSuccess, 1)
			} else {
				//h.DetailInsertSuccess++
				atomic.AddUint32(&h.DetailInsertSuccess, 1)
			}

		} else {
			//Log.Warn("%+v\n", (sqlArr)[i]["update"])
			if IsSummary {
				//h.SummaryUpdateSuccess++
				atomic.AddUint32(&h.SummaryUpdateSuccess, 1)
			} else {
				//h.DetailUpdateSuccess++
				atomic.AddUint32(&h.DetailUpdateSuccess, 1)
			}

		}
		// xxx clear RowStruct

		//RowStruct.UpdateSql = nil
		//RowStruct.InsertSql = nil
		RowStruct.InsertSql.Reset()
		RowStruct.UpdateSql.Reset()
		utils.Clear(RowStruct)
		HtmlRowDataFree.Put(RowStruct)

	}

	//Log.Info("htmlfile : %s ,total sql:%d,success update sql:%d,insert sql:%d\n", filepath.Base(htmlfile), len(*sqlArr), updatenum, insertnum)
	return
}

var DomFree = sync.Pool{
	New: func() interface{} {
		return &goquery.Document{}
	},
}

var SeleFree = sync.Pool{
	New: func() interface{} {
		return &goquery.Selection{}
	},
}

var FiledHeaderFree = sync.Pool{
	New: func() interface{} {
		return make([]string, 0)
	},
}

func (self *HtmlFusion) formatFields(index int, s *string) {
	//rstr := ""
	switch true {
	case index == 0 && strings.Contains(*s, "Client"):
		*s = "DATACLIENT"
		break
	case strings.Contains(*s, "Phase(Write End Time)"):

		*s = strings.Replace(*s, "(Write End Time)", "", -1)

		break
	case index == 7 && self.Cv == 101:

		*s = "End Time or Current Phase" // cv 101
		break

	case index == 10 && self.Cv == 101:
		*s = "Data Written"
		break

	case strings.Contains(*s, "(Compression") && strings.Contains(*s, "Rate)"):
		//rstr = strings.Replace(s, " (Compression Rate)", "", -1)
		*s = "Size of Application"
		break
	case strings.Contains(*s, "Data Transferred"):
		//rstr = strings.Replace(s, " (Compression Rate)", "", -1)
		*s = "Data Transferred"
		break
	case strings.Contains(*s, "(Space Saving Percentage)"):
		*s = strings.Replace(*s, "(Space Saving Percentage)", "", -1)
		break

	case strings.Contains(*s, "(Current)"):
		*s = strings.Replace(*s, "(Current)", "", -1)
		break
	case *s == "Agent /Instance":
		*s = "AgentInstance"
		break
	case *s == "Backup Set /Subclient":
		*s = "BackupSetSubclient"
		break

	default:
		*s = strings.TrimSpace(*s)
	}
	//fmt.Println(rstr)
	*s = strings.TrimSpace(*s)
	*s = strings.Replace(*s, "\n", "", -1)

	//return &rstr
}

func (htmlwheel *HtmlFusion) InitReadHtml(Summary_tbodyHeader *goquery.Selection) {

	htmlwheel.Fail = false
	//htmlwheel.HtmlFileAbs = *htmlfile
	htmlwheel.Cv = g.Cv11
	htmlwheel.ChineseFile = true

	//Data := make([][]string, 0)
	//FiledsHeader := new([]string) // 列头，序号为KEY
	//cv := 0

	f, err := os.OpenFile(htmlwheel.HtmlFileAbs, os.O_RDONLY, 0755)
	fs, _ := f.Stat()
	//fmt.Println(ff.ModTime().Unix())
	if err != nil {
		//Log.Error("file open err:%s\n", err)
		htmlwheel.ResultLog = fmt.Sprintf("file open err:%s\n", err)
		htmlwheel.Fail = true
		return
	}
	defer func() {
		_ = f.Close()
	}()

	//// xxx 临时页面变量
	//var dom *goquery.Document = DomFree.Get().(*goquery.Document)
	//var tmpSele1 *goquery.Selection = SeleFree.Get().(*goquery.Selection)
	//var tmpSele2 *goquery.Selection = SeleFree.Get().(*goquery.Selection)
	//var tmpSele3 *goquery.Selection = SeleFree.Get().(*goquery.Selection)
	//var FiledsHeader []string = FiledHeaderFree.Get().([]string)
	//var Summary_tbodyHeader = SeleFree.Get().(*goquery.Selection)

	htmlwheel.Dom, err = goquery.NewDocumentFromReader(f)

	if err != nil {
		//log.Fatal(err)
		//Log.Error("dom err:%s\n", err)
		htmlwheel.ResultLog = fmt.Sprintf("dom err:%s", err)
		htmlwheel.Fail = true
		return
	}

	// check chinese
	tmpSele1 := htmlwheel.Dom.Find("body:contains(备份作业)") // 正常

	tmpSele2 := htmlwheel.Dom.Find("body:contains(Report)")
	tmpSele3 := htmlwheel.Dom.Find("body:contains('Group By')")

	defer func() {
		//htmlwheel.Dom.Empty()
		//tmpSele1.Empty()
		//tmpSele2.Empty()
		//tmpSele3.Empty()
		//dom = nil
		tmpSele3 = nil
		tmpSele2 = nil
		tmpSele1 = nil

	}()
	htmlwheel.FiledsHeader = make([]string, 0)
	Log.Info(fmt.Sprintf("使用HtmlFusion对象:%p,Dom:%p,tmpSel1:%p,tmpSel2:%p,tmpSel3:%p,Summary_tbodyHeader:%p,FiledsHeader:%p\n",
		htmlwheel, htmlwheel.Dom, tmpSele1, tmpSele2, tmpSele3, Summary_tbodyHeader, htmlwheel.FiledsHeader))

	if tmpSele1.Length() < 1 || tmpSele2.Length() >= 1 && tmpSele3.Length() >= 1 {
		htmlwheel.ChineseFile = false
		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "col not Chinese")
		htmlwheel.Fail = true
		return
	}

	tmpSele1 = htmlwheel.Dom.Find("body:contains(存储策略拷贝中的作业报告)") // 不正常
	tmpSele2 = htmlwheel.Dom.Find("body:contains(未找到符合选定条件的数据)") // 不正常

	if tmpSele1.Length() >= 1 || tmpSele2.Length() >= 1 {
		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "不正确的中文格式")
		htmlwheel.Fail = true
		return
	}

	tmpSele1 = htmlwheel.Dom.Find("body:contains(CommCell)")

	//tmpSele1.EachWithBreak(func(i int, selection *goquery.Selection) bool{
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

	htmlwheel.CommCell = htmlwheel.Dom.Text()
	//fmt.Printf("CommCell :%+v\n", ss)
	//CommCell = strings.TrimSpace(ss)
	htmlwheel.CommCell = strings.Split(htmlwheel.CommCell, "CommCell:")[1]
	htmlwheel.CommCell = strings.Split(htmlwheel.CommCell, "--")[0]
	htmlwheel.CommCell = strings.Replace(htmlwheel.CommCell, " ", "", -1)

	htmlwheel.DetailTbodyDataNum = 13
	htmlwheel.SummaryTbodyDataNum = 10

	detailDomFindStr := "body > table:nth-child(" + strconv.Itoa(htmlwheel.DetailTbodyDataNum) + ") > tbody > tr"
	tmpSele1 = htmlwheel.Dom.Find(detailDomFindStr + ":nth-child(1) > td")

	//tmpSele1 = htmlwheel.Dom.Find("body > table:nth-child(13) > tbody > tr:nth-child(1) > td")

	if tmpSele1.Length() == 0 {
		htmlwheel.Cv = g.Cv101 // 版本是10,HTML结构不同,列头一样
		// 特殊的CV10
		//t_tbodyHeader = htmlwheel.Dom.Find("body > table > tbody:nth-child(1) > tr:nth-child(1) > td[bgcolor='#cccccc']"  )
		//tmpSele1 = htmlwheel.Dom.Find("body > table:nth-child(19) > tbody > tr:nth-child(1) > td")
		htmlwheel.DetailTbodyDataNum = 19
		detailDomFindStr := "body > table:nth-child(" + strconv.Itoa(htmlwheel.DetailTbodyDataNum) + ") > tbody > tr"
		tmpSele1 = htmlwheel.Dom.Find(detailDomFindStr + ":nth-child(1) > td")

		htmlwheel.SummaryTbodyDataNum = 16
	}

	//cv8 := false
	//有两种格式的HTML ,代表是CV8
	if len(tmpSele1.Nodes) == 0 {
		htmlwheel.Cv = g.Cv8
		htmlwheel.DetailTbodyDataNum = 12
		htmlwheel.SummaryTbodyDataNum = 9
		detailDomFindStr := "body > table:nth-child(" + strconv.Itoa(htmlwheel.DetailTbodyDataNum) + ") > tbody > tr"
		tmpSele1 = htmlwheel.Dom.Find(detailDomFindStr + ":nth-child(1) > td")

	}

	// xxx 如果前面都没有列头，goquery前22个元素 找到列头的位置
	if tmpSele1.Length() == 0 {
		for i := 8; i <= 22; i++ { // 摘要表格第一列 不是列头，所以不会找到摘要表格
			detailDomFindStr := "body > table:nth-child(" + strconv.Itoa(i) + ") > tbody > tr"
			tmpSele1 = htmlwheel.Dom.Find(detailDomFindStr + ":nth-child(1) > td")

			//tmpSele1 = htmlwheel.Dom.Find("body > table:nth-child(" + strconv.Itoa(i) + ") > tbody > tr:nth-child(1) > td")
			//fmt.Printf("%d,%d\n", i, len(t_tbodyHeader.Nodes))\
			if len(tmpSele1.Nodes) > 9 { // xxx 找到一个至少有9列的行
				htmlwheel.DetailTbodyDataNum = i
				htmlwheel.Cv = g.Cv100
				break
			}
		}
	}

	SummaryDomFindStr := "body > table:nth-child(" + strconv.Itoa(htmlwheel.SummaryTbodyDataNum) + ") > tbody > tr"
	Summary_tbodyHeader = htmlwheel.Dom.Find(SummaryDomFindStr + ":nth-child(2) > td")

	// 摘要数据 至少10列
	if Summary_tbodyHeader.Length() == 0 {
		for i := 5; i < htmlwheel.DetailTbodyDataNum; i++ { // 摘要表格 肯定在 详细数据前面
			SummaryDomFindStr := "body > table:nth-child(" + strconv.Itoa(i) + ") > tbody > tr"
			Summary_tbodyHeader = htmlwheel.Dom.Find(SummaryDomFindStr + ":nth-child(2) > td")

			//tmpSele1 = htmlwheel.Dom.Find("body > table:nth-child(" + strconv.Itoa(i) + ") > tbody > tr:nth-child(1) > td")
			//fmt.Printf("%d,%d\n", i, len(t_tbodyHeader.Nodes))\
			if Summary_tbodyHeader.Length() > 10 { // xxx 找到一个至少有10列的行
				htmlwheel.SummaryTbodyDataNum = i
				htmlwheel.Cv = g.Cv100
				break
			}
		}
	}

	if len(tmpSele1.Nodes) == 0 || Summary_tbodyHeader.Length() == 0 {
		htmlwheel.Fail = true
		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "前22个元素，没有找到有效摘要数据或详细数据列头的表格")
		return

	}

	// 列头，序号为KEY
	// 循环找出列头
	var FindFiledHeader = func(i int, tmpSele2 *goquery.Selection) {
		defer func() {
			tmpSele2 = nil
		}()
		sa := tmpSele2.Text()
		htmlwheel.formatFields(i, &sa) // todo 格式化 列头
		//*FiledsHeader = append(*FiledsHeader, sa)
		//fmt.Println(i)
		htmlwheel.FiledsHeader = append(htmlwheel.FiledsHeader, sa)
		//FiledsHeader[i] = sa
	}

	tmpSele1.Each(FindFiledHeader)
	//fmt.Println(FiledsHeader)
	// 如果是CV8版本，手动修改列名, 删除 size of Backup cols 这列

	switch htmlwheel.Cv {
	case 8:
		//tmp := make([]string, 16)
		//copy(tmp, FiledsHeader)
		index := 9 // xxx cv8需要删除第9列
		htmlwheel.FiledsHeader = append(htmlwheel.FiledsHeader[0:index], htmlwheel.FiledsHeader[index+1:]...)
		//FiledsHeader = tmp[0:9]
		//fmt.Println(*FiledsHeader)
		//*FiledsHeader = append(*FiledsHeader, []string{"Data Transferred", "Data Written"}...)
		//fmt.Println(*FiledsHeader)
		//FiledsHeader = append(*FiledsHeader, tmp[10:]...)
		//fmt.Println(*FiledsHeader)
		htmlwheel.FiledsHeader[3] = "Job ID (CommCell)(Status)"

		break
	case 101:
		break
	default:
		break
	}

	switch true {
	case len(htmlwheel.FiledsHeader) < g.FieldsLen-6: // col not enough,运维报告中不用判断

		htmlwheel.Fail = true

		htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "col not enough")

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
		tmpSele2 = htmlwheel.Dom.Find("body:contains(备份作业)") // 正常
	} else {
		tmpSele2 = htmlwheel.Dom.Find("body:contains(on)")
	}

	var SearchGenTimeCommCell = func(i int, tmp *goquery.Selection) bool {
		defer func() {
			tmp = nil
		}()
		if i == 0 {
			htmlwheel.GenTime = tmpSele2.Text()
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

	tmpSele2.EachWithBreak(SearchGenTimeCommCell)

	utils.FormatTime(&htmlwheel.GenTime, true)
	htmlwheel.GenTime = fmt.Sprintf("%s000", htmlwheel.GenTime)

	// 与 定义好的 9列进行比对， 这9列必须存在
	//fmt.Println(htmlwheel)
	if !htmlwheel.ChineseFile { // 中文文件不比对
		for i := 0; i < len(g.NeedFields); i++ {
			//fmt.Println((*FiledsHeader)[i])
			if g.NeedFields[i] != (htmlwheel.FiledsHeader)[i] {
				if !strings.Contains((htmlwheel.FiledsHeader)[i], g.NeedFields[i]) { // 包含关键字即可
					//fmt.Printf("%+v\n", NeedFields[i])
					htmlwheel.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", htmlwheel.HtmlFileAbs, "NeedField not enough")
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
	//fmt.Println(htmlwheel.Cv, htmlwheel.SummaryTbodyDataNum, htmlwheel.DetailTbodyDataNum)

	switch htmlwheel.Cv {
	case 8:
		//t_tbodyData = htmlwheel.Dom.Find("body > table:nth-child("+strconv.Itoa(htmlwheel.TbodyDataNum)+") > tbody > tr")
		DetailFields = g.DetailFieldsMapCv8

		break
	case 101:
		//t_tbodyData = htmlwheel.Dom.Find("body > table:nth-child("+strconv.Itoa(htmlwheel.TbodyDataNum)+") > tbody > tr")
		DetailFields = g.DetailFieldsMap

		break
	case 11:
		DetailFields = g.DetailFieldsMap
		if Summary_tbodyHeader.Length() == 18 {
			htmlwheel.Cv = 10
		}

		break
	case 100: // xxx 自定义列头 ，先使用 默认列头
		DetailFields = g.DetailFieldsMap

		break
	default:
		//t_tbodyData = htmlwheel.Dom.Find("body > table:nth-child("+strconv.Itoa(htmlwheel.TbodyDataNum)+") > tbody > tr ")
		DetailFields = g.DetailFieldsMap

	}

	// 特殊的文件
	if Summary_tbodyHeader.Length() < 18 && htmlwheel.Cv != 8 {
		htmlwheel.Cv = g.Cv100
	}

	tmpSele1 = htmlwheel.Dom.Find("body > table:nth-child(" + strconv.Itoa(htmlwheel.DetailTbodyDataNum) + ") > tbody > tr ")

	htmlwheel.DetailTotalRow = tmpSele1.Length() - 1

	htmlwheel.FusionSqlArrs = make([]*HtmlRowData, tmpSele1.Length()-1)

	tmpSele1.Each(func(rowIndex int, rowsele *goquery.Selection) {
		if rowIndex == 0 { // 0行是列名
			return
		}

		cols := rowsele.Find("td")
		//fmt.Println(rowIndex,len(cols.Nodes))
		if htmlwheel.Cv != 8 {
			if (len(DetailFields) != len(cols.Nodes)) && len(cols.Nodes) != 9 { // xxx 9列数据可能是 包含错误信息的行
				//_ = Log.Error("htmfile:%s ,no eq rows:%d\n", htmlfile, len(DetailFields))
				htmlwheel.DetailFailRow++
				return
			}

		} else {
			if len(DetailFields) != len(cols.Nodes)-1 && len(cols.Nodes) != 9 { //
				//Log.Error("htmfile:%s ,no eq rows:%d\n", htmlfile, i)
				htmlwheel.DetailFailRow++
				return
			}

		}

		RowWheel := HtmlRowDataFree.Get().(*HtmlRowData)
		RowWheel.FormatErr = false // 默认格式正常
		var DataMap = make(map[string]interface{}, 28)

		// xxx 判断是否为 9列 的数据  ,示例 HFSERVERLC3_32DAY_4COPY_YUEBAO_Chinese.html

		if len(cols.Nodes) == 9 {
			//fmt.Println(cols.Eq(3).Text(),strings.Contains(cols.Eq(3).Text(), "不适用")) // jobid 包含关键字

			//fmt.Println(strings.Contains(cols.Eq(6).Text(),"("),strings.Contains(cols.Eq(7).Text(),"(")) // 两个时间列不包含括号
			//fmt.Println(cols.Eq(8).Text()) // 错误信息不为空
			if strings.Contains(cols.Eq(3).Text(), "不适用") ||
				(!strings.Contains(cols.Eq(6).Text(), "(") || !strings.Contains(cols.Eq(7).Text(), "(")) &&
					len(cols.Eq(8).Text()) > 4 {
				RowWheel.FormatErr = true

			} else {
				htmlwheel.DetailFailRow++
				return
			}
		}

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

		if RowWheel.FormatErr {
			RowWheel.RunType = "失败"
			RowWheel.SolveType = "0" // 未解决
			cols.Each(func(colnum int, colsele *goquery.Selection) {
				switch true {
				case colnum >= 9:
					return
				case colnum == 8:
					ss1 := colsele.Text()
					RowWheel.RowReasonforfailure = ss1
				case colnum < 8:
					ss1 := colsele.Text()
					ss1 = *RowWheel.FormatErr_FormatValues(colnum, &ss1, htmlwheel.ChineseFile) // todo 格式化 数据,重复使ssl
					DataMap[DetailFields[colnum]] = ss1
					//colsele.Empty()
				}
				//if colnum >= 9 { // 第9列不用
				//	return
				//}
				//fmt.Println(colnum,DetailFields[colnum],colsele.Text())
				//ss1 := colsele.Text()
				//fmt.Println(ss1)
				//ss1 = *RowWheel.FormatErr_FormatValues(colnum, &ss1, htmlwheel) // todo 格式化 数据,重复使ssl
				//DataMap[DetailFields[colnum]] = ss1
				//colsele.Empty()
			})

		} else {
			cols.Each(func(colnum int, colsele *goquery.Selection) {
				if htmlwheel.Cv == 8 && colnum == 9 {
					return
				}
				ss1 := colsele.Text()

				ss1 = *RowWheel.FormatValues(colnum, &ss1, htmlwheel.ChineseFile) // todo 格式化 数据,重复使ssl
				if len(rowsele.Next().Find(" td").Nodes) == 1 {
					//fmt.Println(i, rowsele.Next().Find(" td").Text())
					RowWheel.RowReasonforfailure = rowsele.Next().Find(" td").Text()

					//fmt.Println(rowReasonforfailure)

					if strings.Contains(RowWheel.RowReasonforfailure, "Fail") ||
						strings.Contains(RowWheel.RowReasonforfailure, "ER") || len(strings.TrimSpace(RowWheel.RowReasonforfailure)) > 5 {
						RowWheel.Color = "2"
					}

				} else {
					RowWheel.RowReasonforfailure = "NULL"
				}

				//
				//fmt.Println(strings.Split())
				DataMap[DetailFields[colnum]] = ss1
				//tmpSubData = append(tmpSubData, ss1)
				//fmt.Println(i,ss1)
				//colsele.Empty()
			})
		}

		DataMap[g.DetailFieldsMapPlus[17]] = htmlwheel.CommCell
		DataMap[g.DetailFieldsMapPlus[18]] = RowWheel.StartTime
		DataMap[g.DetailFieldsMapPlus[19]] = RowWheel.Subclient
		DataMap[g.DetailFieldsMapPlus[20]] = RowWheel.ApplicationSize
		DataMap[g.DetailFieldsMapPlus[21]] = filepath.Base(htmlwheel.HtmlFileAbs)
		DataMap[g.DetailFieldsMapPlus[22]] = strconv.FormatInt(time.Now().Unix(), 10) + "000"
		DataMap[g.DetailFieldsMapPlus[28]] = RowWheel.EndTime
		//fmt.Printf("%v,%v\n", &rowReasonforfailure, rowReasonforfailure)

		// xxx 截取最大字符长度
		utils.Rebuildstr(&RowWheel.RowReasonforfailure)
		utils.Split(&RowWheel.RowReasonforfailure, 4000)

		//fmt.Printf("%v,%v\n", &rowReasonforfailure, rowReasonforfailure)
		//fmt.Println(Color)
		//fmt.Println(rowReasonforfailure)
		DataMap[g.DetailFieldsMapPlus[23]] = strings.TrimSpace(RowWheel.RowReasonforfailure)
		DataMap[g.DetailFieldsMapPlus[24]] = RowWheel.Color

		if strings.Contains(RowWheel.RowReasonforfailure, "因同一子客户端有另一作业在运行而未运行") {
			DataMap[g.DetailFieldsMapPlus[25]] = "未运行"
			DataMap[g.DetailFieldsMapPlus[27]] = "0"
		} else {
			DataMap[g.DetailFieldsMapPlus[25]] = RowWheel.RunType
			DataMap[g.DetailFieldsMapPlus[27]] = RowWheel.SolveType
		}

		DataMap[g.DetailFieldsMapPlus[26]] = htmlwheel.GenTime

		DataMap[g.DetailFieldsMapPlus[29]] = strconv.FormatInt(fs.ModTime().Unix(), 10) + "000"

		//DataMap = DataMap // 重复利用DataMap 每行必须要循环生成 SQL
		//htmlwheel.FusionSqlArrs = append(htmlwheel.FusionSqlArrs, RowWheel)
		htmlwheel.FusionSqlArrs[rowIndex-1] = RowWheel

		//if htmlwheel.Cv != 8 && DataMap[DetailFields[9]] != nil && !htmlwheel.Fail  {
		//	//datatransfer := DataMap[DetailFields[9]].(string)
		//	//utils.FormatData(&datatransfer)
		//	//DataMap[g.DetailFieldsMapPlus[29]] = datatransfer
		//}

		RowWheel.DataMap = DataMap

		// xxx 判断 是否 夯住

		// 运行中的需要检查  ， 同一个jobid 按照插入时间降序 第一个的时间差 与 本次是否差别 在 1小时以上，appcationsize是否没有变化
		if RowWheel.RunType == "活动" {
			RowWheel.CheckHang()
		}

		RowWheel.ToUpdateSql(false)
		//fmt.Println(RowWheel.UpdateSql.String())
		//RowWheel.ToInsertSql(false) // 后生成
		//fmt.Println(RowWheel.InsertSql.String())
		// xxx 清空MAP
		//utils.ClearMap(DataMap)
		//rowsele=nil

	})

	return
}
