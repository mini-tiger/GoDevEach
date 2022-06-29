package main

import "net/http"

func foo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foo("))
		next.ServeHTTP(w, r)
		w.Write([]byte(")"))
	})
}

func bar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar("))
		next.ServeHTTP(w, r)
		w.Write([]byte(")"))
	})
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

type pipeline struct {
	middlewares []middleware
}

type middleware func(http.Handler) http.Handler

func newPipeline(ms ...middleware) pipeline {
	return pipeline{ms}
}

func (p pipeline) pipe(ms ...middleware) pipeline {
	return pipeline{append(p.middlewares, ms...)}
}

func (p pipeline) process(h http.Handler) http.Handler {
	for i := range p.middlewares {
		h = p.middlewares[len(p.middlewares)-1-i](h)
	}

	return h
}

func main() {
	http.Handle("/", newPipeline().pipe(foo, bar).process(http.HandlerFunc(test)))
	http.ListenAndServe(":8080", nil)
}
