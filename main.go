package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
	"webserver/metrics"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	ip   string
	port uint

	errorRate uint
	delay     uint

	numRequest uint
)

var errorMessage = struct {
	Message string
}{
	Message: "error occured",
}

func init() {
	flag.UintVar(&errorRate, "error-rate", 0, "error rate in percentage (0 - 100%")
	flag.UintVar(&delay, "delay", 0, "delay in milli-seconds")

	flag.StringVar(&ip, "ip", "0.0.0.0", "Webserver listen address")
	flag.UintVar(&port, "port", 8080, "Webserver listen address")
}

func handle(w http.ResponseWriter, r *http.Request) {
	var (
		status int
	)

	// Increment num requests counter
	numRequest++

	log.WithFields(log.Fields{
		"method": r.Method,
		"status": status,
		"count":  numRequest,
	}).Info("received request")

	// Start measuring duration
	start := time.Now()
	defer func() {
		elapsed := time.Since(start).Seconds()
		log.WithFields(log.Fields{
			"method":              r.Method,
			"status":              status,
			"count":               numRequest,
			"processing_time(ms)": elapsed * 1000,
		}).Info("time taken for processing request")
		metrics.WebserverRequestDurationSeconds.With(
			prometheus.Labels{
				"method": r.Method,
			},
		).Observe(elapsed)
	}()

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		status = http.StatusBadRequest
		metrics.WebserverErrors.With(
			prometheus.Labels{
				"method": r.Method, "status": strconv.Itoa(status),
			},
		).Inc()

		metrics.WebserverRequests.With(
			prometheus.Labels{
				"method": r.Method, "status": strconv.Itoa(status),
			},
		).Inc()

		log.WithFields(log.Fields{
			"method": r.Method,
			"status": status,
			"count":  numRequest,
		}).WithError(err).Error("bad request")

		w.WriteHeader(status)
		return
	}

	if delay > 0 {
		log.WithFields(log.Fields{
			"method":    r.Method,
			"delay(ms)": delay,
			"count":     numRequest,
		}).Info("delaying request processing")
		// Add delay before sending the response
		<-time.After(1000000 * time.Duration(delay))
	}

	if errorRate > 0 {
		if numRequest%(100/errorRate) == 0 {
			status = http.StatusInternalServerError
			metrics.WebserverErrors.With(
				prometheus.Labels{
					"method": r.Method, "status": strconv.Itoa(status),
				},
			).Inc()

			metrics.WebserverRequests.With(
				prometheus.Labels{
					"method": r.Method, "status": strconv.Itoa(status),
				},
			).Inc()

			log.WithFields(log.Fields{
				"method": r.Method,
				"count":  numRequest,
			}).Info("introducing error")

			w.WriteHeader(status)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	size, err := w.Write(body)
	if err != nil {
		status = http.StatusBadRequest
		metrics.WebserverErrors.With(
			prometheus.Labels{
				"method": r.Method, "status": strconv.Itoa(status),
			},
		).Inc()

		status = http.StatusOK
		metrics.WebserverRequests.With(
			prometheus.Labels{
				"method": r.Method, "status": strconv.Itoa(status),
			},
		).Inc()

		log.WithFields(log.Fields{
			"method": r.Method,
			"status": status,
			"count":  numRequest,
		}).WithError(err).Error("failed writing to response body")

		w.WriteHeader(status)
		return
	}

	log.WithFields(log.Fields{
		"method":     r.Method,
		"delay":      delay,
		"count":      numRequest,
		"error_rate": errorRate,
		"size":       size,
	}).Info("sent response")

	status = http.StatusOK
	metrics.WebserverRequests.With(
		prometheus.Labels{
			"method": r.Method, "status": strconv.Itoa(status),
		},
	).Inc()
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	http.HandleFunc("/hello", handle)
	http.HandleFunc("/healthz", health)
	quit := make(chan struct{})
	defer close(quit)

	if port > 65535 {
		err := errors.New("port number must be in range 0 - 65535")
		log.WithFields(log.Fields{
			"ip":   ip,
			"port": port,
		}).WithError(err).Fatal()
	}
	sWebServer := strconv.Itoa(int(port))
	sPromServer := strconv.Itoa(int(port + 1))

	go func() {
		addr := ip + ":" + sPromServer
		log.Printf("Serving prometheus metrics on %s\n", addr)
		log.Fatal(metrics.StartServer(addr))
	}()

	go func() {
		addr := ip + ":" + sWebServer
		log.Printf("Serving requests on %s\n", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()

	interruptHandler()
}

func interruptHandler() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	os.Exit(1)
}
