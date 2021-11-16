package main

import (
	_ "github.com/mattn/go-oci8"
	//"github.com/go-xorm/xorm"
	"fmt"
	"github.com/xormplus/xorm"
	"log"
	"runtime"
)

//const (
//	INTERVAL = 10
//)

//type TF_OP_VEHICLE  struct {
//	Id         int64    `xorm:"notnull pk autoincr INT(11)",json:"id"` // todo xorm 参数必须写在最前面
//	Uuid       string `xorm:"varchar(64) notnull unique",json:"uuid"`
//	Ip         string `json:"ip"`
//	Passwd     string `json:"-"`
//	Username   string `json:"username"`
//	Port       uint    `xorm:"INT(4)"`
//	Bizid      uint    `gorm:"column:bizid"`
//	Createtime time.Time `xorm:"created"` //在调用方法Insert时自动赋值为当前时间，
//	//UpdatedAt time.Time `xorm:"updated"`    //更新时间 自动
//}

var engine *xorm.Engine //定义引擎全局变量

//func sync()  {
//	// todo 同步结构体与数据表,和python的 migrate 一样同步数据库
//	if err := engine.Sync2(new(Hostinfo)); err != nil {
//		log.Fatalf("Fail to sync database: %v\n", err)
//	}
//	if err := engine.Sync2(new(Hostrun)); err != nil {
//		log.Fatalf("Fail to sync database: %v\n", err)
//	}
//}

// 启动程序后就执行
func main() {
	// 创建 ORM 引擎与数据库
	var err error
	engine, err = xorm.NewEngine("oci8", "test1/test1@1.119.132.155:1521/orcl")
	//db,err:=xorm.NewOracle("oci8","test1/test1@1.119.132.155:1521/orcl")

	if err != nil {
		log.Fatalf("Fail to create hrengine: %v\n", err)
	}

	//sql := "SELECT TRAIN_SERIAL,VEHICLE_SERIAL,VEHICLE_ORDER,VEHICLE_TYPE,POS_AB_FLAG FROM TF_OP_VEHICLE WHERE " +
	//	"pass_time > to_date(?, 'YYYY/MM/DD HH24:MI:SS'); "
	//sql := fmt.Sprintf("select train_serial,train_id,station_id,pass_time,vehicle_number,index_id " +
	//	"from tf_op_train where pass_time > to_date('%s', 'YYYY/MM/DD HH24:MI:SS');","2018/9/25 9:38:51")

	sql := fmt.Sprintf("select train_serial,train_id,station_id,pass_time,vehicle_number,index_id "+
		"from tf_op_train where pass_time > to_date('%s', 'YYYY/MM/DD HH24:MI:SS')", "2018/9/25 9:38:51")
	//sql = "select * from tf_op_train"

	fmt.Println(sql)
	results, err := engine.QueryString(sql)
	r, err := engine.QueryResult(sql).Count()
	fmt.Println(r, err)

	fmt.Printf("%T,%+v\n", results, results)

	//sync()

	// todo 创建日志 可以选用  原生SQL，不能写入日志
	//f, err := os.Create("C:\\work\\go-dev\\src\\godev\\mymodels\\数据库\\xrom\\sql.log")
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//log1:=xorm.NewSimpleLogger(f) //todo 将SQL语句写到日志
	//engine.ShowSQL(true)
	//engine.SetMaxIdleConns(5)
	//log1.Debug("111111")
	//results, err := engine.Query("select * from hostinfo")

	////todo insert
	//u := new(Hostinfo) // 创建一个对象
	//u.Port=22 // 给对象的字段赋值
	//u.Username= "abc"
	//u1, err := engine.InsertOne(u) // 执行插入操作
	//if err != nil {
	//	SimplePanic(err)
	//}
	//fmt.Println(u)
	//fmt.Println(u1)

	////todo selelct
	//fmt.Println("=====================================")
	//results1, err := engine.QueryString("select * from hostinfo")
	////results, err := engine.Where("a = 1").QueryString()
	//fmt.Println(results1)
	//fmt.Println("=====================================")
	//results2, err := engine.QueryInterface("select * from hostinfo")
	////results, err := engine.Where("a = 1").QueryInterface()
	//fmt.Println(results2)
	//fmt.Println("=====================================")
	//
	//// 单条数据
	//h:=Hostinfo{}
	//results, err := engine.Where("id >= ?",0).Or("id >= ?",1).Desc("id").Cols("uuid").Get(&h)
	//fmt.Println(results)
	//fmt.Println(h)

	// 多条数据
	//fmt.Println("=====================================")
	//h12 := make([]Hostinfo, 0)
	//err = engine.Where("port = ? or id = ?", 22, 0).Limit(20, 0).Desc("id").Find(&h12)
	//fmt.Println(h12)
	//SimplePanic(err)

	//fmt.Println("=====================================")
	//var ints []int64
	//err = engine.Table("hostinfo").Cols("id").Find(&ints)
	//SimplePanic(err)
	//fmt.Println(ints)
	//
	//fmt.Println("=====================================")
	////todo 一对多 SQL 或者extend
	////	使用join方法，在表Struct中定义关联字段extend
	///*
	//SELECT hostrun.*
	//FROM hostinfo
	//INNER JOIN hostrun
	//ON hostinfo.id = hostrun.host_id
	//*/
	//var ho []Hostrun
	//err = engine.Table("hostinfo").Select("hostrun.*").
	//	Join("INNER", "hostrun", "hostinfo.id = hostrun.host_id").
	//	Where("hostinfo.id > ?",0 ).Limit(10, 0).
	//	Find(&ho)
	//fmt.Printf("%+v\n",ho)
}

func SimplePanic(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(file, line, err)
		runtime.Goexit()
	}
}
