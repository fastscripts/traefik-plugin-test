package traefik_plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderPresent(t *testing.T) {
	cfg := CreateConfig()
	handler, err := New(context.Background(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom-Request", "some-value")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("X-Custom-Response")
	if got != "header-was-present" {
		t.Errorf("expected 'header-was-present', got '%s'", got)
	}
}

func TestHeaderMissing(t *testing.T) {
	cfg := CreateConfig()
	handler, err := New(context.Background(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("X-Custom-Response")
	if got != "header-was-missing" {
		t.Errorf("expected 'header-was-missing', got '%s'", got)
	}
}
