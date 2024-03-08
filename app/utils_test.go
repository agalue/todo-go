package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")
	assert.Equal(t, 1, getID(w, r))
}

func TestGetIDInvalid(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos/x", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "x")
	assert.Equal(t, -1, getID(w, r))
}

func TestGetIDNotFound(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	w := httptest.NewRecorder()
	assert.Equal(t, -1, getID(w, r))
}
