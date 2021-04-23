package g

import "gitee.com/taojun319/tjtools/nmap"

//var FieldsMap map[string]string=map[string]string{"客户端":"ReportClient","主机名":"HOSTNAME","总作业数":"TOTALJOB",
//													"已完成":"COMPLETED","完成但有错误":"COMPLETIONERROR","完成但有警告":"COMPLETIONWARN",
//													"已终止":"TERMINATION","不成功":}
var AllChEngSummaryMap *nmap.SafeMap = nmap.NewSafeMap()
const (
	Cv10  = 10
	Cv100 = 100
	Cv8   = 8
	Cv11  = 11
	Cv101 = 101
)

func init() {

	AllSummary := map[string]interface{}{"客户端": "REPORTCLIENT", "主机名": "HOSTNAME", "总作业数": "TOTALJOB",
		"已完成": "COMPLETED", "完成但有错误": "COMPLETEDWITHERRORS", "完成但有警告": "COMPLETEDWITHWARNINGS",
		"已终止": "KILLE", "不成功": "UNSUCCESSFUL", "运行中": "RUNNING", "已延迟": "DELAYED",
		"未运行": "NORUN", "无计划": "NOSCHEDULE", "提交": "COMMITTED", "应用程序大小（压缩率）": "SIZEOFAPPLICATION",
		"写入数据（空间节省百分比）": "DATAWRITTEN", "开始时间": "STARTTIME", "结束时间": "ENDTIME", "受保护对象": "PROTECTEDOBJECTS",
		"失败对象": "FAILEDOBJECTS", "失败文件夹": "FAILEDFOLDERS"}
	AllChEngSummaryMap.M = AllSummary
	//var a *nmap.SafeMap=&nmap.SafeMap{M:AllChEngSummaryMap}

	//var aa *nmap.SafeMap = &nmap.NewSafeMap(a)
}

const (
	DetailTableName  = "HF_YWREPORTDETAIL"
	SummaryTableName = "HF_YWREPORTSUMMARY"
	//TimeLayout       = "2006-01-02 15:04:05" // status 需要使用

)

var NeedFields []string = []string{"DATACLIENT", "AgentInstance", "BackupSetSubclient", "Job ID (CommCell)(Status)",
	"Type", "Scan Type", "Start Time(Write Start Time)", "End Time or Current Phase", "Size of Application"}
const (
	TimeLayoutChi = "2006/01/02 15:04:05"
	TimeLayoutEng = "01/02/2006 15:04:05"
	BaseCheckHangSql = "select ID,case APPLICATIONSIZE when null then null else ROUND(APPLICATIONSIZE-%s,2) end as diff from %s "+
"where \"Job ID (CommCell)(Status)\" ='%s' and RUNTYPE='活动' and FILESYNCTIME < %s and ROWNUM< 2 "+
"order by FILESYNCTIME desc"
	)

// 摘要表
var SummaryFieldsMap map[int]string = map[int]string{0: "REPORTCLIENT", 1: "HOSTNAME", 2: "TOTALJOB",
	3: "COMPLETED", 4: "COMPLETEDWITHERRORS", 5: "COMPLETEDWITHWARNINGS", 6: "KILLED", 7: "UNSUCCESSFUL", 8: "RUNNING", 9: "DELAYED",
	10: "NORUN", 11: "NOSCHEDULE", 12: "COMMITTED", 13: "SIZEOFAPPLICATION", 14: "DATAWRITTEN", 15: "STARTTIME", 16: "ENDTIME", 17: "PROTECTEDOBJECTS",
	18: "FAILEDOBJECTS", 19: "FAILEDFOLDERS"}

// 摘要表cv10    缺少 1，12
var SummaryFieldsMapCv10 map[int]string = map[int]string{0: "REPORTCLIENT", 1: "TOTALJOB",
	2: "COMPLETED", 3: "COMPLETEDWITHERRORS", 4: "COMPLETEDWITHWARNINGS", 5: "KILLED", 6: "UNSUCCESSFUL", 7: "RUNNING", 8: "DELAYED",
	9: "NORUN", 10: "NOSCHEDULE", 11: "SIZEOFAPPLICATION", 12: "DATAWRITTEN", 13: "STARTTIME", 14: "ENDTIME", 15: "PROTECTEDOBJECTS",
	16: "FAILEDOBJECTS", 17: "FAILEDFOLDERS"}

// 摘要表cv8   缺少 1，5，9，12
var SummaryFieldsMapCv8 map[int]string = map[int]string{0: "REPORTCLIENT", 1: "TOTALJOB",
	2: "COMPLETED", 3: "COMPLETEDWITHERRORS", 4: "KILLED", 5: "UNSUCCESSFUL", 6: "RUNNING",
	7: "NORUN", 8: "NOSCHEDULE", 9: "SIZEOFAPPLICATION", 10: "DATAWRITTEN", 11: "STARTTIME", 12: "ENDTIME", 13: "PROTECTEDOBJECTS",
	14: "FAILEDOBJECTS", 15: "FAILEDFOLDERS"}

var SummaryFieldsMapPlus map[int]string = map[int]string{20: "COMMCELL", 21: "REPORTTIME", // 这两个字段从页面上开头获取,唯一联合字段，也是和详细表 关系的字段
	22: "HTMLFILE", 23: "INSERTTIME"}

var SuccessColors map[string]string = map[string]string{"#CCFFCC": "0", "#66ABDD": "0", "#CC9999": "0", "93C54B": "0", "#CCFFFF": "0",
	"#FFFFFF": "0", "#ccffcc": "0", "#66abdd": "0", "#cc9999": "0", "93c54b": "0", "#ccffff": "0",
	"#ffffff": "0"}

// 详细数据表
var DetailFieldsMap map[int]string = map[int]string{0: "DATACLIENT", 1: "AgentInstance", 2: "BackupSetSubclient", 3: "Job ID (CommCell)(Status)",
	4: "Type", 5: "Scan Type", 6: "Start Time(Write Start Time)", 7: "End Time or Current Phase", 8: "Size of Application", 9: "Data Transferred",
	10: "Data Written", 11: "Data Size Change", 12: "Transfer Time", 13: "Throughput (GB/Hour)", 14: "Protected Objects", 15: "Failed Objects",
	16: "Failed Folders"}

// 详细数据表 cv8,去掉第9,10去掉.第10列更换为  "Data Size Change"
var DetailFieldsMapCv8 map[int]string = map[int]string{0: "DATACLIENT", 1: "AgentInstance", 2: "BackupSetSubclient", 3: "Job ID (CommCell)(Status)",
	4: "Type", 5: "Scan Type", 6: "Start Time(Write Start Time)", 7: "End Time or Current Phase", 8: "Size of Application",
	10: "Data Size Change", 11: "Transfer Time", 12: "Throughput (GB/Hour)", 13: "Protected Objects", 14: "Failed Objects",
	15: "Failed Folders"}

var DetailFieldsMapPlus map[int]string = map[int]string{
	17: "COMMCELL", 18: "START TIME", 19: "DATASUBCLIENT", 20: "APPLICATIONSIZE", //作业状态  ,通过开始时间 和 子客户端，格式化出来的字段
	21: "HTMLFILE", 22: "INSERTTIME", 23: "REASONFORFAILURE", 24: "COLOR", 25: "RUNTYPE", 26: "REPORTTIME", 27: "SOLVETYPE", 28: "END TIME", // 解决时间，工程师， 这几个字段不用插入数据，
	29:"FILESYNCTIME"}

//todo 有问题的行 解决状态默认是未解决
