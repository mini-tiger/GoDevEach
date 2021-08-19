package controllers

import (
	"context"
	funcs "datacenter/funcs"
	"datacenter/g"
	"datacenter/modules"
	"datacenter/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  auth.go
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午3:39
 */

type ResultData struct {
	Pageinfo map[string]interface{}   `json:"pageinfo"`
	List     []map[string]interface{} `json:"list"`
}

func QueryByEs(c *gin.Context) {
	srv := funcs.OAuthSrv

	_, err := srv.ValidationBearerToken(c.Request)
	if err != nil {
		ResponseError(c, err)
		return
	}

	indexName := c.PostForm("indexName")

	var pageIndex, pageSize int
	_, ok := c.GetPostForm("pageIndex")
	if !ok {
		pageIndex = 1
	} else {
		pageIndex, _ = strconv.Atoi(c.PostForm("pageIndex"))
	}
	//queryFields := c.PostForm("queryFields")

	_, ok = c.GetPostForm("pageSize")
	if !ok {
		pageSize = 10
	} else {
		pageSize, _ = strconv.Atoi(c.PostForm("pageSize"))
	}

	// 构建查询语句
	queryCondition, ok := c.GetPostForm("query")
	queryMap := make(map[string]interface{})
	if ok {

		err = json.Unmarshal([]byte(queryCondition), &queryMap)
		if err != nil {
			ResponseError(c, err)
			return
		}
	}

	//fmt.Println(queryCondition)
	//fmt.Println(queryMap)
	//fmt.Printf(indexName)
	var queryKeyLen int = len(queryMap)

	//if len(queryMap) > 1 {
	//	fmt.Println(1111)
	//}

	//var queryCond *elastic.BoolQuery=elastic.NewBoolQuery()

	searchFusion := funcs.EsClient.Search().Index(indexName)

	searchSource := elastic.NewSearchSource()

	// page
	searchSource.From((pageIndex - 1) * pageSize).Size(pageSize)
	searchFusion.SearchSource(searchSource)

	// querycrond
	switch true {
	case queryKeyLen == 1:
		if _, or_ok := queryMap["or"]; or_ok {

			boolQuery := elastic.NewBoolQuery()
			for key, value := range queryMap["or"].(map[string]interface{}) {
				termQuery := elastic.NewMatchQuery(key, value)
				//fmt.Println(key, value)
				boolQuery.Should(termQuery)
			}
			searchFusion.Query(boolQuery)
			searchSource.Query(boolQuery)

		} else {
			for key, v := range queryMap {

				termQuery := elastic.NewMatchQuery(key, v)
				//fmt.Println(key, value)
				//queryCond.Should(termQuery)
				//searchFusion.Query(termQuery)
				searchSource.Query(termQuery)
			}

		}

		break
	case queryKeyLen > 1:
		boolQuery := elastic.NewBoolQuery()
		for key, value := range queryMap {
			termQuery := elastic.NewMatchQuery(key, value)
			//fmt.Println(key, value)
			boolQuery.Must(termQuery)
		}
		//searchFusion.Query(boolQuery)
		searchSource.Query(boolQuery)

	}

	// sort
	sortBy, sortbyok := c.GetPostForm("sortBy")

	if sortbyok {
		sortQuery := elastic.NewFieldSort(sortBy).Asc()
		if c.PostForm("sortDirection") == "-1" {
			sortQuery = sortQuery.Desc()
		}
		searchSource.SortBy(sortQuery)
		searchFusion.SortBy(sortQuery)
	}

	// show fields
	fields, ok := c.GetPostForm("queryFields")
	//fmt.Println(fields)
	if ok {

		fsc := elastic.NewFetchSourceContext(true).Include(strings.Split(fields, ",")...)
		searchFusion.FetchSourceContext(fsc)
		searchSource.FetchSourceContext(fsc)
	}

	//termQuery := elastic.NewTermQuery("park_name", "实测交大停车场") // 不会对搜索词进行分词处理，而是作为一个整体与目标字段进行匹配

	//d := elastic.NewMatchQuery("park_name", "交大科技大厦") // 会将搜索词分词

	//andbool:=elastic.NewBoolQuery().Must(termQuery,d) // and

	//orbool := elastic.NewBoolQuery().Should(termQuery, d) // or

	//xxx andbool:=elastic.NewBoolQuery().Must(termQuery,d) // and

	//xxx orbool := elastic.NewBoolQuery().Should(termQuery, d) // or

	// xxx 多条件方法一
	//timeQ := elastic.NewRangeQuery("@timestamp").From(from).To(End)
	//componentQ := elastic.NewTermQuery("component", *component)
	//deploymentQ := elastic.NewTermQuery("deploymentName", deploymentName)
	//
	//generalQ := elastic.NewBoolQuery()
	//generalQ = generalQ.Must(timeQ).Must(componentQ).Must(deploymentQ)

	//xxx 多条件方法二
	// 创建bool查询
	//boolQuery := elastic.NewBoolQuery().Must()
	//
	//// 创建term查询
	//termQuery := elastic.NewTermQuery("Author", "tizi")
	//matchQuery := elastic.NewMatchQuery("Title", "golang es教程")
	//
	//// 设置bool查询的should条件, 组合了两个子查询
	//// 表示搜索Author=tizi或者Title匹配"golang es教程"的文档
	//boolQuery.Should(termQuery, matchQuery)

	// xxx https://www.coder.work/article/1023959
	//

	//searchResult, err := funcs.EsClient.Search().
	//	Index(indexName).
	//	//Type(typeName).
	//	Query(termQuery).
	//	Sort("update_date", true). // 按id升序排序
	//	From(0).Size(10). // 拿前10个结果
	//	Pretty(true).
	//	FetchSourceContext(fsc).
	//	Do(context.Background()) // 执行

	//searchFusion := funcs.EsClient.Search().Index(indexName)
	//searchFusion:=elastic.NewSearchService(funcs.EsClient)

	//searchFusion.Index(indexName)
	//searchFusion.Source(searchSource)

	searchResult, err := searchFusion.Do(context.Background())

	if err != nil {
		panic(err)
	}
	total := searchResult.TotalHits()
	fmt.Printf("Found %d subjects\n", total)

	var result []map[string]interface{}
	if total > 0 {
		for _, hit := range searchResult.Hits.Hits {
			item := make(map[string]interface{})

			err := json.Unmarshal(hit.Source, &item)
			result = append(result, item)
			if err != nil {
				fmt.Printf("err:%v\n", err)
				continue
			}
			//fmt.Printf("doc %+v\n", item)
		}
	} else {
		ResponseError(c, errors.New("Not Found Data!"))
		return
		//fmt.Println("Not found !")
	}

	//pageinfo:=make(map[string]interface{})
	//pageinfo["pageIndex"]=pageIndex
	//pageinfo["pageSize"]=pageSize
	//pageinfo["totalCount"]=total

	body, _ := searchSource.Source()
	mjson, _ := json.MarshalIndent(body, "", "\t")
	//fmt.Println("return data:",len(result))
	g.GetLog().Debug("QueryByEs query :%+v\n", string(mjson))

	ResponseSuccess(c, ResultData{Pageinfo: map[string]interface{}{
		"pageIndex": pageIndex, "pageSize": pageSize, "totalCount": total},
		List: result})

}

func Token(c *gin.Context) {

	s := funcs.OAuthSrv
	gt, tgr, err := s.ValidationTokenRequest(c.Request)
	if err != nil {
		ResponseError(c, err)
		return
	}
	tgr.ClientSecret = utils.Md5V3(tgr.ClientSecret)
	//fmt.Println(tgr)
	ti, err := s.GetAccessToken(context.Background(), gt, tgr)

	if err != nil {
		ResponseError(c, err)
		return
	}

	if err != nil {
		ResponseError(c, err)
		return
	}
	//ResponseSuccess(c,	s.GetTokenData(ti))
	c.JSON(http.StatusOK, s.GetTokenData(ti))

}

func InitClient(c *gin.Context) {

	secret := c.PostForm("secret")
	resource_ids := c.PostForm("resource_ids")
	scope := c.PostForm("scope")
	authorized_grant_types := c.PostForm("authorized_grant_types")
	web_server_redirect_uri := c.PostForm("web_server_redirect_uri")
	authorities := c.PostForm("authorities")
	access_token_validity := c.PostForm("access_token_validity")
	refresh_token_validity := c.PostForm("refresh_token_validity")
	additional_information := c.PostForm("additional_information")
	autoapprove := c.PostForm("autoapprove")
	client_id := c.PostForm("client_id")

	var db *gorm.DB
	// 查询存在
	var clientDetail modules.OauthClientDetails
	db = modules.MysqlDb.Table("oauth_client_details").Select("*").Where("client_id=?", client_id).First(&clientDetail)

	var u modules.User // 准备用户数据
	u.Name = client_id
	u.Password = utils.Md5V3(secret)
	u.Source = "client_users"
	u.Roles = "admin"
	u.Status = "1"

	// 1. 存在则创建 项目 默认用户
	if db.RowsAffected > 0 {
		modules.MysqlDb.Table("user").Select("id").Where(" name=? and  source='client_users'", client_id).First(&u)
		if u.Id == 0 {

			db = modules.MysqlDb.Table("user").Create(&u)
			//fmt.Printf("%+v\n", u)
			//fmt.Printf("%+v\n", db)
			if db.RowsAffected == 0 {
				ResponseError(c, errors.New("user table insert Fail"))
				return
			}
		}
	} else {
		clientDetail.ClientSecret = utils.Md5V3(secret)
		clientDetail.ClientId = client_id
		clientDetail.WebServerRedirectUri = web_server_redirect_uri
		atv, _ := strconv.Atoi(access_token_validity)

		clientDetail.AccessTokenValidity = atv
		clientDetail.AdditionalInformation = additional_information
		clientDetail.Authorities = authorities
		clientDetail.AuthorizedGrantTypes = authorized_grant_types
		clientDetail.Autoapprove = autoapprove
		rtv, _ := strconv.Atoi(refresh_token_validity)
		clientDetail.RefreshTokenValidity = rtv
		clientDetail.ResourceIds = resource_ids
		clientDetail.Scope = scope
		db = modules.MysqlDb.Table("oauth_client_details").Create(&clientDetail)
		//fmt.Printf("%+v\n", u)
		//fmt.Printf("%+v\n", db)
		if db.RowsAffected == 0 {
			ResponseError(c, errors.New("oauth_client_details table insert Fail"))
			return
		} else {
			db = modules.MysqlDb.Table("user").Create(&u)
			//fmt.Printf("%+v\n", u)
			//fmt.Printf("%+v\n", db)
			if db.RowsAffected == 0 {
				ResponseError(c, errors.New("user table insert Fail"))
				return
			}
		}
	}
	ResponseSuccess(c, clientDetail)
	return

}
func DeleteClient(c *gin.Context) {
	//fmt.Println(c.PostForm("client_secret"))
	secret := utils.Md5V3(c.PostForm("client_secret"))
	client_id := c.PostForm("client_id")
	//fmt.Println(client_id)
	//fmt.Println(secret)
	var db *gorm.DB
	clientDetailDB := modules.MysqlDb.Table("oauth_client_details")
	userDB := modules.MysqlDb.Table("user")
	//fmt.Println(secret)

	// 查询存在
	var clientDetail modules.OauthClientDetails
	db = clientDetailDB.Select("*").Where("client_id=? and client_secret=?", client_id, secret).First(&clientDetail)
	if db.RowsAffected > 0 {
		clientDetailDB.Delete(&clientDetail)

		userDB.Where("name=? and source='client_users'", client_id).Delete(&modules.User{})
		//modules.MysqlDb.Raw("DELETE from user  where name=? and source='client_users'", client_id).Rows()
		ResponseSuccess(c, "success")
		return
	}
	ResponseError(c, errors.New("Client not Exist"))

}

func UpdateClient(c *gin.Context) {

	secret := utils.Md5V3(c.PostForm("secret"))
	resource_ids := c.PostForm("resource_ids")
	scope := c.PostForm("scope")
	authorized_grant_types := c.PostForm("authorized_grant_types")
	web_server_redirect_uri := c.PostForm("web_server_redirect_uri")
	authorities := c.PostForm("authorities")
	access_token_validity := c.PostForm("access_token_validity")
	refresh_token_validity := c.PostForm("refresh_token_validity")
	additional_information := c.PostForm("additional_information")
	autoapprove := c.PostForm("autoapprove")
	client_id := c.PostForm("client_id")

	var db *gorm.DB
	clientDetailDB := modules.MysqlDb.Table("oauth_client_details")
	userDB := modules.MysqlDb.Table("user")

	// 查询存在
	var clientDetail modules.OauthClientDetails
	db = clientDetailDB.Select("*").Where("client_id=?", client_id).First(&clientDetail)

	var u modules.User // 准备用户数据

	//fmt.Println(clientDetail)
	// 1.
	if db.RowsAffected > 0 {
		//fmt.Println(clientDetail)
		clientDetail.ClientSecret = secret
		clientDetail.ClientId = client_id
		clientDetail.WebServerRedirectUri = web_server_redirect_uri
		atv, _ := strconv.Atoi(access_token_validity)

		clientDetail.AccessTokenValidity = atv
		clientDetail.AdditionalInformation = additional_information
		clientDetail.Authorities = authorities
		clientDetail.AuthorizedGrantTypes = authorized_grant_types
		clientDetail.Autoapprove = autoapprove
		rtv, _ := strconv.Atoi(refresh_token_validity)
		clientDetail.RefreshTokenValidity = rtv
		clientDetail.ResourceIds = resource_ids
		clientDetail.Scope = scope

		//fmt.Println(clientDetail)
		db = clientDetailDB.Save(&clientDetail)

		if db.RowsAffected > 0 { //没有改动 这里是0
			// 查找默认用户 没有则创建
			db = userDB.Select("id").Where(" name=? and  source='client_users'", client_id).First(&u)

			if db.RowsAffected == 0 {
				u.Name = client_id
				u.Password = secret
				u.Source = "client_users"
				u.Roles = "admin"
				u.Status = "1"
				db = modules.MysqlDb.Table("user").Create(&u)
				if db.RowsAffected == 0 {
					ResponseError(c, errors.New("user table insert Fail"))
					return
				}
			} else {
				fmt.Printf("%+v\n", u)
				db.Model(&u).Select("name", "password").Updates(modules.User{Name: client_id, Password: secret})
				fmt.Printf("%+v\n", db)
			}
		}
	} else {
		ResponseError(c, errors.New("clientDetail Not Found"))
	}
	ResponseSuccess(c, clientDetail)
	return

}

func Register(c *gin.Context) {

	username := c.PostForm("username")
	//fmt.Println(username)

	source := c.PostForm("source")
	client := c.PostForm("client")
	password := c.PostForm("password")
	role := c.PostForm("role")
	//fmt.Println(source)

	var u modules.User
	//rows, err := modules.MysqlDb.Raw("SELECT  *  from user where name=? and  source=?",[]string{"name","111"}...).Scan(&u)
	modules.MysqlDb.Table("user").Select("*").Where(" name=? and  source=?", username, source).Find(&u)

	if u.Id == 0 {
		u.Name = username
		u.Client = client
		u.Password = utils.Md5V3(password)
		u.Source = source
		u.Roles = role
		db := modules.MysqlDb.Table("user").Create(&u)
		//fmt.Printf("%+v\n", u)
		//fmt.Printf("%+v\n", db)
		if db.RowsAffected == 0 {
			ResponseError(c, errors.New("user table insert Fail"))
			return
		}
		g.GetLog().Printf("新建用户 %+v\n", u)
	}

	ResponseSuccess(c, u)
	return

}

//func InsertClientStoreUser(clientid, userid string) (err error) {
//	var clientdetail modules.OauthClientDetails
//	modules.MysqlDb.Table("oauth_client_details").Select("*").Where(" client_id=? ", clientid).Find(&clientdetail)
//	//fmt.Println(clientdetail)
//	if clientdetail.ClientSecret == "" {
//		err = errors.New("oauth_client_details table select Fail")
//		return
//	}
//
//	// clientStore insert
//	err = funcs.ClientStore.Set(userid, &models.Client{
//		ID:     clientdetail.ClientId,
//		Secret: clientdetail.ClientSecret,
//		Domain: clientdetail.WebServerRedirectUri,
//	})
//	return err
//}
