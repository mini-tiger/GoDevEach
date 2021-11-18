package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  one2many
 * @Version: 1.0.0
 * @Date: 2021/11/16 上午9:31
 */
// xxx one(KpiInfo) many(kpiset)
// xxx https://gobea.cn/blog/detail/260JLw5Z.html

type MonitorType struct {
	Id               int       `gorm:"primaryKey" json:"id"`
	KpiSetID         []KpiSet  `json:"kpi_set_id"`
	MonitorTypeName  string    `gorm:"column:monitor_type_name; size:128" json:"monitor_type_name"`
	MonitorTypeIdent string    `gorm:"column:monitor_type_ident; size:64" json:"monitor_type_ident"`
	MonitorTypeDesc  string    `gorm:"column:monitor_type_desc; size:512" json:"monitor_type_desc"`
	CreatedAt        time.Time `sql:"index"`
}

func (MonitorType) TableName() string {
	return "aaaa_itm_t_monitor_type"
}

type KpiInfo struct {
	ID uint64 `gorm:"primaryKey" json:"id"`
	//MonitorTypeId int      `gorm:"column:monitor_type_id" json:"monitor_type_id"` //多余
	KpiSetId uint64 `json:"kpi_set_id"` // 与 kpiset 关联字段
	KpiKey   string `gorm:"column:kpi_key; size:128" json:"kpi_key"`
	KpiName  string `gorm:"column:kpi_name; size:128" json:"kpi_name"`
	KpiUint  string `gorm:"column:kpi_uint; size:32" json:"kpi_uint"`
}

func (KpiInfo) TableName() string {
	return "AAA_itm_t_kpi"
}

type KpiSet struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	KpiInfoID      []KpiInfo `json:"kpi_info_id"`                                   // kpiinfo 多 数据库无此字段
	MonitorTypeId  int       `gorm:"column:monitor_type_id" json:"monitor_type_id"` // MonitorType 一
	KpiSetName     string    `gorm:"column:kpi_set_name; size:128" json:"kpi_set_name"`
	KpiSetStatus   int       `gorm:"column:kpi_set_status" json:"kpi_set_status"`
	MonitorCycleId int       `gorm:"column:monitor_cycle_id" json:"monitor_cycle_id"`
	KpiSetIdent    string    `gorm:"column:kpi_set_ident; size:64" json:"kpi_set_ident"`
	KpiSetDesc     string    `gorm:"column:kpi_set_desc; size:512" json:"kpi_set_desc"`
	//CreateTime     int64  `gorm:"column:create_time" json:"create_time"`
	KpiSetType string    `gorm:"column:kpi_set_type; size:32" json:"kpi_set_type"`
	CreatedAt  time.Time `sql:"index"`
	//UpdatedAt time.Time `sql:"index"`
}

func (KpiSet) TableName() string {
	return "AAAA_itm_t_kpi_set"
}

var Db *gorm.DB

func init() {
	var err error
	//admin := config.ViperCfg.Mysql
	//dsn := admin.Username+":"+admin.Password+"@tcp("+admin.Path+")/"+admin.Dbname+"?"+admin.Config
	//dsn:="root:hello@tcp(172.16.71.17:3306)/itgo_monitor?charset=utf8&parseTime=True&loc=Local"
	//if MysqlClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
	//	log.Fatal("mysql init err:", err)
	//}
	Db, err = gorm.Open("mysql", "root:hello@tcp(172.16.71.17:3306)/itgo_monitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
	//if Db, err = gorm.Open(mysql.New(mysql.Config{
	//	//DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
	//	DSN: dsn,
	//	//DefaultStringSize: 64, // default size for string fields
	//
	//	DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
	//	DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
	//	DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
	//	SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	//}), &gorm.Config{}); err != nil {
	//	log.Fatal("mysql init err:", err)
	//}

	// 迁移数据库表
	DBTables()
}

// 注册数据库表专用
func DBTables() {
	//if err := Db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
	//	&KpiInfo{},
	//	&KpiSet{},
	//
	//
	//); err != nil {
	//	log.Println("migrate mysql table err:", err)
	//}
	//if !Db.Migrator().HasTable(&KpiInfo{}){
	//	Db.Migrator().CreateTable(&KpiInfo{})
	//}
	//if !Db.Migrator().HasTable(&KpiSet{}){
	//	Db.Migrator().CreateTable(&KpiSet{})
	//}
	//// xxx 创建外键
	//if !Db.Migrator().HasConstraint(&KpiInfo{},"KpiSetId"){
	//	Db.Migrator().CreateConstraint(&KpiInfo{},"KpiSetId")
	//}
	Db.AutoMigrate(&MonitorType{}, &KpiSet{}, &KpiInfo{})
}

func main() {

}
