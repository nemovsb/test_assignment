package metrics

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusServer struct {
	server http.Server
}

func NewPometheusServer(port string) *PrometheusServer {
	handler := promhttp.Handler()
	http.Handle("/metrics", handler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return &PrometheusServer{server: server}
}

func (s *PrometheusServer) Run() error {
	log.Println("Prometheus server run")
	return s.server.ListenAndServe()
}

func (s *PrometheusServer) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}
