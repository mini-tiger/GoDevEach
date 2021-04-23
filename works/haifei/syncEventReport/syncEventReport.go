package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"path/filepath"
	"time"
	"tjtools/utils"
	"tjtools/logDiy"
	"godev/works/haifei/syncEventReport/g"
	"encoding/json"
	"log"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-oci8"
)

//https://github.com/360EntSecGroup-Skylar/excelize
// 中文 https://xuri.me/excelize/zh-hans/
// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#File.GetSheetName

const (
	//ReadFile    = "10.155.2.4_yuebao_English.html"
	//htmlfileReg = "C:\\work\\go-dev\\src\\godev\\works\\haifei\\*.html"
	//htmlBakDir  = "C:\\work\\go-dev\\src\\godev\\works\\haifei\\bak\\"
	//timeinter   = 10
	timeLayout = "2006-01-02 15:04:05"
	FieldsLen  = 21 // 列一共有几列,使用时 要 减一
)

var MoveFileChan chan string = make(chan string, 0)
var startRunTime = time.Now().Unix()
var (
	CommCell        string // 客户端名
	StartTime       string
	ApplicationSize string
	Subclient       string
	TimeArr         []string = make([]string, 0)
)

var Log *logDiy.Log1
var c Config

type Config struct {
	HtmlfileReg string `json:"htmlfileReg"`
	HtmlBakDir  string `json:"htmlBakDir"`
	Logfile     string `json:"logfile"`
	LogMaxDays  int    `json:"logMaxDays"`
	Timeinter   int    `json:"timeinter"`
	OracleDsn   string `json:"oracledsn"`
}

func readconfig() {
	cfgstr := g.ParseConfig("/home/go/src/godev/works/haifei/syncEventReport/syncEventReport.json")
	err := json.Unmarshal([]byte(cfgstr), &c)
	if err != nil {
		log.Fatalln("parse config file fail:", err)
	}
}

func main() {

	readconfig()

	fmt.Printf("%+v", c)

	logDiy.InitLog1(c.Logfile, c.LogMaxDays)
	Log = logDiy.Logger()

	go waitMovefile()
	for {
		files, _ := filepath.Glob(c.HtmlfileReg)
		//fmt.Println(files) // contains a list of all files in the current directory
		d, h, m, s := utils.GetTime(time.Now().Unix() - startRunTime)
		Log.Printf("本次读取到 %d个文件,%v\n", len(files), files)
		Log.Printf("开始运行时间 %s,已经运行%d天%d小时%d分钟%d秒", time.Unix(startRunTime, 0).Format(timeLayout), d, h, m, s)
		if len(files) >= 1 {
			waitHtmlFile(files)
		}
		time.Sleep(time.Duration(c.Timeinter) * time.Second)
	}
}

func MoveFile(file string) {
	err := os.MkdirAll(c.HtmlBakDir, 0777)
	if err != nil {
		Log.Error("mkdir htmlbak Error DIR: %s, err:%s", c.HtmlBakDir, err)
	}
	_, fileName := filepath.Split(file)
	Log.Printf("move file old:%s,new:%s\n", file, filepath.Join(c.HtmlBakDir, string(os.PathSeparator), fileName))
	err = os.Rename(file, filepath.Join(c.HtmlBakDir, string(os.PathSeparator), fileName))
	if err != nil {
		Log.Error("osRename htmlbak Error DIR: %s, err:%s", c.HtmlBakDir, err)
	}
}

func waitMovefile() {
	for {
		select {
		case file, ok := <-MoveFileChan:
			if ok {
				go MoveFile(file)
			} else {
				break
			}
		}
	}
}

func waitHtmlFile(files []string) {
	for _, file := range files {
		Log.Printf("begin covert file:%s\n", file)
		resultNum := ReadHtml(file)
		switch resultNum {
		case 1:
			Log.Error("htmlFile %s, ERR:%s", file, "col num not 7 or 8")
		case 2:
			Log.Error("htmlFile %s, ERR:%s", file, "col not chinese")
		case 3:
			Log.Error("htmlFile %s, ERR:%s", file, "sql Gen err,sqlarr len is 0")
		case 4:
			Log.Error("htmlFile %s, ERR:%s", file, "insert or update have err")
		default:
			MoveFileChan <- file
			Log.Printf("finish covert file:%s\n", file)
		}

	}

}

func formatFields(index int, s string) (rstr string) {
	switch true {
	case index == 0 && strings.Contains(s, "Client"):
		rstr = "DATACLIENT"
		break
	case strings.Contains(s, "Phase(Write End Time)"):
		rstr = strings.Replace(s, "(Write End Time)", "", -1)
		break
	case strings.Contains(s, " (Compression Rate)"):
		rstr = strings.Replace(s, " (Compression Rate)", "", -1)
		break
	case strings.Contains(s, "(Space Saving Percentage)"):
		rstr = strings.Replace(s, "(Space Saving Percentage)", "", -1)
		break
	case strings.Contains(s, "(Current)"):
		rstr = strings.Replace(s, "(Current)", "", -1)
		break
	case s == "Agent /Instance":
		rstr = "AgentInstance"
		break
	case s == "Backup Set /Subclient":
		rstr = "BackupSetSubclient"
		break

	default:
		rstr = strings.TrimSpace(s)
	}
	return
}

func formatValues(index int, s string) (rstr string) {
	//fmt.Println(index, s)
	if strings.Contains(s, "N/A") {
		rstr = strings.Replace(s, "N/A", "", -1)
		s = rstr
	} else {
		rstr = s
	}
	if strings.Contains(s, "%") {
		rstr = strings.Replace(s, "%", "%%", -1)
		s = rstr
	} else {
		rstr = s
	}

	switch true {
	//case strings.Contains(s, "N/A"):
	//	rstr = strings.Replace(s, "N/A", "", -1)
	//	s = rstr
	case index == 2: // subclient

		if strings.Contains(s, "/") {
			Subclient = strings.Split(s, "/")[1]
		} else {
			Subclient = s
		}
		//fmt.Println("333333333333",index, s, Subclient)
		break

	}
	return
}

func formatTimeArr(index int, s string) {
	switch true {

	case index == 3: // time 赋值 全局变量，对应自定义列
		if strings.Contains(s, "(") {
			StartTime = strings.Split(s, "(")[0]
		} else {
			StartTime = s
		}
		StartTime = strings.TrimSpace(StartTime)
		StartTime = fmt.Sprint("to_date('" + StartTime + "','yyyy/mm/dd hh24:mi:ss')")

		TimeArr = append(TimeArr, StartTime)

		//StartTime = fmt.Sprintf("to date(%s mm/dd/yyyy hh24:mi:ss)",StartTime)
		//StartTime = "2019-08-12 04:00:00"
		break
	}
}

func ReadHtml(htmlfile string) (resultNum int) {

	f, err := os.OpenFile(htmlfile, os.O_RDONLY, 0755)
	if err != nil {
		Log.Println(err)
	}

	dom, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		//log.Fatal(err)
		Log.Printf("[Error]:%s", err)
	}
	//fmt.Println(dom)

	CommCells := dom.Find("body:contains(CommCell)")
	CommCells.Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			ss := strings.Split(strings.Split(selection.Text(), "CommCell:")[1], "--")[0]
			//fmt.Printf("CommCell :%+v\n", ss)
			CommCell = strings.TrimSpace(ss)
		}

	})
	//fmt.Println(CommCell)

	FiledsHeader := new([]string) // 列头，序号为KEY
	t_tbodyHeader := dom.Find("body >  table > tbody > tr:nth-child(1) > td")

	cv8 := false
	//有两种格式的HTML ,代表是CV8
	if len(t_tbodyHeader.Nodes) == 7 {
		cv8 = true
		//t_tbodyHeader = dom.Find("body > table:nth-child(12) > tbody > tr:nth-child(1) > td")
	}

	// 循环找出列头
	t_tbodyHeader.Each(func(i int, s *goquery.Selection) {
		sa := s.Text()
		sa = formatFields(i, sa)
		*FiledsHeader = append(*FiledsHeader, sa)
	})

	// 如果是CV8版本，手动修改列名, del size of Backup cols
	//if cv8 {
	//	tmp := make([]string, 16)
	//	copy(tmp[:], (*FiledsHeader)[:])
	//
	//	*FiledsHeader = tmp[0:9]
	//	//fmt.Println(*FiledsHeader)
	//	//*FiledsHeader = append(*FiledsHeader, []string{"Data Transferred", "Data Written"}...)
	//	//fmt.Println(*FiledsHeader)
	//	*FiledsHeader = append(*FiledsHeader, tmp[10:]...)
	//	//fmt.Println(*FiledsHeader)
	//	(*FiledsHeader)[3] = "Job ID (CommCell)(Status)"
	//}
	//fmt.Println("222222222",*FiledsHeader)
	//for i := 0; i < len(*FiledsHeader); i++ {
	//	fmt.Printf("%d, value:%s\n", i, (*FiledsHeader)[i])
	//}

	switch true {
	case len(*FiledsHeader) < 7 || len(*FiledsHeader) > 8: // col not enough
		//fmt.Println(FiledsHeader)
		//fmt.Println(len(*FiledsHeader))
		return 1
	case (*FiledsHeader)[0] != "严重性": // col not English
		return 2
	}

	//列名从英文 改成中文
	if cv8 {
		*FiledsHeader = []string{"COMMCELL", "SEVERITY", "EVENT ID", "JOB ID", "PROGRAM", "COMPUTER", "DESCRIPTION", "TIME"}
	} else {
		*FiledsHeader = []string{"COMMCELL", "SEVERITY", "EVENT ID", "JOB ID", "PROGRAM", "COMPUTER", "EVENT CODE", "DESCRIPTION", "TIME"}
	}

	//fmt.Println(cv8)
	//fmt.Printf("%+v\n", *FiledsHeader)

	//*FiledsHeader = append(*FiledsHeader, []string{"START TIME", "COMMCELL", "APPLICATIONSIZE", "DATASUBCLIENT"}...)

	var t_tbodyData *goquery.Selection

	//if cv8 {
	//t_tbodyData = dom.Find("body > table:nth-child(12) > tbody > tr")
	//} else {
	t_tbodyData = dom.Find("body > table > tbody > tr")
	//}

	//fmt.Println(t_tbodyData.Html())

	headlen := len(*FiledsHeader)
	Data := make([][]string, 0)
	//循环表格数据
	t_tbodyData.Each(func(i int, s *goquery.Selection) {
		if i == 0 { // 0行是列名
			return
		}
		tmpSubData := make([]string, 0)
		tmpSubData = append(tmpSubData, []string{CommCell}...)
		sa := s.Find("td")
		sa.Each(func(i int, selection *goquery.Selection) {
			ss1 := selection.Text()
			//switch true {
			//case strings.Contains(ss1,"N/A"):
			//	ss1=strings.Replace(ss1,"N/A","",-1)
			//}
			if i == 3 {
				formatTimeArr(i, ss1)
			} else {
				ss1 = formatValues(i, ss1)
				tmpSubData = append(tmpSubData, ss1)
			}


			//fmt.Println(i,ss1)
		})
		//tmpSubData = append(tmpSubData, []string{StartTime, CommCell, ApplicationSize, Subclient}...)

		//if cv8 {
		//	//fmt.Printf("%d,%d\n",len(tmpSubData),len(*FiledsHeader))
		//	if len(tmpSubData) == headlen+1 { // 因为列头删除了一列
		//		tmp := tmpSubData[:]
		//		tmpSubData = tmp[0:9]
		//		//tmpSubData = append(tmpSubData, []string{"", ""}...)
		//		tmpSubData = append(tmpSubData, tmp[10:]...)
		//		//fmt.Printf("%d,%v\n",len(tmpSubData),tmpSubData)
		//		Data = append(Data, tmpSubData)
		//	}

		////} else {
		//fmt.Println(len(tmpSubData),headlen)
		if len(tmpSubData) +1 == headlen { // 数据列数要与 列头一样长
			Data = append(Data, tmpSubData)
		}
		//}

	})

	f.Close()
	//fmt.Println(FiledsHeader)

	return GenSqls(Data, FiledsHeader, htmlfile) // 这里应该包括数据库插入 ，可以是后台 go 后面

}

func sliceToString(sl []string) (returnstring string) {
	sl1 := sl[0:len(sl)]
	for i := 0; i < len(sl1); i++ {
		if i == len(sl1)-1 {
			returnstring = returnstring
			returnstring = returnstring + "\"" + sl1[i] + "\"" + ")"
		} else {
			returnstring = returnstring
			returnstring = returnstring + "\"" + sl1[i] + "\"" + ","
		}
	}
	return
}

func sliceToStringValue(sl []string) (returnstring string) {
	sl1 := sl[0:len(sl)]
	for i := 0; i < len(sl1); i++ {
		if i == len(sl1)-1 {
			returnstring = returnstring
			returnstring = returnstring + "\"" + sl1[i] + "\""
		} else {
			returnstring = returnstring
			returnstring = returnstring + "\"" + sl1[i] + "\"" + ","
		}
	}
	return
}

func updatesliceToString(sl []string, value []string) (returnstring string) {
	wherestring := ""
	if (len(sl) -1 == len(value)) {
		for i := 0; i < len(sl); i++ {
			if len(value) > i{
				if (value[i] == "") { // 空值路过
					continue
				}
			}
			if i == len(sl)-1 {
				//returnstring = returnstring + "\"" + sl[i] + "\"" + "='" + value[i] + "'"
				//wherestring = wherestring + "\"" + sl[i] + "\"" + "='" + value[i] + "'"
				returnstring = returnstring + "\"" + sl[i] + "\"" + "=%s"
				wherestring = wherestring + "\"" + sl[i] + "\"" + "=%s"
			} else {
				returnstring = returnstring + "\"" + sl[i] + "\"" + "='" + value[i] + "',"
				wherestring = wherestring + "\"" + sl[i] + "\"" + "='" + value[i] + "' and "
			}

		}
	}
	returnstring = returnstring + " where " + wherestring
	return
}

func GenSqls(Data [][]string, header *[]string, htmlfile string) (resultNum int) {
	baseSql := "insert into HF_EVENTREPORT(" + sliceToString(*header) + " values("
	updateSql := "update HF_EVENTREPORT set "
	//fmt.Println(Data)
	//fmt.Println(updateSql)
	//fmt.Println(baseSql)
	sqlArr := make([]map[string]string, 0)
	for _, value := range Data {
		tmpmap := make(map[string]string, 0)
		//fmt.Println(index, baseSql+sliceToString(value))
		valueStr := strings.Replace(sliceToStringValue(value), "\"", "'", -1)
		//fmt.Println(valueStr)
		//sqlArr = append(sqlArr, baseSql+valueStr)
		tmpmap["insert"] = baseSql + valueStr
		tmpmap["update"] = updateSql + updatesliceToString(*header, value)
		sqlArr = append(sqlArr, tmpmap)
	}

	return stmtSql(sqlArr, htmlfile)
}

func stmtSql(sqlArr []map[string]string, htmlfile string) (resultNum int) {
	os.Setenv("NLS_LANG", "")
	//if len(os.Args) != 2 {
	//	log.Fatalln(os.Args[0] + " user/password@host:port/sid")
	//}

	db, err := sql.Open("oci8", c.OracleDsn)
	//fmt.Printf("%+v\n",db)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	updatenum := 0
	insertnum := 0
	if len(sqlArr) == 0 {
		return 3
	}
	for i := 0; i < len(sqlArr); i++ {
		//fmt.Println(fmt.Sprintf(sqlArr[i]["insert"], TimeArr[i]))
		//fmt.Println(sql["insert"])
		//r,err:=db.Exec(fmt.Sprintf("insert into BCD(\"Client\",\"AgentInstance\",\"BackupSetSubclient\") values('a','b',%d)", 1))
		Result, err := db.Exec(fmt.Sprintf(sqlArr[i]["update"], TimeArr[i],TimeArr[i]))
		//fmt.Println(Result.LastInsertId())
		//fmt.Println(Result.RowsAffected())
		//fmt.Println(Result, err)
		if err != nil {
			Log.Error("HtmlFile:%s ,Sql Exec Err:%s", htmlfile, err)
			resultNum = 4
			continue
		}
		updatenum = updatenum + 1
		r, e := Result.RowsAffected()
		if err != nil {
			Log.Error("HtmlFile:%s ,Sql Exec Err:%s", htmlfile, e)
			resultNum = 4
			continue
		}

		if r == 0 {
			//fmt.Println(sqlArr[i]["insert"]+",%s)")
			_, err := db.Exec(fmt.Sprintf(sqlArr[i]["insert"]+",%s)", TimeArr[i]))
			if err != nil {
				Log.Error("HtmlFile:%s ,Sql Exec Err:%s", htmlfile, err)
				resultNum = 4
				continue
			}
			insertnum = insertnum + 1
		}

	}
	Log.Printf("htmlfile : %s ,total sql:%d,success update sql:%d,insert sql:%d\n", htmlfile, len(sqlArr), updatenum, insertnum)
	return

}
