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

type Logger struct {
	handler          http.Handler
	totalRequests    *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
}

func NewLogger(mux *http.ServeMux) *Logger {
	mux.Handle("/metrics", promhttp.Handler())
	logger := &Logger{
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
	prometheus.MustRegister(logger.totalRequests, logger.latencyHistogram)
	return logger
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	recorder := &StatusRecorder{
		ResponseWriter: w,
		Status:         200,
	}
	l.handler.ServeHTTP(recorder, r)
	duration := time.Since(start)
	slog.Info("query executed",
		slog.String("source", r.RemoteAddr),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", recorder.Status),
		slog.Duration("duration", duration),
	)
	l.totalRequests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(recorder.Status)).Inc()
	l.latencyHistogram.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(recorder.Status)).Observe(duration.Seconds())
}
