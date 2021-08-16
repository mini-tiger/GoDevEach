package main

import (
	"context"
	"elasticsearch/g"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  33
 * @Version: 1.0.0
 * @Date: 2021/8/14 上午11:30
 */

var (
	subject   Subject
	indexName = "parking_park_order"
	//typeName  = g.TypeName
	servers = g.Servers
)

type Subject struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

func Search(client *elastic.Client, ctx context.Context) {

	// Term搜索

	termQuery := elastic.NewTermQuery("park_name", "实测交大停车场") // 不会对搜索词进行分词处理，而是作为一个整体与目标字段进行匹配

	d := elastic.NewMatchQuery("park_name", "交大科技大厦") // 会将搜索词分词

	//andbool:=elastic.NewBoolQuery().Must(termQuery,d) // and

	orbool := elastic.NewBoolQuery().Should(termQuery, d) // or

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

	searchResult, err := client.Search().
		Index(indexName).
		//Type(typeName).
		Query(orbool).
		Sort("update_date", true). // 按id升序排序
		From(0).Size(10).          // 拿前10个结果
		Pretty(true).
		Do(ctx) // 执行
	if err != nil {
		panic(err)
	}
	total := searchResult.TotalHits()
	fmt.Printf("Found %d subjects\n", total)

	if total > 0 {
		//for _, item := range searchResult.Each(reflect.TypeOf(subject)) {
		//	if t, ok := item.(Subject); ok {
		//		fmt.Printf("Found: Subject(id=%d, title=%s)\n", t.ID, t.Title)
		//	}
		//}

		for _, hit := range searchResult.Hits.Hits {
			item := make(map[string]interface{})
			err := json.Unmarshal(hit.Source, &item)
			if err != nil {
				fmt.Printf("err:%v\n", err)
				continue
			}
			fmt.Printf("doc %+v\n", item)
		}
	} else {
		fmt.Println("Not found!")
	}
}

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(servers...),
		elastic.SetSniff(false), //docker es

	)
	if err != nil {
		panic(err)
	}
	Search(client, ctx)
}
