package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cermakm/example-apps/echoserver/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	host = "0.0.0.0"
	port = "8080"
)

func initMetricsServer(srv *http.Server, doneCh chan<- bool) {
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	addr := host + ":" + "9178"
	srv.Addr = addr
	srv.Handler = metricsMux

	log.Printf("serving metrics server on: http://%s\n", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("error occurred while serving metrics server: %s", err)
		}
		log.Println("stopped metrics server")
		doneCh <- true
	}()
}

func main() {
	ctx := context.Background()
	doneCh := make(chan bool, 2)

	metricsServer := &http.Server{}
	initMetricsServer(metricsServer, doneCh)

	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)

	if val, exists := os.LookupEnv("SERVER_HOST"); exists {
		host = val
	}
	if val, exists := os.LookupEnv("SERVER_PORT"); exists {
		port = val
	}

	addr := host + ":" + port
	s := &http.Server{
		Addr:    addr,
		Handler: middleware.NewMetricsMiddleware(mux),
	}

	log.Printf("serving echoserver on: http://%s\n", addr)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("error occurred while serving echoserver: %s", err)
		}
		log.Println("stopped main server")

		metricsServer.Shutdown(ctx)
		doneCh <- true
	}()

	<-doneCh
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Engeto: Kubernetes Example Application\n\n"))

	fmt.Fprintf(w, "request received: %#v", r)
}
