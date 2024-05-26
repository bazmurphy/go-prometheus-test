package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// define a counter metric
	requestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
	)

	// define a histogram metric
	requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	// register the metrics with Prometheus
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// setup the prometheus request duration timer
		timer := prometheus.NewTimer(requestDuration)
		//and get the request/response duration when the http handler finishes
		defer timer.ObserveDuration()

		log.Println("request to / received...")

		// add a random sleep delay (between 0 and 1000 milliseconds)
		sleepDuration := time.Duration(rand.Intn(2000)) * time.Millisecond
		log.Printf("sleeping for %v...\n", sleepDuration)

		time.Sleep(sleepDuration)

		log.Println("incrementing request counter...")

		// increment the prometheus request counter
		requestCounter.Inc()

		log.Println("writing response...")

		message := fmt.Sprintf("hello from the app\nwaited %s sleepDuration for a response\nprometheus scrapes metrics from the app at /metrics\n", sleepDuration)

		w.Write([]byte(message))
	})

	// set up a http handler for prometheus to scrape metrics from
	http.Handle("/metrics", promhttp.Handler())

	// start the HTTP server
	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
