// Package simple_traefik_plugin is a Traefik middleware plugin
// that checks for a specific request header and sets a response header accordingly.
package traefik_plugin

import (
	"context"
	"net/http"
)

// Config holds the plugin configuration.
type Config struct {
	// RequestHeader is the name of the request header to check.
	RequestHeader string `json:"requestHeader,omitempty"`
	// ResponseHeader is the name of the response header to set.
	ResponseHeader string `json:"responseHeader,omitempty"`
	// ResponseValue is the value set when the request header is present.
	ResponseValue string `json:"responseValue,omitempty"`
	// MissingValue is the value set when the request header is missing.
	MissingValue string `json:"missingValue,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		RequestHeader:  "X-Custom-Request",
		ResponseHeader: "X-Custom-Response",
		ResponseValue:  "header-was-present",
		MissingValue:   "header-was-missing",
	}
}

// HeaderChecker is the plugin struct.
type HeaderChecker struct {
	next   http.Handler
	config *Config
	name   string
}

// New creates a new plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HeaderChecker{
		next:   next,
		config: config,
		name:   name,
	}, nil
}

// responseWriterWrapper wraps http.ResponseWriter to inject a header before writing.
type responseWriterWrapper struct {
	http.ResponseWriter
	headerValue string
	headerName  string
	wroteHeader bool
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	if !w.wroteHeader {
		w.Header().Set(w.headerName, w.headerValue)
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.Header().Set(w.headerName, w.headerValue)
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(b)
}

// ServeHTTP implements the http.Handler interface.
func (h *HeaderChecker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	value := h.config.MissingValue
	if req.Header.Get(h.config.RequestHeader) != "" {
		value = h.config.ResponseValue
	}

	wrapper := &responseWriterWrapper{
		ResponseWriter: rw,
		headerName:     h.config.ResponseHeader,
		headerValue:    value,
	}

	h.next.ServeHTTP(wrapper, req)
}
