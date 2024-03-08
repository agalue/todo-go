package app

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type Observer struct {
	handler          http.Handler
	totalRequests    *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
}

func NewObserver(mux *http.ServeMux) *Observer {
	mux.Handle("/metrics", promhttp.Handler())
	obs := &Observer{
		handler: mux,
		totalRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_total_requests",
			Help: "The total number of completed requests",
		}, []string{"method", "route", "status"}),
		latencyHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_duration_seconds",
			Help: "Histogram of response time for requests in seconds",
		}, []string{"method", "route", "status"}),
	}
	prometheus.MustRegister(obs.totalRequests, obs.latencyHistogram)
	return obs
}

func (o *Observer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	recorder := &StatusRecorder{
		ResponseWriter: w,
		Status:         200,
	}
	o.handler.ServeHTTP(recorder, r)
	duration := time.Since(start)
	slog.Info("query executed",
		slog.String("source", r.RemoteAddr),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", recorder.Status),
		slog.Duration("duration", duration),
	)
	o.totalRequests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(recorder.Status)).Inc()
	o.latencyHistogram.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(recorder.Status)).Observe(duration.Seconds())
}
