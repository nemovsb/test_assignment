package ginhandlers

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestCount prometheus.Counter
}

func NewMetrics() *Metrics {

	metric := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "requests",
		})
	prometheus.MustRegister(metric)
	return &Metrics{
		RequestCount: metric,
	}
}
