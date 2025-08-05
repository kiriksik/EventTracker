package telemetry

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	EventsTotal *prometheus.CounterVec
	ErrorsTotal *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		EventsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "eventtracker",
			Name:      "events_total",
			Help:      "Total number of received events",
		}, []string{"source"}), // e.g. "http", "ws"
		ErrorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "eventtracker",
			Name:      "errors_total",
			Help:      "Total number of errors",
		}, []string{"component"}),
	}
}

func (m *Metrics) Register() {
	prometheus.MustRegister(m.EventsTotal, m.ErrorsTotal)
}

func StartMetricsServer(port string) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("[Telemetry] Prometheus metrics available at : " + port + "/metrics")
		_ = http.ListenAndServe(port, nil)
	}()
}
