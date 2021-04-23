package models

import (
	"database/sql"

	"log"
)

/**
 * @Author: Tao Jun
 * @Description: model
 * @File:  mysqlconn
 * @Version: 1.0.0
 * @Date: 2021/4/21 下午2:36
 */

// 定义一个全局对象db
var DB *sql.DB

func init() {
	// DSN:Data Source Name
	dsn := "kaleido:123456@tcp(192.168.43.152:3306)/kaleido?charset=utf8&parseTime=True&loc=Local"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db open err:%s\n", err)
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = DB.Ping()
	if err != nil {
		log.Fatalln("db conn err")
	}
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(20)

}
