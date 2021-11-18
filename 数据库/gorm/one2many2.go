package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	ID            uint64 `gorm:"primaryKey" json:"id"`
	MonitorTypeId int    `gorm:"column:monitor_type_id" json:"monitor_type_id"`
	KpiSetId      int    `gorm:"column:kpi_set_id" json:"kpi_set_id"` // 与 kpiset 关联字段
	KpiKey        string `gorm:"column:kpi_key; size:128" json:"kpi_key"`
	KpiName       string `gorm:"column:kpi_name; size:128" json:"kpi_name"`
	KpiUint       string `gorm:"column:kpi_uint; size:32" json:"kpi_uint"`
}

func (KpiInfo) TableName() string {
	return "AAA_itm_t_kpi"
}

type KpiSet struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	KpiInfoID      []*KpiInfo `json:"kpi_info_id"` // kpiinfo 多 数据库无此字段
	MonitorTypeId  int        `gorm:"column:monitor_type_id" json:"monitor_type_id"`
	KpiSetName     string     `gorm:"column:kpi_set_name; size:128" json:"kpi_set_name"`
	KpiSetStatus   int        `gorm:"column:kpi_set_status" json:"kpi_set_status"`
	MonitorCycleId int        `gorm:"column:monitor_cycle_id" json:"monitor_cycle_id"`
	KpiSetIdent    string     `gorm:"column:kpi_set_ident; size:64" json:"kpi_set_ident"`
	KpiSetDesc     string     `gorm:"column:kpi_set_desc; size:512" json:"kpi_set_desc"`
	//CreateTime     int64  `gorm:"column:create_time" json:"create_time"`
	KpiSetType string    `gorm:"column:kpi_set_type; size:32" json:"kpi_set_type"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
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
	dsn := "root:hello@tcp(172.16.71.17:3306)/itgo_monitor?charset=utf8&parseTime=True&loc=Local"
	//if MysqlClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
	//	log.Fatal("mysql init err:", err)
	//}
	//Db, err = gorm.Open("mysql", "root:hello@tcp(172.16.71.17:3306)/itgo_monitor?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	if Db, err = gorm.Open(mysql.New(mysql.Config{
		//DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
		DSN: dsn,
		//DefaultStringSize: 64, // default size for string fields

		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{}); err != nil {
		log.Fatal("mysql init err:", err)
	}

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

	if !Db.Migrator().HasTable(&KpiSet{}) {
		Db.Migrator().CreateTable(&KpiSet{})
	}
	if !Db.Migrator().HasTable(&KpiInfo{}) {
		Db.Migrator().CreateTable(&KpiInfo{})
	}
	//// xxx 创建外键
	if !Db.Migrator().HasConstraint(&KpiSet{}, "KpiInfoID") {
		Db.Migrator().CreateConstraint(&KpiSet{}, "KpiInfoID")
	}
	//Db.AutoMigrate(&KpiInfo{}, &KpiSet{})

}

func main() {
	db := Db
	var kpiset KpiSet
	////xxx create
	//// 一起创建
	kpiset = KpiSet{
		KpiSetDesc: "Abc",
		KpiInfoID: []*KpiInfo{
			{
				KpiKey: "123",
			},
			{
				KpiKey: "456",
			},
		},
	}
	db.Create(&kpiset)
	////db.Save(&kpiinfo)
	////xxx 先 多  在  一
	//
	KpiInfos := []*KpiInfo{
		{
			KpiKey: "123",
		},
		{
			KpiKey: "456",
		},
	}
	num := time.Now().Unix()
	////
	////// 先创建 一
	kpiset = KpiSet{KpiSetDesc: fmt.Sprintf("Abc%d", num)}
	kpiset.KpiInfoID = KpiInfos // 多 绑定到   一
	db.Create(&kpiset)
	//
	//// ========
	//// = xxx 查询 =
	//// ========
	//
	kpiset = KpiSet{KpiSetDesc: fmt.Sprintf("Abc%d", num)}
	//
	//// 先找一
	db.Where(&kpiset).First(&kpiset)
	//
	//// 通过一 找 多 方法A
	var kpiinfoss *[]*KpiInfo = &kpiset.KpiInfoID

	//Db.Preload("KpiInfoID").First(&kpiset)
	db.Preload(clause.Associations).Find(&kpiset)

	//db.Model(&kpiset).Association("KpiInfoId").Find(&kpiset)

	//n := db.Model(&kpiinfoss).Association("KpiInfoId").Count()
	fmt.Printf("%+v\n", kpiset)
	//// xxx 总数
	fmt.Printf("关联数量:%d,%+v\n", len(*kpiinfoss), kpiinfoss)
	time.Sleep(10 * time.Second)
	//
	////xxx  通过 多找一(缺)  方法B  https://gorm.io/zh_CN/docs/preload.html

	//
	//// 预加载分两条查询语句
	////xxx  Preload   通过一  , 返回 一 和  多那边的 数据

	fmt.Println("===============================================")
	kpiset1 := KpiSet{KpiSetDesc: fmt.Sprintf("Abc%d", num)}
	//var kpiInfo112 = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}
	db.Where(&kpiset1).Take(&kpiset1)
	//fmt.Printf("%+v\n", kpiset1)

	var kpiinfoss1 []*KpiInfo = kpiset.KpiInfoID
	db.Model(&kpiset1).Association("KpiInfoID").Find(&kpiinfoss1)
	//db.Where(&kpiInfo112).Find(&kpiInfo113) // 一这边 匹配多个，有多个 Abc1匹配
	//
	//db.Model(&kpiInfo113).
	//	Preload("KpiSetId").
	//	Find(&kpiInfo113)
	//
	fmt.Printf("%+v\n", kpiset1)
	fmt.Printf("%+v\n", *(kpiinfoss1[0]))

	fmt.Println("===============================================")
	//
	////xxx 多 找一
	//kpiset44 := KpiSet{}
	//db.First(&kpiset44, "id = 1005") // 从数据库找一条数据
	//
	////db.Model(&kpiset44).Related(&kpiset44, "kpi_info_id") // Related 返回 一条数据
	//
	//fmt.Printf("关联 一那边 的id:%d \n", kpiset44.KpiInfoID)
	//fmt.Printf("%+v\n", kpiset44)
	//
	////for _, user := range kpiset44.KpiInfoID {
	////	fmt.Print(user.Name + "  ")
	////}
	//fmt.Println("===============================================")
	//
	//// ========
	//// = xxx 更新 =
	//// ========
	//// 使用创建新关联替换当前关联

	//KpiInfos2 := []*KpiInfo{
	//	&KpiInfo{
	//		KpiKey: fmt.Sprintf("Abc%d", num),
	//	},
	//	&KpiInfo{
	//		KpiKey: fmt.Sprintf("Abc%d", num),
	//	},
	//}
	//
	var kpiset3 = KpiSet{}
	db.Where("kpi_set_desc = ?", fmt.Sprintf("Abc%d", num)).First(&kpiset3)

	fmt.Println(kpiset3)

	//err:=db.Model(&kpiset3).Association("KpiInfoID").Replace(KpiInfos2) // xxx 更新 bug

	// xxx 建议 删除 旧的数据
	kpiset3 = KpiSet{}
	db.Where("kpi_set_desc = ?", fmt.Sprintf("Abc%d", num)).First(&kpiset3)
	Db.Preload("KpiInfoID").First(&kpiset3)
	// 删除前
	fmt.Println("删除前")
	fmt.Printf("%+v,%d,%+v,%+v\n", kpiset3.ID, len(kpiset3.KpiInfoID), *(kpiset3.KpiInfoID[0]), *(kpiset3.KpiInfoID[1]))

	// 删除后
	//db.Model(&kpiset3).Association("KpiInfoID").Clear() //清空关联
	//xxx soft del
	//db.Model(&kpiset3).Association("KpiInfoID").Delete(kpiset3.KpiInfoID)
	// time.Sleep(10*time.Second)

	//xxx  hard del
	db.Debug().Unscoped().Delete(kpiset3.KpiInfoID)

	fmt.Println("===============================================")
	kpiset3 = KpiSet{}
	db.Where("kpi_set_desc = ?", fmt.Sprintf("Abc%d", num)).First(&kpiset3)
	Db.Preload("KpiInfoID").First(&kpiset3)
	fmt.Println("删除后")
	fmt.Printf("%+v,%d\n", kpiset3.ID, len(kpiset3.KpiInfoID))

	// kpiset3 = KpiSet{ID: 267}
	////db.Where("kpi_set_desc = ?", fmt.Sprintf("Abc%d", 1637117092)).Take(&kpiset3)
	//Db.Preload("KpiInfoID").First(&kpiset)
	//fmt.Printf("%+v\n", kpiset3)
	//
	//kpiinfoss2:=[]*KpiInfo{}
	//db.Model(&kpiset3).Association("KpiInfoID").Find(&kpiinfoss2)
	// fmt.Printf("%+v\n",kpiinfoss2[0])

	//fmt.Println(kpiset3.CreatedAt.Format(time.RFC3339))
	//
	//// 多个
	//// languages := []*Language{}
	//// db.Where("id IN (?)", []uint{1, 2}).Find(&languages)
	//// db.Model(&user).Association("Languages").Replace(languages)
	//// fmt.Println(user)
	//
	//fmt.Println("===============================================")
	//// ========
	//// = xxx 清除 =
	//// ========
	//// 清空对关联的引用，不会删除关联数据本身
	//var kpiInfo116 = KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}
	//
	//// 删除关联的引用，不会删除关联本身
	//// &kpiInfo116.KpiSetId 普通查询为空，需要关联查询得到
	//db.Where(&kpiInfo116).First(&kpiInfo116)
	//
	////db.Model(&kpiInfo116).Association("KpiSetId").Clear()
	//
	////db.Model(&kpiInfo116).Association("KpiSetId").Delete(&kpiInfo116.KpiSetId)
	//
	//fmt.Printf("%+v,%d\n", kpiInfo116.ID, len(kpiInfo116.KpiSetId))
	//
	//// 清除关联后  使用Related 还是可以查到
	//db.Model(&kpiInfo116).Related(&kpiInfo116.KpiSetId).Find(&kpiInfo116.KpiSetId)
	//fmt.Printf("%+v,%d\n", kpiInfo116.ID, len(kpiInfo116.KpiSetId))
	//
	//// ========
	//// =xxx 删除 =
	//// ========
	//
	//kpiInfo117 := KpiInfo{KpiKey: fmt.Sprintf("Abc%d", num)}
	//
	//// 软删除 https://jasperxu.com/gorm-zh/crud.html#d
	//// 如果模型有DeletedAt字段，它将自动获得软删除功能！ 那么在调用Delete时不会从数据库中永久删除，而是只将字段DeletedAt的值设置为当前时间,否则 真正删除
	//
	//// 软删除时，清空引用后再执行删除，没有Deletet字段,真正删除
	//
	////db.Debug().Delete(&kpiInfo116) // 一 这边删除
	////默认软删除，Unscoped()记录真正删除
	//
	////db.Debug().Unscoped().Delete(&kpiInfo116) // 一 这边删除
	//// fmt.Println(user)
	//
	////fmt.Println(len(kpisetss11))
	////db.Debug().Unscoped().Delete(&kpisetss11) // 一 这边删除,没有
	//
	//// xxx 先查找多，并物理删除
	//db.Where(&kpiInfo117).Find(&kpiInfo117)
	//var kpisetss333 []KpiSet
	//db.Model(&kpiInfo117).Association("KpiSetId").Find(&kpisetss333)
	//
	////db.Debug().Unscoped().Delete(kpisetss333)
	//
	//fmt.Println(len(kpisetss333))
	//for _, value := range kpisetss333 {
	//	db.Debug().Unscoped().Delete(value)
	//}
	//
	//// xxx 物理删除一 这边的数据
	//db.Debug().Select("KpiSetId").Unscoped().Delete(&kpiInfo117)
	////db.Select("KpiSetId").Delete(&kpiInfo117)

}
