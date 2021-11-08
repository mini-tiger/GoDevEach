package main

import (
	"flag"
	"log"
	"net/http"
	"prometheusDemo/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

// https://www.cnblogs.com/makelu/p/11082485.html
// https://prometheus.io/docs/guides/go-application/
func main() {
	flag.Parse()
	services.RecordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
