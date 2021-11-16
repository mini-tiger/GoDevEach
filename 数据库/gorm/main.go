package main

import (
	"fmt"
	"gormDemo/models"
	"gormDemo/utils"
	"strings"
	"time"
)

type EndpointCounter struct {
	ID         uint `gorm:"primary_key"`
	EndpointID int
	Counter    string
	Step       int
	Type       string
	Ts         int
	TCreate    time.Time
	TModify    time.Time
}

type Endpoint struct {
	ID               uint              `gorm:"primary_key"`
	Endpoint         string            `json:"endpoint"`
	Ts               int               `json:"_"`
	TCreate          time.Time         `json:"-"`
	TModify          time.Time         `json:"-"`
	EndpointCounters []EndpointCounter `gorm:"ForeignKey:EndpointIDE"`
}

func falcon() {
	utils.Init_db()
	db := utils.Conn().Uic

	var u models.User
	db.Where("name=?", "admin").Find(&u)
	fmt.Println(u)

	fmt.Printf("%s\n", "====================================")
	db = utils.Conn().Graph

	var counters []EndpointCounter

	eids := fmt.Sprintf("(%d.%d)", 517, 141) // 需要带括号的格式 in
	dt := db.Table("endpoint_counter").Select("endpoint_id, counter, step, type").Where(fmt.Sprintf("endpoint_id IN %s", eids))
	// Select where 级联调用
	dt = dt.Where("counter regexp ?", strings.TrimSpace("df"))
	//todo 将结果再次where，limit 分页用  ， scan将结果存入 counteres切片，每个元素按照Struct存储
	dt = dt.Limit(50).Offset(0).Scan(&counters)
	for _, v := range counters {
		fmt.Printf("%+v\n", v)
	}

	fmt.Printf("%s\n", "====================================")

	type DBRows struct {
		Endpoint  string
		CounterId int
		Counter   string
		Type      string
		Step      int
	}
	inputs := []string{}
	inputs = append(inputs, []string{"falcon-win12-1", "localhost.localdomain"}...)
	rows := []DBRows{}
	dt = db.Raw(
		`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoint as a, endpoint_counter as b
		where b.endpoint_id = a.id
		AND a.endpoint in (?)`, inputs) // #这里也可以加scan

	dt.Limit(10).Scan(&rows) //limit 限制输出的条数
	for _, v := range rows {
		fmt.Printf("%+v\n", v)
	}

	fmt.Printf("%s\n", "====================================")
	// todo delete
	//del_r := Endpoint{}
	tx := db.Begin() // todo 做update,delete,create时候要 用begin
	dt = tx.Table("endpoint").Where("endpoint in (?)", inputs).Delete(&Endpoint{})
	//fmt.Println(Endpoint{})
	fmt.Println(dt.RowsAffected) //返回删除了几行

	//dt = tx.Exec(`delete from tag_endpoint where endpoint_id in //执行原始 SQL
	//		(select id from endpoint where endpoint in (?))`, inputs)

	tx.Rollback() //回滚
	//tx.Commit() // todo 这行代表提交到数据库
	fmt.Printf("%s\n", "====================================")

	// todo create
	tx = db.Begin() // todo 做update,delete,create时候要 用begin
	c := Endpoint{2,
		"test",
		111,
		time.Now(),
		time.Now(),
		[]EndpointCounter{}}
	dt = tx.Save(c) //还没有提交
	fmt.Println(dt.RowsAffected)
	tx.Commit() //提交

	fmt.Printf("%s\n", "====================================")

	// todo update
	tx = db.Begin() // todo 做update,delete,create时候要 用begin

	e := Endpoint{}                                             //
	db.Table("endpoint").Where("endpoint = ?", "test").Scan(&e) // 先找到需要修改的记录
	//fmt.Println(e)
	id := e.ID //将结果e中主键ID 提取
	//fmt.Println(id)
	update_data := map[string]interface{}{ //需要修改的字段，不写字段不修改
		"endpoint": "test2",
		"ts":       333,
	}
	dt = db.Model(&e).Where("id = ?", id).Update(update_data) //这里不用commit
	fmt.Println(dt.RowsAffected)
}

func main() {
	//falcon()

	//一对多
	utils.Init_db_mulit()
	db := utils.Conn().Nodeman

	var hi models.Hostinfo
	//var hrun []models.Hostrun
	db.Where("id=?", "1").Find(&hi)
	fmt.Println(hi)

	a := db.Model(&hi).Association("Run").Count()
	fmt.Println(hi)
	fmt.Println(a)

	var hrun []models.Hostrun

	hi1 := models.Hostinfo{ID: 1}                 // ID为1的hostinfo 对应多条数据
	db.Model(&hi1).Association("Run").Find(&hrun) // 一个查找多个
	fmt.Println(hrun)

	var hr []models.Hostrun
	db.Table("hostrun").Find(&hr) // 查找hostrun 表里所有记录
	for _, v := range hr {        //通过多这边  找到对应一那边
		tmp := models.Hostinfo{ID: v.Host_id}
		db.Table("hostinfo").Find(&tmp)
		fmt.Printf("当前hostrun 对应的HOSTID：%d\n", v.Host_id)
		fmt.Printf("对应的hostinfo %+v\n", tmp)

	}

	//fmt.Println(hi11)
	//
	//aa:=db.Model(&hi).Where("id=?", "1")
	//hostid:=aa.Value
	//fmt.Println(hostid)
	//db.Where("host_id=?",).Find(&hrun)
	//fmt.Println(hrun)

}
