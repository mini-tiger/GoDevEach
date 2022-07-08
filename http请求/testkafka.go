package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  github
 * @Version: 1.0.0
 * @Date: 2021/11/8 下午5:30
 */
//type Repository struct {
//	ID              int        `json:"id"`
//	NodeID          string     `json:"node_id"`
//	Name            string     `json:"name"`
//	FullName        string     `json:"full_name"`
//	Owner           *Developer `json:"owner"`
//	Private         bool       `json:"private"`
//	Description     string     `json:"description"`
//	Fork            bool       `json:"fork"`
//	Language        string     `json:"language"`
//	ForksCount      int        `json:"forks_count"`
//	StargazersCount int        `json:"stargazers_count"`
//	WatchersCount   int        `json:"watchers_count"`
//	OpenIssuesCount int        `json:"open_issues_count"`
//}

func main() {
	client := resty.New()

	var result []map[string]interface{}
	client.R().
		//SetAuthToken("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").
		SetHeader("X-auth-token", "ef98065f-0f78-4164-b36b-bab7ed1dc97f").
		//SetQueryParams(map[string]string{
		//	"per_page":  "3",
		//	"page":      "1",
		//	"sort":      "created",
		//	"direction": "asc",
		//}).
		//SetPathParams(map[string]string{
		//	"org": "posts",
		//}).
		SetResult(&result).
		Get("http://172.22.50.25:30972/topics?pageIndex=1&pageSize=10&clusterId=4a6105e7-ee31-4694-908e-68211e025e2a")

	//fmt.Println(result)
	if len(result) == 0 {
		fmt.Println("err 401")
	}

	var sum float64
	for i, repo := range result {
		if repo["name"].(string) == "zhangdaweitest" {
			fmt.Printf("repo%d: %#v,%v\n", i+1, repo, repo["totalLogSize"].(float64)/1024/1024/1024)
		}

		sum = sum + repo["totalLogSize"].(float64)

	}
	fmt.Println(sum / 1024 / 1024 / 1024)
}
