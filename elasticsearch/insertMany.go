package main

import (
	"context"
	"elasticsearch/g"
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strconv"
)

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(g.Servers...),
		elastic.SetSniff(false), //docker es

	)
	if err != nil {
		panic(err)
	}
	// xxx 方法一
	subjects := []g.Subject{
		g.Subject{
			ID:     1,
			Title:  "肖恩克的救赎",
			Genres: []string{"犯罪", "剧情"},
		},
		g.Subject{
			ID:     2,
			Title:  "千与千寻1",
			Genres: []string{"剧情", "喜剧", "爱情", "战争"},
		},
	}

	bulkRequest := client.Bulk()
	for _, subject := range subjects {
		doc := elastic.NewBulkIndexRequest().Index(g.IndexName).Type(g.TypeName).Id(strconv.Itoa(subject.ID)).Doc(subject)
		bulkRequest = bulkRequest.Add(doc)
	}

	response, err := bulkRequest.Do(ctx)
	if err != nil {
		panic(err)
	}
	failed := response.Failed()
	l := len(failed)
	if l > 0 {
		fmt.Printf("Error(%d),%v", l, response.Errors)
	}

	// xxx 方法二
	subject3 := g.Subject{
		ID:     3,
		Title:  "这个杀手太冷",
		Genres: []string{"剧情", "动作", "犯罪"},
	}
	subject4 := g.Subject{
		ID:     4,
		Title:  "阿甘正传",
		Genres: []string{"剧情", "爱情"},
	}

	subject5 := subject3
	subject5.Title = "这个杀手不太冷"

	index1Req := elastic.NewBulkIndexRequest().Index(g.IndexName).Type(g.TypeName).Id("3").Doc(subject3)
	index2Req := elastic.NewBulkIndexRequest().OpType("create").Type(g.TypeName).Index(g.IndexName).Id("4").Doc(subject4)
	delete1Req := elastic.NewBulkDeleteRequest().Index(g.IndexName).Type(g.TypeName).Id("1")
	update2Req := elastic.NewBulkUpdateRequest().Index(g.IndexName).Type(g.TypeName).Id("3").Doc(subject5)

	bulkRequest = client.Bulk()
	bulkRequest = bulkRequest.Add(index1Req)
	bulkRequest = bulkRequest.Add(index2Req)
	bulkRequest = bulkRequest.Add(delete1Req)
	bulkRequest = bulkRequest.Add(update2Req)

	_, err = bulkRequest.Refresh("wait_for").Do(ctx)
	if err != nil {
		panic(err)
	}

	if bulkRequest.NumberOfActions() == 0 {
		fmt.Println("Actions all clear!")
	}

	searchResult, err := client.Search().
		Index(g.IndexName).
		Sort("id", false). // 按id升序排序
		Pretty(true).
		Do(ctx) // 执行
	if err != nil {
		panic(err)
	}
	var subject g.Subject
	for _, item := range searchResult.Each(reflect.TypeOf(subject)) {
		if t, ok := item.(g.Subject); ok {
			fmt.Printf("Found: Subject(id=%d, title=%s)\n", t.ID, t.Title)
		}
	}

}
