package db

import (
	"tjtools/db/oracle"
	"godev/works/train/g"
	_ "github.com/mattn/go-oci8"
	"os"
)
//var DB *sql.DB
//
//func Init() {
//	var err error
//	DB, err = sql.Open("mysql", g.Config().Database)
//	if err != nil {
//		log.Fatalln("open db fail:", err)
//	}
//
//	DB.SetMaxOpenConns(g.Config().MaxConns)
//	DB.SetMaxIdleConns(g.Config().MaxIdle)
//
//	err = DB.Ping()
//	if err != nil {
//		log.Fatalln("ping db fail:", err)
//	}
//}




func Init() {
	os.Setenv("NLS_LANG", "")
	Engine, err := oracle.NewEngine(g.Config().OracleDsn, "runsql.log")
	if err!=nil{
		g.Logger().Error("db init Fail")
		return
	}
	g.Logger().Printf("db success")
	//Engine.ShowExecTime(true)
	//Engine.ShowSQL(true)
	//err = Engine.Ping()
	Engine.ShowExecTime(true)
	if err != nil {
		g.Logger().Error("ping db fail:", err)
		}
	g.Engine = Engine

}



