package proxy

import (
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	requestsTotal *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {

	m := &metrics{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "requests_total",
				Help: "Number of requests",
			},
			[]string{"outcome"},
		),
	}
	reg.MustRegister(m.requestsTotal)
	return m
}
