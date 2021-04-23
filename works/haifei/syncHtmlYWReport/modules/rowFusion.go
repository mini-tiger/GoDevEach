package modules

import (
	"fmt"
	"haifei/syncHtmlYWReport/g"
	"haifei/syncHtmlYWReport/utils"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var HtmlRowDataFree = sync.Pool{
	New: func() interface{} {
		return &HtmlRowData{}
	},
}

type HtmlRowData struct {
	StartTime           string
	EndTime             string
	ApplicationSize     string
	SubClient           string
	RowReasonforfailure string
	Subclient           string
	Color               string
	RunType             string
	SolveType           string
	FormatErr           bool // 是否是 只有9列，但也需要导入的数据
	DataMap             map[string]interface{}
	InsertSql           strings.Builder
	UpdateSql           strings.Builder
	//Sqls				map[string]string
}

var fl float64
var err error
var out [][]string

func (Row *HtmlRowData) FormatErr_FormatValues(index int, s *string, ChineseFile bool) *string {
	//fmt.Println(index, s)
	utils.Rebuildstr(s)

	switch true {

	case index == 2: // subclient

		switch true {
		case strings.Contains(*s, "/"):

			//xxx 取最后一个 /后面 的所有字符
			*s = (*s)[strings.LastIndex(*s, "/")+1:]
			//
			////fallthrough
			if strings.Contains(*s, "Logcommand line") {
				Row.Subclient = "Logcommand line"
				break
			}

			if strings.Contains(*s, "command line") ||
				(strings.Contains(*s, "command") && strings.Contains(*s, "line")) {
				Row.Subclient = "command line"
				break
			}

			Row.Subclient = *s

			break
		default:
			Row.Subclient = *s

		}

		break

	case index == 3: // jobid
		*s = strings.ReplaceAll(strings.ReplaceAll(*s, "\n", ""), " ", "")
		if strings.Contains(*s, "不适用") {
			*s = strings.Split(*s, "不适用")[0]
		}

		return s

	case index == 6: // starttime 赋值 全局变量，对应自定义列

		if strings.Contains(*s, "2") && strings.Contains(*s, ":") {
			if strings.Contains(*s, "(") {
				Row.StartTime = strings.Split(*s, "(")[0]
			} else {
				Row.StartTime = *s
			}
			Row.StartTime = strings.TrimSpace(Row.StartTime)
			utils.FormatTime(&Row.StartTime, ChineseFile)
			Row.StartTime = fmt.Sprintf("%s000", Row.StartTime)
		}

		break
	case index == 7:
		if strings.Contains(*s, "2") && strings.Contains(*s, ":") {
			if strings.Contains(*s, "(") {
				Row.EndTime = strings.Split(*s, "(")[0]
			} else {
				Row.EndTime = *s
			}
			Row.EndTime = strings.TrimSpace(Row.EndTime)
			utils.FormatTime(&Row.EndTime, ChineseFile)
			Row.EndTime = fmt.Sprintf("%s000", Row.EndTime)
		}

		break
	case index == 8:

		Row.RowReasonforfailure = *s

		break
	}
	return s
}

func (Row *HtmlRowData) FormatValues(index int, s *string, ChineseFile bool) *string {
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
				Row.Subclient = "Logcommand line"
				break
			}

			if strings.Contains(*s, "command line") ||
				(strings.Contains(*s, "command") && strings.Contains(*s, "line")) {
				Row.Subclient = "command line"
				break
			}

			Row.Subclient = *s

			break
		default:
			Row.Subclient = *s

		}

		break

	case index == 3: // jobid
		*s = strings.ReplaceAll(strings.ReplaceAll(*s, "\n", ""), " ", "")
		out = regexp.MustCompile(`\(([a-zA-Z)]+)\)`).FindAllStringSubmatch(*s, -1)

		if len(out) != 1 {
			Row.RunType = "NULL"
		}

		if len(out[0]) != 2 {
			Row.RunType = "NULL"
		}

		Row.SolveType = "2" // 默认完成
		switch out[0][1] {
		case "A":
			Row.RunType = "活动"
		case "C":
			Row.RunType = "已完成"
		case "CMTD":
			Row.RunType = "提交"
		case "D":
			Row.RunType = "已延迟"
			Row.SolveType= "0"
		case "CWE":
			Row.RunType = "完成但有错误"
			Row.SolveType = "0" // 未解决
		case "CWW":
			Row.RunType = "完成但有警告"
		case "K":
			Row.RunType = "已终止"
			Row.SolveType = "0"
		case "F":
			Row.RunType = "失败"
			Row.SolveType = "0"
		case "STK":
			Row.RunType = "卡住"
			Row.SolveType = "0"
		case "NR":
			Row.RunType = "未运行"
			Row.SolveType = "0"
		default:
			Row.RunType = "NULL"
			Row.SolveType = "0"
		}

		return s

	case index == 6: // starttime 赋值 全局变量，对应自定义列
		//StartTimeValue := ""

		// 1.不包括 2 和 : 格式不对不转换

		if strings.Contains(*s, "2") && strings.Contains(*s, ":") {
			if strings.Contains(*s, "(") {
				Row.StartTime = strings.Split(*s, "(")[0]
			} else {
				Row.StartTime = *s
			}
			Row.StartTime = strings.TrimSpace(Row.StartTime)
			utils.FormatTime(&Row.StartTime, ChineseFile)
			Row.StartTime = fmt.Sprintf("%s000", Row.StartTime)
		}

		break
	case index == 7:
		if strings.Contains(*s, "2") && strings.Contains(*s, ":") {
			if strings.Contains(*s, "(") {
				Row.EndTime = strings.Split(*s, "(")[0]
			} else {
				Row.EndTime = *s
			}
			Row.EndTime = strings.TrimSpace(Row.EndTime)
			utils.FormatTime(&Row.EndTime, ChineseFile)
			if strings.Contains(Row.EndTime, "-") {
				Row.EndTime = "NULL"
			} else {
				Row.EndTime = fmt.Sprintf("%s000", Row.EndTime)
			}

		}

		break
	case index == 8:
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
		case strings.Contains(Row.ApplicationSize, "字节"):
			Row.ApplicationSize = strings.Split(Row.ApplicationSize, " 字节")[0]
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

func (Row *HtmlRowData) CheckHang() {
	basesql := fmt.Sprintf(g.BaseCheckHangSql, Row.DataMap["APPLICATIONSIZE"], g.DetailTableName,
		Row.DataMap["Job ID (CommCell)(Status)"], Row.DataMap["FILESYNCTIME"])
	//fmt.Println(basesql)
	var ID, diff string
	DB.QueryRow(basesql).Scan(&ID, &diff)
	if diff == "0" {
		Row.DataMap["HANG"] = "1"
		Row.DataMap["SOLVETYPE"] = "0"
		Row.DataMap["RELID"] = ID
	}
	//fmt.Println(basesql, diff)
}

func (Row *HtmlRowData) ToUpdateSql(IsSummary bool) {

	Row.UpdateSql = strings.Builder{}

	//if len(h.SortKeys) <= 0 || (len(Row.DataMap) != len(h.SortKeys)) { // 可能有9列的情况，要判断是否与前一次的KEYS 一样长
	//	var keys []string = make([]string, 0)
	//	//
	//	for k := range Row.DataMap {
	//		keys = append(keys, k)
	//	}
	//	sort.Strings(keys)
	//
	//	h.SortKeys = keys
	//}

	//fmt.Println(len(h.SortKeys))
	//fmt.Println(h.SortKeys)

	if IsSummary {
		//updateSql := "update " + g.SummaryTableName + " set "
		//*Row.UpdateSql = updateSql + updatesliceToString(DataMap)
		Row.UpdateSql.WriteString("update ")
		Row.UpdateSql.WriteString(g.SummaryTableName)
		Row.UpdateSql.WriteString(" set ")
		Row.UpdateSql.WriteString(updatesliceToString(Row.DataMap))
		return
	}
	//var returnstring string

	//if len(Row.SortKeys) < 8{ // 如果存在则不重新生成

	//}
	//RowSortKey = RowSortKey[0:0]
	//for k := range DataMap {
	//	RowSortKey = append(RowSortKey, k)
	//}
	//sort.Strings(RowSortKey)

	//var colsStr string
	//fmt.Printf("memAdd:%v\n",&colsStr)
	var tmpstr strings.Builder

	defer tmpstr.Reset()

	//Row.UpdateSql.WriteString( "update " + g.DetailTableName + " set ")

	for k, v := range Row.DataMap {

		// fixme 插入源文件与插入时间跳过
		if v == "NULL" || v == "" || v == "INSERTTIME" || v == "HTMLFILE" || v == nil || k == "HANG" { // HANG状态不更新,重复文件会 判断有误
			continue
		}

		//*Row.UpdateSql = (*Row.UpdateSql) + "\"" + v + "\"='" + ((DataMap)[v]).(string) + "',"
		//wherestring = wherestring + "\"" + v + "\"='" + (DataMap)[v] + "' and "
		tmpstr.WriteString("\"" + k + "\"='" + (v).(string) + "',")

	}
	//fmt.Printf("memAdd:%v\n",&colsStr)

	// fixme  wheresql只保留 COMMCELL  Jobid, 替换上面的 where
	whereSql := fmt.Sprintf(" 1=1 and \"COMMCELL\"='%s' and \"Job ID (CommCell)(Status)\" = '%s' and" +
		"( \"REPORTTIME\" = '%s' OR \"HTMLFILE\"='%s')",
		(Row.DataMap)["COMMCELL"], (Row.DataMap)["Job ID (CommCell)(Status)"], (Row.DataMap)["REPORTTIME"],(Row.DataMap)["HTMLFILE"])

	//for _, v := range h.SortKeys {
	//	if v == "COMMCELL" || v == "Job ID (CommCell)(Status)" {
	//		wherestring = wherestring + " and \"" + v + "\"='" + (DataMap)[v] + "'"
	//	}
	//}
	// xxx 有可能 最后一个逗号没有去除，可能最后一个是空被跳过

	if tmpstr.String()[tmpstr.Len()-1:] == "," {
		//*Row.UpdateSql = "update " + g.DetailTableName + " set " + strings.TrimRight(tmpstr.String(), ",") + " where " + whereSql
		Row.UpdateSql.WriteString(fmt.Sprintf("update %s set %s where %s", g.DetailTableName, strings.TrimRight(tmpstr.String(), ","), whereSql))
	} else {

		//*Row.UpdateSql = "update " + g.DetailTableName + " set " + tmpstr.String() + " where " + whereSql
		Row.UpdateSql.WriteString(fmt.Sprintf("update %s set %s where %s", g.DetailTableName, tmpstr.String(), whereSql))
	}
	//fmt.Println(*Row.UpdateSql)
	//strings.TrimRight(Row.UpdateSql.String(), ",")

	//Row.UpdateSql.WriteString(" where ")
	//Row.UpdateSql.WriteString(whereSql)
	//returnstring = colsStr + " where " + wherestring
	//return &returnstring
}

func (Row *HtmlRowData) ToInsertSql(IsSummary bool) () {
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

	Row.InsertSql = strings.Builder{}
	//Row.InsertSql = ""
	//var i uint32 = 0
	//fmt.Println("ddddddddd", Row.DataMap)

	for k, v := range Row.DataMap {

		if v == "NULL" || v == "" || v == nil {

			continue
		}

		colsStr.WriteString("\"")
		colsStr.WriteString(k)
		colsStr.WriteString("\",")
		//*Row.InsertSql = (*Row.InsertSql) + "'" + ((DataMap)[v]).(string) + "'" + ","
		tmpstr.WriteString("'")
		tmpstr.WriteString(v.(string))
		tmpstr.WriteString("',")

	}
	//fmt.Println(len(keys),len(sl)
	// xxx 有可能 最后一个逗号没有去除，可能最后一个是空值被跳过
	tmpTableName := g.DetailTableName
	if IsSummary {
		tmpTableName = g.SummaryTableName
	}

	if tmpstr.String()[tmpstr.Len()-1:] == "," {
		//*Row.InsertSql = "insert into " + tmpTableName + " ( " + strings.TrimRight(colsStr.String(), ",") + ")" +
		//	" values(" + strings.TrimRight(tmpstr.String(), ",") + ")"
		Row.InsertSql.WriteString(fmt.Sprintf("insert into %s ( %s )  values(%s)",
			tmpTableName, strings.TrimRight(colsStr.String(), ","), strings.TrimRight(tmpstr.String(), ",")))
	} else {
		//*Row.InsertSql = "insert into " + tmpTableName + " ( " + colsStr.String() + " values(" + tmpstr.String()
		Row.InsertSql.WriteString(fmt.Sprintf("insert into %s ( %s  values(%s", tmpTableName, colsStr.String(), tmpstr.String()))
	}

	//return returnstring
}
