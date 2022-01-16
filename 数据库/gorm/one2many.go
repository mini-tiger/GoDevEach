package main

import (
	"fmt"
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
type KpiInfo struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	MonitorTypeId int       `gorm:"column:monitor_type_id" json:"monitor_type_id"`
	KpiSetId      []*KpiSet `json:"kpi_set_id"` //数据库不会 有这个字段
	KpiKey        string    `gorm:"column:kpi_key; size:128" json:"kpi_key"`
	KpiName       string    `gorm:"column:kpi_name; size:128" json:"kpi_name"`
	KpiUint       string    `gorm:"column:kpi_uint; size:32" json:"kpi_uint"`
}

func (KpiInfo) TableName() string {
	return "AAA_itm_t_kpi"
}

type KpiSet struct {
	ID             uint64 `gorm:"primaryKey" json:"id"`
	KpiInfoID      uint64 `gorm:"column:kpi_info_id;not null" json:"kpi_info_id"` // 与 kpiset 关联字段
	MonitorTypeId  int    `gorm:"column:monitor_type_id" json:"monitor_type_id"`
	KpiSetName     string `gorm:"column:kpi_set_name; size:128" json:"kpi_set_name"`
	KpiSetStatus   int    `gorm:"column:kpi_set_status" json:"kpi_set_status"`
	MonitorCycleId int    `gorm:"column:monitor_cycle_id" json:"monitor_cycle_id"`
	KpiSetIdent    string `gorm:"column:kpi_set_ident; size:64" json:"kpi_set_ident"`
	KpiSetDesc     string `gorm:"column:kpi_set_desc; size:512" json:"kpi_set_desc"`
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
	Db, err = gorm.Open("mysql", "root:FuTongde1gesjk@tcp(172.16.71.31:3306)/itgo_monitor?charset=utf8&parseTime=True&loc=Local")
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
	Db.AutoMigrate(&KpiInfo{}, &KpiSet{})
}

func main() {
	db := Db
	var kpiinfo KpiInfo
	//xxx create
	// 一起创建
	kpiinfo = KpiInfo{
		KpiKey: "Abc",
		KpiSetId: []*KpiSet{
			{
				KpiSetType: "123",
			},
			{
				KpiSetType: "456",
			},
		},
	}
	db.Create(&kpiinfo)
	//db.Save(&kpiinfo)
	//xxx 先 多  在  一

	KpiSets := []*KpiSet{
		&KpiSet{
			KpiSetType: "123",
		},
		&KpiSet{
			KpiSetType: "456",
		},
	}
	num := time.Now().Unix()
	//
	//// 先创建 一
	kpiinfo = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}
	kpiinfo.KpiSetId = KpiSets // 多 绑定到   一
	db.Create(&kpiinfo)

	// ========
	// = xxx 查询 =
	// ========

	kpiinfo = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}

	// 先找一
	db.Where(&kpiinfo).First(&kpiinfo)

	// 通过一 找 多 方法A
	var kpisetss []KpiSet
	db.Model(&kpiinfo).Association("KpiSetId").Find(&kpisetss)

	n := db.Model(&kpiinfo).Association("KpiSetId").Count()
	fmt.Printf("%+v\n", kpiinfo)
	// xxx 总数
	fmt.Printf("关联数量:%d,%+v\n", n, kpisetss)

	//xxx  通过一 找 多 方法B
	var kpiinfo321 = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}
	var kpiset321s []KpiSet = make([]KpiSet, 0)

	db.Where(&kpiinfo321).First(&kpiinfo321)
	db.Model(&kpiinfo321).Related(&kpiset321s)
	var ss []uint64 = make([]uint64, 0)

	for _, vv := range kpiset321s {
		ss = append(ss, vv.ID)
	}
	fmt.Printf("%+v\n", kpiinfo321)
	fmt.Printf("%v\n", ss)
	//time.Sleep(100*time.Second)

	// 预加载分两条查询语句
	//xxx  Preload   通过一  , 返回 一 和  多那边的 数据
	fmt.Println("===============================================")
	var kpiInfo111 = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}

	db.Where(&kpiInfo111).Find(&kpiInfo111)

	db.Model(&kpiInfo111).
		Preload("KpiSetId").
		Find(&kpiInfo111)

	fmt.Printf("%+v\n", kpiInfo111)

	fmt.Println("===============================================")
	var kpiInfo113 = []KpiInfo{}
	var kpiInfo112 = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}

	db.Where(&kpiInfo112).Find(&kpiInfo113) // 一这边 匹配多个，有多个 Abc1匹配

	db.Model(&kpiInfo113).
		Preload("KpiSetId").
		Find(&kpiInfo113)

	fmt.Printf("%+v\n", kpiInfo113)
	fmt.Println("===============================================")

	//xxx 多 找一
	kpiset44 := KpiSet{}
	db.First(&kpiset44, "id = 1005") // 从数据库找一条数据

	db.Model(&kpiset44).Related(&kpiset44, "kpi_info_id") // Related 返回 一 那边的数据

	fmt.Printf("关联 一那边 的id:%d \n", kpiset44.KpiInfoID)
	fmt.Printf("%+v\n", kpiset44)

	//for _, user := range kpiset44.KpiInfoID {
	//	fmt.Print(user.Name + "  ")
	//}
	fmt.Println("===============================================")

	// ========
	// = xxx 更新 =
	// ========
	// 使用创建新关联替换当前关联
	// 单个
	KpiSets123 := []*KpiSet{
		&KpiSet{
			KpiSetType: "1231",
		},
		&KpiSet{
			KpiSetType: "4561",
		},
	}

	var kpiInfo114 = KpiInfo{}
	db.Where("kpi_key=?", fmt.Sprintf("Abc%d", num)).First(&kpiInfo114)
	fmt.Println(kpiInfo114)
	// xxx 建议 删除 旧的数据
	db.Model(&kpiInfo114).Association("KpiSetId").Replace(KpiSets123) // 更新   多那这的数据不会删除
	fmt.Printf("%+v\n", kpiInfo114)

	// 多个
	// languages := []*Language{}
	// db.Where("id IN (?)", []uint{1, 2}).Find(&languages)
	// db.Model(&user).Association("Languages").Replace(languages)
	// fmt.Println(user)

	fmt.Println("===============================================")
	// ========
	// = xxx 清除 =
	// ========
	// 清空对关联的引用，不会删除关联数据本身
	var kpiInfo116 = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}

	// 删除关联的引用，不会删除关联本身
	// &kpiInfo116.KpiSetId 普通查询为空，需要关联查询得到
	db.Where(&kpiInfo116).First(&kpiInfo116)

	//db.Model(&kpiInfo116).Association("KpiSetId").Clear()

	//db.Model(&kpiInfo116).Association("KpiSetId").Delete(&kpiInfo116.KpiSetId)

	fmt.Printf("%+v,%d\n", kpiInfo116.ID, len(kpiInfo116.KpiSetId))

	// 清除关联后
	db.Model(&kpiInfo116).Related(&kpiInfo116.KpiSetId).Find(&kpiInfo116.KpiSetId)
	fmt.Printf("%+v,%d\n", kpiInfo116.ID, len(kpiInfo116.KpiSetId))

	// ========
	// =xxx 删除 =
	// ========

	kpiInfo117 := KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}

	// 软删除 https://jasperxu.com/gorm-zh/crud.html#d
	// 如果模型有DeletedAt字段，它将自动获得软删除功能！ 那么在调用Delete时不会从数据库中永久删除，而是只将字段DeletedAt的值设置为当前时间,否则 真正删除

	// 软删除时，清空引用后再执行删除，没有Deletet字段,真正删除

	//db.Debug().Delete(&kpiInfo116) // 一 这边删除
	//默认软删除，Unscoped()记录真正删除

	//db.Debug().Unscoped().Delete(&kpiInfo116) // 一 这边删除
	// fmt.Println(user)

	//fmt.Println(len(kpisetss11))
	//db.Debug().Unscoped().Delete(&kpisetss11) // 一 这边删除,没有

	// xxx 先查找多，并物理删除
	db.Where(&kpiInfo117).Find(&kpiInfo117)
	var kpisetss333 []KpiSet
	db.Model(&kpiInfo117).Association("KpiSetId").Find(&kpisetss333)

	//db.Debug().Unscoped().Delete(kpisetss333)

	fmt.Println(len(kpisetss333))
	for _, value := range kpisetss333 {
		db.Debug().Unscoped().Delete(value)
	}

	// xxx 物理删除一 这边的数据
	db.Debug().Select("KpiSetId").Unscoped().Delete(&kpiInfo117)
	//db.Select("KpiSetId").Delete(&kpiInfo117)

}
