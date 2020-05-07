// Code generated by go-prometrics-gen. DO NOT EDIT.
package metrics

import (
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (

	// Name:Webserver
	// Title: Webserver RED Metrics

	// WebserverRequests : Counter for HTTP requests handled by the webserver
	WebserverRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webserver_requests",
			Help: "Counter for the number of requests handled",
		},
		[]string{"method", "status"},
	)

	// WebserverErrors : Counter for errored non-200 OK response codes returned by the webserver
	WebserverErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webserver_errors",
			Help: "Counter for the number of errors (non 200 OK)",
		},
		[]string{"method", "status"},
	)
	// WebserverRequestDurationSeconds : Histogram for the handlers organized based on the HTTP request method, and HTTP request path
	WebserverRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "webserver_request_duration_seconds",
			Help:    "Histogram for the runtime of handler methods",
			Buckets: []float64{5e-05, 0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1},
		},
		[]string{"method"},
	)
)

// StartServer starts a blocking web-server with ListenAndServe
func StartServer(address string) error {
	if address == "" {
		return errors.New("prometheus server address must be specified")
	}

	// Start the prometheus metrics server on its unique port
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(address, promhttp.Handler()); err != nil {
		return err
	}
	return nil
}
