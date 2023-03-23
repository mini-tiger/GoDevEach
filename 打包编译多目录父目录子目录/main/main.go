package main

import (
	"build_dir/public"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	fmt.Println(public.Util1())
	serveHTTP()
}
func serveHTTP() {
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from %s: %s %q", r.RemoteAddr, r.Method, r.URL)
	fmt.Fprintf(w, "go-daemon: %q", html.EscapeString(r.URL.Path))
}
