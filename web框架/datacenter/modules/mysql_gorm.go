package modules

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	//"time"
)

/**
 * @Author: Tao Jun
 * @Description: modules
 * @File:  mysql_gorm
 * @Version: 1.0.0
 * @Date: 2021/3/22 下午2:49
 */
// db连接
//var db *gorm.DB
//var DbPool *sql.DB
var MysqlDb *gorm.DB

func MysqlInitConn() (err error) {
	MysqlDb, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:+mysql2016@tcp(192.168.43.179:3306)/myapp_test?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                         // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                        // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                        // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                        // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                       // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		return err
	}

	MysqlDb.DB()
	//XXX 连接池由sql.db包提供
	DbPool, err := MysqlDb.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	DbPool.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	DbPool.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	DbPool.SetConnMaxLifetime(time.Hour)
	return DbPool.Ping()

}

func DoQuerySort(rows *sql.Rows) (result map[int]map[string]string) {
	// 返回所有列
	columns, _ := rows.Columns()
	// 这里标识一行所有列的值，用 []byte 表示
	vals := make([][]byte, len(columns))
	// 这里标识一行填充数据
	scans := make([]interface{}, len(columns))
	// 这里scans 引用vals 把数据填充到 []byte 里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result = make(map[int]map[string]string)
	for rows.Next() {
		// 填充数据
		rows.Scan(scans...)
		// 每行数据
		row := make(map[string]string)
		// 把vals中的数据赋值到row中
		for k, v := range vals {
			key := columns[k]
			// 这里把[]byte数据转化成string
			row[key] = string(v)
		}
		// 放入结果集
		result[i] = row
		i++
	}
	return

}

func DoQuery(rows *sql.Rows) (result []map[string]string) {
	// 返回所有列
	columns, _ := rows.Columns()
	// 这里标识一行所有列的值，用 []byte 表示
	vals := make([][]byte, len(columns))
	// 这里标识一行填充数据
	scans := make([]interface{}, len(columns))
	// 这里scans 引用vals 把数据填充到 []byte 里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	result = make([]map[string]string, 0)
	for rows.Next() {
		// 填充数据
		rows.Scan(scans...)
		// 每行数据
		row := make(map[string]string)
		// 把vals中的数据赋值到row中
		for k, v := range vals {
			key := columns[k]
			// 这里把[]byte数据转化成string
			row[key] = string(v)
		}
		// 放入结果集
		result = append(result, row)

	}
	return

}
