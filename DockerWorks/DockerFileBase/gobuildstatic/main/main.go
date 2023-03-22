package main

import (
	"build_dir/submain"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func main() {
	fmt.Println(submain.Sub())
	serveHTTP()
}

func httpTest() {
	reader := strings.NewReader(`{"title":"The Go Standard Library","content":"It contains many packages."}`)
	r, _ := http.NewRequest(http.MethodPost, "/topic/", reader)

	w := httptest.NewRecorder()
	httpHandler(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response code is %v\n", resp.StatusCode)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response body: %v\n", string(data))
	}
}

func serveHTTP() {
	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 100; i++ {
			time.Sleep(1 * time.Second)
			httpTest()
		}
	}()

	http.HandleFunc("/", httpHandler)

	http.ListenAndServe("127.0.0.1:8080", nil)
	select {}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from %s: %s %q", r.RemoteAddr, r.Method, r.URL)
	fmt.Fprintf(w, "%s go-daemon: %q", time.Now().String(), html.EscapeString(r.URL.Path))
}
