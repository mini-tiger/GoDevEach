package main

import (
	"fmt"
	"log"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Description: Air热重载
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/4/28 上午9:19
 */
// xxx https://github.com/cosmtrek/air
/* xxx 使用
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
cp ./bin/air /bin/air

air 替换 go run main.go
*/
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world1!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
