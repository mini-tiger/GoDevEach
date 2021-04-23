package main

import (
	"fmt"
	"github.com/jackdanger/collectlinks"
	"net/http"
)

func main() {
	resp, _ := http.Get("http://btbtdy3.com/down/35803-0-0.html")
	links := collectlinks.All(resp.Body)
	fmt.Println(links)
}
