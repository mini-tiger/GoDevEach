package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

//https://github.com/gocolly/colly/tree/master/_examples

// https://www.ulovecode.com/2020/04/28/Go/Colly%E5%AE%98%E6%96%B9%E6%96%87%E6%A1%A3%E5%AD%A6%E4%B9%A0%E4%BB%8E%E5%85%A5%E9%97%A8%E5%88%B0%E5%85%A5%E5%9C%9F/

func main() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*.douban.*", Parallelism: 5})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML(".hd", func(e *colly.HTMLElement) {
		log.Println(strings.Split(e.ChildAttr("a", "href"), "/")[4],
			strings.TrimSpace(e.DOM.Find("span.title").Eq(0).Text()))
	})

	c.OnHTML(".paginator a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.Visit("https://movie.douban.com/top250?start=0&filter=")
	c.Wait()
}
