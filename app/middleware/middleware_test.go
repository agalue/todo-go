package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"gotest.tools/v3/assert"
)

func TestObserver(t *testing.T) {
	count := 0
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count++
		if r.URL.Path == "/fail" {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	obs := NewObserver(context.Background(), router)
	srv := httptest.NewServer(obs)
	defer srv.Close()

	for range 5 {
		_, err := http.Get(srv.URL + "/test")
		assert.NilError(t, err)
		_, err = http.Get(srv.URL + "/fail")
		assert.NilError(t, err)
	}

	assert.Equal(t, 10, count)
	assert.Equal(t, float64(5), testutil.ToFloat64(obs.totalRequests.WithLabelValues("GET", "/test", "200")))
	assert.Equal(t, float64(5), testutil.ToFloat64(obs.totalRequests.WithLabelValues("GET", "/fail", "400")))
}
