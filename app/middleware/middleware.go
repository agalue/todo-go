package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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
	traceProvider    *sdktrace.TracerProvider
}

func newTraceProvider(ctx context.Context) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		slog.Warn("Cannot Initialize OpenTelemetry tracing via gRPC", slog.String("error", err.Error()))
		return tp
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp.RegisterSpanProcessor(bsp)
	return tp
}

func NewObserver(ctx context.Context, mux *http.ServeMux) *Observer {
	mux.Handle("/metrics", promhttp.Handler())
	tp := newTraceProvider(ctx)
	obs := &Observer{
		handler:       otelhttp.NewHandler(mux, "app"),
		traceProvider: tp,
		totalRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_total_requests",
			Help: "The total number of completed requests",
		}, []string{"method", "route", "status"}),
		latencyHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_duration_seconds",
			Help: "Histogram of response time for requests in seconds",
		}, []string{"method", "route", "status"}),
	}
	if err := prometheus.Register(obs.totalRequests); err != nil {
		slog.Warn("cannot register", slog.String("metric", "totalRequests"))
	}
	if err := prometheus.Register(obs.latencyHistogram); err != nil {
		slog.Warn("cannot register", slog.String("metric", "latencyHistogram"))
	}
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

func (o *Observer) Shutdown() {
	if o.traceProvider != nil {
		if err := o.traceProvider.Shutdown(context.Background()); err != nil {
			slog.Error("cannot shutdown tracer", slog.String("error", err.Error()))
		}
	}
}
