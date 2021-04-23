package modules

import (
	"database/sql"
	"fmt"
	"gitee.com/taojun319/tjtools/file"
	"github.com/ccpaging/nxlog4go"
	"haifei/syncHtml/g"
	"haifei/syncHtml/utils"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
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

type HtmlRowData struct {
	StartTime           string
	ApplicationSize     string
	SubClient           string
	RowReasonForFailure string
	Color               string
	RunType             string
	//DataMap             map[string]interface{}
	InsertSql *string
	UpdateSql *string
	//Sqls				map[string]string
}

type HtmlFusion struct {
	HtmlFileAbs    string         // 文件名绝对路径
	NewHtmlFileAbs string         //备份文件名绝对路径
	Cv             int            //文件版本
	ResultLog      string         //错误提示字符串
	Fail           bool           // 是否成功，流程是否往下进行
	FusionSqlArrs  []*HtmlRowData //字段和数据的MAP
	//SqlArrs        []map[string]string//SQL语句MAP，分插入和更新
	//SqlArrs        []map[string]string //SQL语句MAP，分插入和更新
	SortKeys    []string
	ChineseFile bool
	//DetailSqlArr  [](*map[string]string)
	//FiledsHeader []string
	InsertSuccess int // 插入和更新 成功多少条
	UpdateSuccess int
	InsertFail    int
	UpdateFail    int
	FailRow       int // 行数据不正常的个数
	TotalRow      int
	CommCell      string
	GenTime       string
	TbodyDataNum  int // 数据在goquery 节点的位置
	DB            *sql.DB
}

var HtmlFusionFree = sync.Pool{
	New: func() interface{} {
		return &HtmlFusion{}
	},
}

var HtmlRowDataFree = sync.Pool{
	New: func() interface{} {
		return &HtmlRowData{}
	},
}

var Log *nxlog4go.Logger
var c *g.Config

//var RowSortKey = make([]string, 25)
var logstr string

func LoadLogAndCfg() {
	Log = g.GetLog()
	c = g.GetConfig()
}

func FormatFields(index int, s *string, cv int) {
	//rstr := ""
	switch true {
	case index == 0 && strings.Contains(*s, "Client"):
		*s = "DATACLIENT"
		break
	case strings.Contains(*s, "Phase(Write End Time)"):

		*s = strings.Replace(*s, "(Write End Time)", "", -1)

		break
	case index == 7 && cv == 101:

		*s = "End Time or Current Phase" // cv 101
		break

	case index == 10 && cv == 101:
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

func (Row *HtmlRowData) FormatValues(index int, s *string, h *HtmlFusion) *string {
	//fmt.Println(index, s)
	utils.Rebuildstr(s)
	//if strings.Contains(*s, "N/A") {
	//	rstr = strings.Replace(*s, "N/A", "", -1)
	//	s = &rstr
	//} else {
	//	rstr = *s
	//}
	//if strings.Contains(*s, "%") {
	//	rstr = strings.Replace(*s, "%", "%%", -1)
	//	s = &rstr
	//} else {
	//	rstr = *s
	//}

	switch true {
	//case strings.Contains(s, "N/A"):
	//	rstr = strings.Replace(s, "N/A", "", -1)
	//	s = rstr
	case index == 2: // subclient

		//tmpstr := s
		//fmt.Println(tmpstr)
		//fmt.Println(strings.Contains(tmpstr,"/"),strings.Split(s, "/")[1])
		switch true {
		case strings.Contains(*s, "/"):

			//if strings.Contains(*s, "N/A/") {
			//	*s = strings.ReplaceAll(*s, "N/A/", "")
			//} else {
			//	*s = strings.Split(*s, "/")[1]
			//}
			//
			//if *s == "" {
			//	var reg = regexp.MustCompile("/.*")
			//	match := reg.FindString(*s)
			//	*s = match
			//}
			//xxx 取最后一个 /后面 的所有字符
			*s = (*s)[strings.LastIndex(*s, "/")+1:]
			//
			////fallthrough
			if strings.Contains(*s, "Logcommand line") {
				Row.SubClient = "Logcommand line"
				break
			}

			if strings.Contains(*s, "command line") ||
				(strings.Contains(*s, "command") && strings.Contains(*s, "line")) {
				Row.SubClient = "command line"
				break
			}

			Row.SubClient = *s

			break
		default:
			Row.SubClient = *s

		}

		break

	case index == 3: // jobid
		*s = strings.ReplaceAll(strings.ReplaceAll(*s, "\n", ""), " ", "")
		rex := regexp.MustCompile(`\(([a-zA-Z)]+)\)`)
		out := rex.FindAllStringSubmatch(*s, -1)

		if len(out) != 1 {
			Row.RunType = "NULL"
		}

		if len(out[0]) != 2 {
			Row.RunType = "NULL"
		}

		switch out[0][1] {
		case "A":
			Row.RunType = "活动"
		case "C":
			Row.RunType = "已完成"
		case "CMTD":
			Row.RunType = "提交"
		case "D":
			Row.RunType = "已延迟"
		case "CWE":
			Row.RunType = "完成但有错误"
		case "CWW":
			Row.RunType = "完成但有警告"
		case "K":
			Row.RunType = "已终止"
		case "F":
			Row.RunType = "失败"
		case "STK":
			Row.RunType = "卡住"
		case "NR":
			Row.RunType = "未运行"
		default:
			Row.RunType = "NULL"
		}

		return s

	case index == 6: // starttime 赋值 全局变量，对应自定义列
		//StartTimeValue := ""
		if strings.Contains(*s, "(") {
			Row.StartTime = strings.Split(*s, "(")[0]
		} else {
			Row.StartTime = *s
		}

		Row.StartTime = strings.TrimSpace(Row.StartTime)
		//StartTimeValue = fmt.Sprint("to_date('" + StartTimeValue + "','mm/dd/yyyy hh24:mi:ss')")
		utils.FormatTime(&Row.StartTime, h.ChineseFile)
		//StartTimeArr = append(StartTimeArr, fmt.Sprintf("%s000", StartTimeValue))
		Row.StartTime = fmt.Sprintf("%s000", Row.StartTime)

		//StartTime = fmt.Sprintf("to date(%s mm/dd/yyyy hh24:mi:ss)",StartTime)
		//StartTime = "2019-08-12 04:00:00"
		break
	case index == 8:
		//as := ""
		var fl float64
		var err error
		if strings.Contains(*s, "(") {
			Row.ApplicationSize = strings.Split(*s, "(")[0]
		} else {
			Row.ApplicationSize = *s
		}

		switch true {
		case strings.Contains(Row.ApplicationSize, "TB"):
			Row.ApplicationSize = strings.Split(Row.ApplicationSize, " TB")[0]
			fl, err = strconv.ParseFloat(Row.ApplicationSize, 64)
			if err != nil {
				_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
				break
			}

			//am:=strconv.FormatFloat(fl*1024,'E',-1,64)
			Row.ApplicationSize = fmt.Sprintf("%.2f", fl*1024*1024)
		case strings.Contains(Row.ApplicationSize, "GB"):
			Row.ApplicationSize = strings.Split(Row.ApplicationSize, " GB")[0]
			fl, err = strconv.ParseFloat(Row.ApplicationSize, 64)
			if err != nil {
				_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
				break
			}

			//am:=strconv.FormatFloat(fl*1024,'E',-1,64)
			Row.ApplicationSize = fmt.Sprintf("%.2f", fl*1024)

		case strings.Contains(Row.ApplicationSize, "MB"):
			Row.ApplicationSize = strings.Split(Row.ApplicationSize, " MB")[0]
			fl, err = strconv.ParseFloat(Row.ApplicationSize, 64)
			if err != nil {
				_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
				break
			}
			Row.ApplicationSize = fmt.Sprintf("%.2f", fl)
		case strings.Contains(Row.ApplicationSize, "KB"):
			Row.ApplicationSize = strings.Split(Row.ApplicationSize, " KB")[0]
			fl, err = strconv.ParseFloat(Row.ApplicationSize, 64)
			if err != nil {
				_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
				break
			}
			Row.ApplicationSize = fmt.Sprintf("%.2f", fl/1024)
		case strings.Contains(Row.ApplicationSize, "Bytes"):
			Row.ApplicationSize = strings.Split(Row.ApplicationSize, " Bytes")[0]
			fl, err = strconv.ParseFloat(Row.ApplicationSize, 64)
			if err != nil {
				_ = Log.Error("应用程序大小转换float失败 err:%s\n", err)
				break
			}
			Row.ApplicationSize = fmt.Sprintf("%.2f", fl/1024/1024)
		case strings.Contains(Row.ApplicationSize, "Not Run because of another job running for the same subclient"):
			break
		default:

			floatNum, _ := regexp.MatchString(`(\d+)\.(\d+)`, Row.ApplicationSize)
			//fmt.Println(i,err)

			intNum, _ := regexp.MatchString(`(\d+)`, Row.ApplicationSize)
			if floatNum || intNum {
				break
			}

			application, e := strconv.Atoi(strings.TrimSpace(Row.ApplicationSize))
			if e != nil {
				_ = Log.Error("change type string:%s,err:%s\n", Row.ApplicationSize, e)
			}
			if application == 0 {
				Row.ApplicationSize = "0"
			}

		}

		break
	}
	return s
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
		//fmt.Println("dddddddd")
		NewName := utils.NewDateFileName(h.HtmlFileAbs, bakDir)
		//fmt.Println("dddddddddddd",NewName)
		err := os.Rename(h.HtmlFileAbs, NewName)
		if err != nil {
			_ = Log.Error("osRename html File Error DIR: %s, err:%s\n", bakDir, err)
		} else {
			Log.Info("move file old:%s,new:%s\n", h.HtmlFileAbs, NewName)
		}
	} else {
		_, fileName := filepath.Split(h.HtmlFileAbs)
		Log.Info("move file old:%s,new:%s\n", h.HtmlFileAbs, filepath.Join(bakDir, fileName))
		err = os.Rename(h.HtmlFileAbs, filepath.Join(bakDir, fileName))
		if err != nil {
			_ = Log.Error("osRename html File Error DIR: %s, err:%s\n", bakDir, err)
		}
	}

}

func (h *HtmlFusion) ResultLogMoveFile() {
	//file := h.HtmlFileAbs
	if h.Fail {
		//MoveFailFileChan <- h.HtmlFileAbs

		//Log.Error(h.ResultLog)
		_ = Log.Error("%s,HtmFile:%s, version:%d,totalRows:%d,FailRow:%d,InsertSuccess:%d,UpdateSuccess:%d,UpdateFail:%d,InsertFail:%d,len(FusionMapSum):%d,CommCell:%s,GenTime:%s\n",
			h.ResultLog,
			filepath.Base(h.HtmlFileAbs),
			h.Cv,
			h.TotalRow,
			h.FailRow,
			h.InsertSuccess,
			h.UpdateSuccess,
			h.UpdateFail,
			h.InsertFail,
			len(h.FusionSqlArrs),
			h.CommCell, h.GenTime)
		h.MoveFile(&c.HtmlFailDir)
	} else {
		//Log.Info("finish covert file:%s\n", h.HtmlFileAbs)
		//MoveBakFileChan <- h.HtmlFileAbs

		Log.Info("HtmFile:%s, version:%d,totalRows:%d,FailRow:%d,InsertSuccess:%d,UpdateSuccess:%d,UpdateFail:%d,InsertFail:%d,len(FusionMapSum):%d,CommCell:%s,GenTime:%s\n",
			filepath.Base(h.HtmlFileAbs),
			h.Cv,
			h.TotalRow,
			h.FailRow,
			h.InsertSuccess,
			h.UpdateSuccess,
			h.UpdateFail,
			h.InsertFail,
			len(h.FusionSqlArrs),
			h.CommCell, h.GenTime)

		if h.UpdateSuccess > 0 {
			_ = Log.Warn("!!!!!!!!!HtmlFile:%s, update gt 0\n", h.HtmlFileAbs)
		}
		h.MoveFile(&c.HtmlBakDir)

	}
}

func (Row *HtmlRowData) ToUpdateSql(h *HtmlFusion, DataMap map[string]interface{}) {

	//var returnstring string

	//if len(Row.SortKeys) < 8{ // 如果存在则不重新生成
	if len(h.SortKeys) <= 0 {
		var keys  = make([]string, 0, len(DataMap))
		//
		for k := range DataMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		h.SortKeys = keys
	}

	//}
	//RowSortKey = RowSortKey[0:0]
	//for k := range DataMap {
	//	RowSortKey = append(RowSortKey, k)
	//}
	//sort.Strings(RowSortKey)

	//var colsStr string
	//fmt.Printf("memAdd:%v\n",&colsStr)
	var tmpStr strings.Builder

	defer tmpStr.Reset()
	Row.UpdateSql = new(string)

	//Row.UpdateSql.WriteString( "update " + g.DetailTableName + " set ")

	for i, v := range h.SortKeys {

		// fixme 插入源文件与插入时间跳过
		if (DataMap)[v] == "NULL" || (DataMap)[v] == "" || v == "INSERTTIME" || v == "HTMLFILE" {
			continue
		}

		switch true {
		case i == len(DataMap)-1:
			//returnstring = returnstring
			//returnstring = returnstring + "\"" + sl1[i] + "\"" + ")"
			//*Row.UpdateSql = (*Row.UpdateSql) + "\"" + v + "\"='" + ((DataMap)[v]).(string) + "'"
			tmpStr.WriteString("\"" + v + "\"='" + ((DataMap)[v]).(string) + "'")
			//wherestring = wherestring + "\"" + v + "\"='" + (DataMap)[v] + "'"

			break
		case i < len(DataMap)-1:
			//*Row.UpdateSql = (*Row.UpdateSql) + "\"" + v + "\"='" + ((DataMap)[v]).(string) + "',"
			//wherestring = wherestring + "\"" + v + "\"='" + (DataMap)[v] + "' and "
			tmpStr.WriteString("\"" + v + "\"='" + ((DataMap)[v]).(string) + "',")

			break
		}

	}
	//fmt.Printf("memAdd:%v\n",&colsStr)

	// fixme  wheresql只保留 COMMCELL  Jobid, 替换上面的 where
	whereSql := fmt.Sprintf(" 1=1 and \"COMMCELL\"='%s' and \"Job ID (CommCell)(Status)\" = '%s' ",
		(DataMap)["COMMCELL"], (DataMap)["Job ID (CommCell)(Status)"])

	//for _, v := range h.SortKeys {
	//	if v == "COMMCELL" || v == "Job ID (CommCell)(Status)" {
	//		wherestring = wherestring + " and \"" + v + "\"='" + (DataMap)[v] + "'"
	//	}
	//}
	// xxx 有可能 最后一个逗号没有去除，可能最后一个是空被跳过
	if tmpStr.String()[tmpStr.Len()-1:] == "," {
		*Row.UpdateSql = "update " + g.DetailTableName + " set " + strings.TrimRight(tmpStr.String(), ",") + ") where " + whereSql
	} else {
		*Row.UpdateSql = "update " + g.DetailTableName + " set " + tmpStr.String() + " where " + whereSql
	}

	//strings.TrimRight(Row.UpdateSql.String(), ",")

	//Row.UpdateSql.WriteString(" where ")
	//Row.UpdateSql.WriteString(whereSql)
	//returnstring = colsStr + " where " + wherestring
	//return &returnstring
}

func (Row *HtmlRowData) ToInsertSql(h *HtmlFusion, DataMap map[string]interface{}) () {
	//sl1 := sl[0:len(sl)]
	//var returnstring string
	//var keys []string = make([]string, 0)
	//for k, _ := range DataMap {
	//	keys = append(keys, k)
	//
	//}
	//sort.Strings(keys)

	//valueStr := ""
	var tmpstr strings.Builder
	var colsStr strings.Builder
	defer func() {
		tmpstr.Reset()
		colsStr.Reset()

	}()

	Row.InsertSql = new(string)
	//Row.InsertSql = ""
	for i, v := range h.SortKeys {
		if (DataMap)[v] == "NULL" || (DataMap)[v] == "" {
			continue
		}
		switch true {
		case i == len(DataMap)-1:
			//returnstring = returnstring
			//returnstring = returnstring + "\"" + sl1[i] + "\"" + ")"
			//*colsStr = *colsStr + "\"" + v + "\")"
			colsStr.WriteString("\"")
			colsStr.WriteString(v)
			colsStr.WriteString("\")")
			//*Row.InsertSql = (*Row.InsertSql) + "'" + ((DataMap)[v]).(string) + "')"
			tmpstr.WriteString("'" + ((DataMap)[v]).(string) + "')")

			break
		case i < len(DataMap)-1:
			//returnstring = returnstring
			//*colsStr = *colsStr + "\"" + v + "\"" + ","
			colsStr.WriteString("\"")
			colsStr.WriteString(v)
			colsStr.WriteString("\",")
			//*Row.InsertSql = (*Row.InsertSql) + "'" + ((DataMap)[v]).(string) + "'" + ","
			tmpstr.WriteString("'" + ((DataMap)[v]).(string) + "'" + ",")
			break
		}

	}
	//fmt.Println(len(keys),len(sl))
	// xxx 有可能 最后一个逗号没有去除，可能最后一个是空值被跳过
	if tmpstr.String()[tmpstr.Len()-1:] == "," {
		*Row.InsertSql = "insert into " + g.DetailTableName + " ( " + strings.TrimRight(colsStr.String(), ",") + ")" +
			" values(" + strings.TrimRight(tmpstr.String(), ",") + ")"
	} else {
		*Row.InsertSql = "insert into " + g.DetailTableName + " ( " + colsStr.String() + " values(" + tmpstr.String()
	}

	//return returnstring
}
func (h *HtmlFusion) StmtSql() {
	if h.Fail {
		return
	}

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

		if RowStruct == nil {
			continue
		}

		//Log.Error("%s\n", *RowStruct.UpdateSql)
		Result, err := h.DB.Exec(*RowStruct.UpdateSql)
		//fmt.Println(Result.LastInsertId())
		//fmt.Println(Result.RowsAffected())
		//fmt.Println("===========", r, err)
		if err != nil {
			_ = Log.Error("HtmlFile:%s ,Sql Exec Err:%s,,,Updatesql:%s\n",
				h.HtmlFileAbs, strings.Replace(err.Error(), "\n", "", 1), *RowStruct.UpdateSql)
			h.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "insert or update have err")
			h.Fail = true
			continue
		}

		r, err := Result.RowsAffected()
		if err != nil {
			_ = Log.Error("HtmlFile:%s ,Sql Exec Err:%s\n", h.HtmlFileAbs, strings.Replace(err.Error(), "\n", "", 1))
			h.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "insert or update have err")
			h.Fail = true
			h.UpdateFail++
			continue
		}

		if r == 0 {
			_, err := h.DB.Exec(*RowStruct.InsertSql)
			if err != nil {
				_ = Log.Error("HtmlFile:%s ,Sql Exec Err:%s,Sql:%s\n", h.HtmlFileAbs,
					strings.Replace(err.Error(), "\n", "", 1),
					*RowStruct.InsertSql)
				h.ResultLog = fmt.Sprintf("htmlFile %s, ERR:%s\n", h.HtmlFileAbs, "insert or update have err")
				h.Fail = true
				h.InsertFail++
				continue
			}
			h.InsertSuccess++
		} else {
			//Log.Warn("%+v\n", (sqlArr)[i]["update"])
			h.UpdateSuccess++
		}
		// xxx clear RowStruct

		//RowStruct.UpdateSql = nil
		//RowStruct.InsertSql = nil
		utils.Clear(RowStruct)
		HtmlRowDataFree.Put(RowStruct)
	}



	//Log.Info("htmlfile : %s ,total sql:%d,success update sql:%d,insert sql:%d\n", filepath.Base(htmlfile), len(*sqlArr), updatenum, insertnum)
	return

}
