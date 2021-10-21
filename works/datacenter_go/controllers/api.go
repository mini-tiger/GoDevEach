package controllers

import (
	"context"
	funcs "datacenter/funcs"
	"datacenter/g"
	"datacenter/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/olivere/elastic/v7"
	"net/url"
	"strconv"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  api
 * @Version: 1.0.0
 * @Date: 2021/9/13 下午5:24
 */

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

	if pageIndex <= 0 || pageSize <= 0 {
		ResponseError(c, errors.New("page params err"))
		return
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

	// xxx querycrond
	//switch true {
	//case queryKeyLen == 1:
	//	if _, or_ok := queryMap["or"]; or_ok {
	//
	//		boolQuery := elastic.NewBoolQuery()
	//		for key, value := range queryMap["or"].(map[string]interface{}) {
	//			termQuery := elastic.NewMatchQuery(key, value)
	//			//fmt.Println(key, value)
	//			boolQuery.Should(termQuery)
	//		}
	//		//searchFusion.Query(boolQuery)
	//		searchSource.Query(boolQuery)
	//	} else {
	//		boolQuery := elastic.NewBoolQuery()
	//		for key, value := range queryMap {
	//			termQuery := elastic.NewMatchQuery(key, value)
	//			//fmt.Println(key, value)
	//			boolQuery.Must(termQuery)
	//		}
	//		//searchFusion.Query(boolQuery)
	//		searchSource.Query(boolQuery)
	//	}
	//
	//	break
	//case queryKeyLen > 1:
	//	boolQuery := elastic.NewBoolQuery()
	//	for key, value := range queryMap {
	//		termQuery := elastic.NewMatchQuery(key, value)
	//		//fmt.Println(key, value)
	//		boolQuery.Must(termQuery)
	//	}
	//	//searchFusion.Query(boolQuery)
	//	searchSource.Query(boolQuery)
	//}

	if _, or_ok := queryMap["or"]; or_ok && queryKeyLen == 1 {
		boolQuery := elastic.NewBoolQuery()
		for key, value := range queryMap["or"].(map[string]interface{}) {
			termQuery := elastic.NewMatchQuery(key, value)
			//fmt.Println(key, value)
			boolQuery.Should(termQuery)
		}
		//searchFusion.Query(boolQuery)
		searchSource.Query(boolQuery)
	} else {
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
		//searchFusion.SortBy(sortQuery)
	}

	// show fields
	fields, ok := c.GetPostForm("queryFields")
	//fmt.Println(fields)
	if ok {

		fsc := elastic.NewFetchSourceContext(true).Include(strings.Split(fields, ",")...)
		//searchFusion.FetchSourceContext(fsc)
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
		ResponseError(c, err)
		return
	}
	total := searchResult.TotalHits()
	//fmt.Printf("Found %d subjects\n", total)

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
	g.GetLog().Debug("IndexName:%s , QueryByEs query :%+v\n", indexName, string(mjson))

	ResponseSuccess(c, &ResultData{PageInfo: map[string]interface{}{
		"pageIndex": pageIndex, "pageSize": pageSize, "totalCount": total},
		List: result})

}

func QueryByEsPlus(c *gin.Context) {
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

	_, ok = c.GetPostForm("pageSize")
	if !ok {
		pageSize = 10
	} else {
		pageSize, _ = strconv.Atoi(c.PostForm("pageSize"))
	}

	if pageIndex <= 0 || pageSize <= 0 {
		ResponseError(c, errors.New("page params err"))
		return
	}
	// 构建查询语句
	queryCondition, ok := c.GetPostForm("query")
	queryMap := make(map[string]string)
	if ok {

		err = json.Unmarshal([]byte(queryCondition), &queryMap)
		if err != nil {
			ResponseError(c, err)
			return
		}
	}

	searchFusion := funcs.EsClient.Search().Index(indexName)

	searchSource := elastic.NewSearchSource()

	// page
	searchSource.From((pageIndex - 1) * pageSize).Size(pageSize)
	searchFusion.SearchSource(searchSource)

	//
	//var queryKeyLen int = len(queryMap)
	//
	//switch true {
	//case queryKeyLen == 1:
	//		for key, v := range queryMap {
	//			termQuery := elastic.NewWildcardQuery(key, v)
	//			searchSource.Query(termQuery)
	//		}
	//	break
	//case queryKeyLen > 1:
	//	boolQuery := elastic.NewBoolQuery()
	//	for key, value := range queryMap {
	//		termQuery := elastic.NewWildcardQuery(key, value)
	//		//fmt.Println(key, value)
	//		boolQuery.Must(termQuery)
	//	}
	//	searchSource.Query(boolQuery)
	//
	//}
	//

	//fmt.Println(regexTpl)

	queryType := c.DefaultPostForm("querytype", "term")

	boolQuery := elastic.NewBoolQuery()
	for key, value := range queryMap {
		var termQuery elastic.Query
		switch true {
		case queryType == "term":
			termQuery = elastic.NewTermQuery(key, value)
		case queryType == "regex":
			termQuery = elastic.NewWildcardQuery(key, fmt.Sprintf(utils.FixReg(c), value))
		default:
			termQuery = elastic.NewTermQuery(key, value)
		}

		//fmt.Println(key, value)
		boolQuery.Must(termQuery)
	}
	searchSource.Query(boolQuery)

	//searchFusion.Query(boolQuery)

	// sort
	sortBy, sortByOk := c.GetPostForm("sortBy")

	if sortByOk {
		sortQuery := elastic.NewFieldSort(sortBy).Asc()
		if c.PostForm("sortDirection") == "-1" {
			sortQuery = sortQuery.Desc()
		}
		searchSource.SortBy(sortQuery)
		//searchFusion.SortBy(sortQuery)
	}

	// show fields
	fields, ok := c.GetPostForm("queryFields")
	//fmt.Println(fields)
	if ok {

		fsc := elastic.NewFetchSourceContext(true).Include(strings.Split(fields, ",")...)
		//searchFusion.FetchSourceContext(fsc)
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
		ResponseError(c, err)
		return
	}
	total := searchResult.TotalHits()
	//fmt.Printf("Found %d subjects\n", total)

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
	g.GetLog().Debug("IndexName:%s , QueryByEs query :%+v\n", indexName, string(mjson))

	ResponseSuccess(c, &ResultData{PageInfo: map[string]interface{}{
		"pageIndex": pageIndex, "pageSize": pageSize, "totalCount": total},
		List: result})

}

//type EsUrlEntry struct {
//	Data interface{} `json:"data"`
//	Path string      `json:"path"`
//}

func QueryBySourceEs(c *gin.Context) {

	//bytes, err := c.GetRawData() // 接收json数据
	//if err != nil {
	//	ResponseError(c, err) //统一返回
	//	return
	//}
	//
	//
	//

	srv := funcs.OAuthSrv

	_, err := srv.ValidationBearerToken(c.Request)
	if err != nil {
		ResponseError(c, err)
		return
	}

	var eue map[string]interface{}

	if err := c.ShouldBindWith(&eue, binding.JSON); err != nil {
		ResponseError(c, err) //统一返回
		return
	}

	//fmt.Printf("%+v\n", eue)
	res, err := funcs.EsClient.PerformRequest(context.Background(), elastic.PerformRequestOptions{
		Method:          "POST",
		Path:            eue["path"].(string),
		Params:          url.Values{},
		Body:            eue["data"],
		Headers:         nil,
		MaxResponseSize: 0,
	})
	if err != nil {
		ResponseError(c, err)
		return
	}

	if g.GetConfig().IsDebug() {
		mjson, _ := json.MarshalIndent(eue["data"], "", "\t")
		g.GetLog().Debug("IndexName:%s , QueryBySourceEs query :%+v\n", eue["path"], string(mjson))
	}

	ResponseSuccess(c, res.Body)
	return
}
