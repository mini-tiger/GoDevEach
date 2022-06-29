package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

var (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx)
	}
}

func main() {
	r := router.New()
	g1 := r.Group("/v1")
	g1.GET("/", Index)
	g1.GET("/hello/{name}", Hello)

	g2 := r.Group("/v2")
	g2.GET("/hello2/{name}", Hello)
	fmt.Println(r.List())
	//log.Fatal(fasthttp.ListenAndServe(":8081", r.Handler))
	if err := fasthttp.ListenAndServe(":8081", CORS(r.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
