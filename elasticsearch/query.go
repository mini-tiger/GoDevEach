package main

import (
	"context"
	"elasticsearch/g"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"reflect"
)

func PrintQuery(src interface{}) {
	//fmt.Println("*****")
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func SearchFunc(client *elastic.Client, ctx context.Context, termQuery *elastic.TermQuery) {

	searchResult, err := client.Search().
		Index(indexName).
		Type(typeName).
		Query(termQuery).
		Sort("id", true). // 按id升序排序
		From(0).Size(10). // 拿前10个结果
		Pretty(true).
		Do(ctx) // 执行
	if err != nil {
		panic(err)
	}
	total := searchResult.TotalHits()
	fmt.Printf("Found %d subjects\n", total)
	if total > 0 {
		for _, item := range searchResult.Each(reflect.TypeOf(subject)) {
			if t, ok := item.(Subject); ok {
				fmt.Printf("Found: Subject(id=%d, title=%s)\n", t.ID, t.Title)
			}
		}

	} else {
		fmt.Println("Not found!")
	}
}

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(g.Servers...))
	if err != nil {
		panic(err)
	}

	query := elastic.NewTermQuery("genres", "动画")
	src, err := query.Source()
	if err != nil {
		panic(err)
	}
	PrintQuery(src)
	SearchFunc(client, ctx, query)

	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(elastic.NewTermQuery("genres", "剧情"))
	boolQuery = boolQuery.Filter(elastic.NewTermQuery("id", 1))
	src, err = boolQuery.Source()
	if err != nil {
		panic(err)
	}
	PrintQuery(src)

	rangeQuery := elastic.NewRangeQuery("born").
		Gte("2012/01/01").
		Lte("now").
		Format("yyyy/MM/dd")
	src, err = rangeQuery.Source()
	PrintQuery(src)
}
