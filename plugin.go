// Package headerprocessor - Minimaler Traefik Middleware Plugin
package headerprocessor

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Config Plugin-Konfiguration
type Config struct {
	// Welchen Request-Header wir prüfen
	CheckHeader string `json:"checkHeader,omitempty"`

	// Welchen Wert muss der Header haben, damit etwas passiert (optional, leer = egal ob gesetzt)
	CheckValue string `json:"checkValue,omitempty"`

	// Header der ans Backend weitergegeben wird
	AddRequestHeader string `json:"addRequestHeader,omitempty"`
	AddRequestValue  string `json:"addRequestValue,omitempty"`

	// Response-Header der zurück zum Client geht
	AddResponseHeader string `json:"addResponseHeader,omitempty"`
	AddResponseValue  string `json:"addResponseValue,omitempty"`
}

// CreateConfig Erstellt Default-Konfiguration
func CreateConfig() *Config {
	return &Config{}
}

// HeaderProcessor Der Middleware-Typ
type HeaderProcessor struct {
	next     http.Handler
	name     string

	checkHeader       string
	checkValue        string
	addReqHeader      string
	addReqValue       string
	addRespHeader     string
	addRespValue      string
}

// New Erstellt die Middleware
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.CheckHeader == "" {
		return nil, fmt.Errorf("checkHeader darf nicht leer sein")
	}

	return &HeaderProcessor{
		next:          next,
		name:          name,
		checkHeader:   config.CheckHeader,
		checkValue:    config.CheckValue,
		addReqHeader:  config.AddRequestHeader,
		addReqValue:   config.AddRequestValue,
		addRespHeader: config.AddResponseHeader,
		addRespValue:  config.AddResponseValue,
	}, nil
}

// ServeHTTP Middleware-Logik
func (h *HeaderProcessor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Request Header auswerten
	headerValue := req.Header.Get(h.checkHeader)

	triggered := false
	if h.checkValue == "" {
		// Nur prüfen ob gesetzt
		triggered = headerValue != ""
	} else {
		triggered = headerValue == h.checkValue
	}

	if triggered {
		log.Printf("[Middleware %s] Triggered! %s=%s", h.name, h.checkHeader, headerValue)

		// Header ans Backend hinzufügen
		if h.addReqHeader != "" {
			value := h.addReqValue
			if value == "" {
				value = headerValue // z.B. gleichen Wert weiterreichen
			}
			req.Header.Set(h.addReqHeader, value)
		}
	}

	// Response Header setzen (über Wrapper, damit er nach der Antwort gesetzt wird)
	if h.addRespHeader != "" {
		rw.Header().Set(h.addRespHeader, h.addRespValue)
	}

	h.next.ServeHTTP(rw, req)
}