package modules

import (
	"database/sql"
	"gitee.com/taojun319/tjtools/db/oracle"
	"haifei/syncHtmlYWReport/g"
	"time"
)

var DB *sql.DB

func GetDBConn() {
	for {
		db, err := sql.Open("oci8", g.GetConfig().OracleDsn)
		if ConnErr := oracle.CheckOracleConn(&g.GetConfig().OracleDsn); ConnErr != nil || err != nil {
			_ = Log.Error("DB Conn err:%v ,SqlExecError:%s,DB Args:%+v ,%d 秒重试 \n",
				err, ConnErr,
				g.GetConfig().OracleDsn,
				g.GetConfig().TimeInter)
			time.Sleep(time.Duration(g.GetConfig().TimeInter) * time.Second)
		} else {
			DB = db
			break
		}
	}
}

func CloseConn() {
	if DB.Ping() == nil {
		_ = DB.Close()
	}
}
