package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/uptrace/bunrouter"
)

func main() {

	router := bunrouter.New()

	router.WithGroup("/api/", func(group *bunrouter.Group) {

		group.GET("/index", index)
		group.GET("/hello", hello)

	})

	router.WithGroup("/api2/", func(group *bunrouter.Group) {

		group.GET("/index", index2)
		group.GET("/hello", hello2)

	})

	log.Println("listening on http://localhost:8080")
	log.Println(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, req bunrouter.Request) error {

	fmt.Fprintln(w, "index")
	return nil
}

func hello(w http.ResponseWriter, req bunrouter.Request) error {

	fmt.Fprintln(w, "hello")
	return nil
}

func index2(w http.ResponseWriter, req bunrouter.Request) error {

	fmt.Fprintln(w, "index2")
	return nil
}

func hello2(w http.ResponseWriter, req bunrouter.Request) error {

	fmt.Fprintln(w, "hello2")
	return nil
}
