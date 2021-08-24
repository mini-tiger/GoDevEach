package controllers

import (
	"fmt"
	"goweb/models"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  mysqlSearch
 * @Version: 1.0.0
 * @Date: 2021/4/21 下午2:06
 */

func registerMysqlControllers() {
	http.Handle("/mysqlselect", http.HandlerFunc(MysqlSelect))
}

type sysAlgo struct {
	Algoname     string
	State        int
	Algofullname string
}

func MysqlSelect(writer http.ResponseWriter, request *http.Request) {

	sql := "select algoname,`state`,algofullname from sys_algo where description like ?"

	rows, err := models.DB.Query(sql, "%回归")

	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	var resultData []interface{} = make([]interface{}, 0)
	for rows.Next() {
		var s sysAlgo
		err := rows.Scan(&s.Algoname, &s.State, &s.Algofullname)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("mysql select %+v\n", s)
		resultData = append(resultData, s)
	}

	ResponseJsonDataSuccess(writer, request, resultData)
}
