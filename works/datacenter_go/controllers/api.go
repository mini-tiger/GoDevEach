package controllers

import (
	"datacenter/modules"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"path"
	"runtime"
)

func SingleUpLoad(c *gin.Context) {

	_, filename, _, _ := runtime.Caller(0)
	//fmt.Println(filename)
	dir := path.Dir(path.Dir(filename))
	file, _ := c.FormFile("file")
	//log.Println(file.Filename)
	//log.Println(path.Join(dir, file.Filename))
	if err := c.SaveUploadedFile(file, path.Join(dir, file.Filename)); err != nil {
		ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"status":     http.StatusOK,
		"statusText": "ok",

		"data": fmt.Sprintf("upload file:%s success", file.Filename),
	})
}

func MysqlFind(c *gin.Context) {
	bytes, err := c.GetRawData() // 接收json数据
	if err != nil {
		ResponseError(c, err) //统一返回
		return
	}
	data := string(bytes)
	fmt.Println("request params:", data)

	ResponseSuccess(c, nil)

	rows, err := modules.MysqlDb.Raw("SELECT  *  from greed_track_chang_info_in_1 where ?", "1=1").Rows()
	if err != nil {
		ResponseError(c, err) //统一返回
		return
	}
	//a:=modules.DoQuerySort(rows)

	b := modules.DoQuery(rows) //xxx  返回map[string]string,解析字段名
	//if err != nil {
	//	ResponseError(c, err) //统一返回
	//	return
	//}
	fmt.Println(b)
	c.JSON(http.StatusOK, gin.H{

		"status":     http.StatusOK,
		"statusText": "ok",

		"data": b,
	})

}

func PostSpeed1(c *gin.Context) {
	fmt.Println(c.Request.Body)
	fmt.Println(c.Request.Header)
	bytes, err := c.GetRawData() // 接收json数据
	if err != nil {
		ResponseError(c, err) //统一返回
		return
	}
	//data := string(bytes)
	fmt.Println(string(bytes))

	if !gjson.GetBytes(bytes, "user").Exists() || !gjson.GetBytes(bytes, "password").Exists() {
		ResponseError(c, errors.New("params err")) //统一返回
		return
	}

	var mapResult map[string]interface{} // json To map
	err = json.Unmarshal(bytes, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	u := gjson.GetBytes(bytes, "user")
	if u.String() == "taojun" {
		c.JSON(http.StatusOK, gin.H{
			"status":        http.StatusOK,
			"statusText":    "ok",
			"requestParams": mapResult, // response  reqParams
		})
		return
	}
	ResponseError(c, errors.New("user not taojun")) //统一返回
}

func PostSpeed2(c *gin.Context) {

	var login modules.Login

	if errA := c.ShouldBindWith(&login, binding.JSON); errA != nil {
		ResponseError(c, errors.New("params err")) //统一返回
	}
	if login.User == "taojun" {
		c.JSON(http.StatusOK, gin.H{
			"status":     http.StatusOK,
			"statusText": "ok",
		})
		return
	}
	ResponseError(c, errors.New("user not taojun")) //统一返回
}

func GetHisMissionDetail(c *gin.Context) {
	bytes, err := c.GetRawData() // 接收json数据
	if err != nil {
		ResponseError(c, err) //统一返回
		return
	}
	data := string(bytes)

	//fmt.Println(data) //打印json 数据

	if !gjson.Get(data, "demand").Exists() || !gjson.Get(data, "data.taskid").Exists() {
		ResponseError(c, err) //统一返回
		return
	}
	// 指定获取要操作的数据集
	mongoClient, err := modules.NewMongoConn()
	if err != nil {
		ResponseError(c, err) //统一返回
		return
	}
	//fmt.Println(mongoClient.CollectionCount("server_auto","dms-content"))

	findOptions := options.Find()
	//findOptions.SetLimit(2)

	var TableName string
	//let mongotasktablename = ""

	switch gjson.Get(data, "demand").String() {
	case "ship":
		// mongotasktablename = "schedulingtask_ship"
		TableName = "chuanboshengchan" //船舶计划,  data 每条17列
		// mongoOutTableName = mysqlOutTableName
		break
	case "site":
		// mongotasktablename = "schedulingtask_site"
		TableName = "site_plan" //场地计划,  data 每条17列
		// mongoOutTableName = mysqlOutTableName
		break
	}
	var filter = make(map[string]interface{})

	filter = map[string]interface{}{"schemaName": TableName, "data.taskid.iv": gjson.Get(data, "data.taskid").Int()}

	// xxx mongodb 的筛选 和 输出 都可以是 bson.M  bson.D 等

	result, err := mongoClient.CollectionFilter("server_auto", "dms-content", bson.M(filter), findOptions)

	//result, err := mongoClient.CollectionFilter("server_auto", "dms-content",bson.D{{
	//	"data.taskid.iv",
	//	bson.D{{
	//		"$in",
	//		bson.A{"Alice", "Bob"},
	//	}},
	//}}, findOptions)

	if err != nil {
		ResponseError(c, err) //统一返回
		return
	}
	//fmt.Println(result[0])
	mongoClient.DisableConn()

	c.JSON(http.StatusOK, gin.H{

		"status":     http.StatusOK,
		"statusText": "ok",

		"data": result,
	})
}
