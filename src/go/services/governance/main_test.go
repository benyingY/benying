package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildServer(t *testing.T) {
	router := buildServer()
	if router == nil {
		t.Fatal("expected router")
	}
}

func TestServiceName(t *testing.T) {
	if serviceName == "" {
		t.Fatal("service name should not be empty")
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := buildServer()
	request, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload struct {
		Status  string `json:"status"`
		Service string `json:"service"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Status != "ok" {
		t.Fatalf("expected status ok, got %q", payload.Status)
	}
	if payload.Service != serviceName {
		t.Fatalf("expected service %q, got %q", serviceName, payload.Service)
	}
}
