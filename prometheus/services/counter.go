package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)
import "github.com/prometheus/client_golang/prometheus/promauto"

func RecordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			myGague.Add(11)
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	myGague = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "my_example_gauge_data",
		Help:        "my example gauge data",
		ConstLabels: map[string]string{"error": ""},
	})
)
