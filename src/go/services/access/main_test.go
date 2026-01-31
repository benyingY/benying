package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func setupAuthToken(t *testing.T) string {
	t.Helper()
	privateKey, publicKeyPEM := generateTestKeypair(t)
	os.Setenv("ACCESS_JWT_PUBLIC_KEY", publicKeyPEM)
	t.Cleanup(func() {
		os.Unsetenv("ACCESS_JWT_PUBLIC_KEY")
	})

	claims := authClaims{
		Subject:   "user-1",
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	return signToken(t, privateKey, claims)
}

func setupInvalidToken(t *testing.T) string {
	t.Helper()
	_, publicKeyPEM := generateTestKeypair(t)
	os.Setenv("ACCESS_JWT_PUBLIC_KEY", publicKeyPEM)
	t.Cleanup(func() {
		os.Unsetenv("ACCESS_JWT_PUBLIC_KEY")
	})

	privateKey, _ := generateTestKeypair(t)
	claims := authClaims{
		Subject:   "user-1",
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	return signToken(t, privateKey, claims)
}

func generateTestKeypair(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})
	return privateKey, string(publicPEM)
}

func signToken(t *testing.T, privateKey *rsa.PrivateKey, claims authClaims) string {
	t.Helper()
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("failed to marshal header: %v", err)
	}
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("failed to marshal claims: %v", err)
	}
	headerSegment := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadSegment := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := headerSegment + "." + payloadSegment
	hashed := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	signatureSegment := base64.RawURLEncoding.EncodeToString(signature)
	return signingInput + "." + signatureSegment
}

func TestInvokePropagatesRequestID(t *testing.T) {
	router := buildServer()
	token := setupAuthToken(t)
	body, err := json.Marshal(invokeRequest{
		Input: "hello",
		Metadata: map[string]string{
			"source": "test",
		},
	})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(body))
	req.Header.Set("X-Request-Id", "req-12345")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	responseRequestID := resp.Header().Get("X-Request-Id")
	if responseRequestID != "req-12345" {
		t.Fatalf("expected request id to be propagated, got %q", responseRequestID)
	}

	var payload invokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Status != "ok" {
		t.Fatalf("expected status ok, got %q", payload.Status)
	}
	if payload.RequestID != "req-12345" {
		t.Fatalf("expected request id in payload, got %q", payload.RequestID)
	}
	if payload.Service != serviceName {
		t.Fatalf("expected service name %q, got %q", serviceName, payload.Service)
	}
	if payload.Input != "hello" {
		t.Fatalf("expected input to be echoed, got %q", payload.Input)
	}
	if payload.Output != "hello" {
		t.Fatalf("expected output to match input, got %q", payload.Output)
	}
	if payload.Metadata["source"] != "test" {
		t.Fatalf("expected metadata to be echoed, got %v", payload.Metadata)
	}
}

func TestInvokeGeneratesRequestIDWhenMissing(t *testing.T) {
	router := buildServer()
	token := setupAuthToken(t)
	body, err := json.Marshal(invokeRequest{
		Input: "hello",
	})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	responseRequestID := resp.Header().Get("X-Request-Id")
	if responseRequestID == "" {
		t.Fatalf("expected generated request id, got empty")
	}

	var payload invokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.RequestID != responseRequestID {
		t.Fatalf("expected payload request id %q to match header", responseRequestID)
	}
}

func TestInvokeRejectsMissingInput(t *testing.T) {
	router := buildServer()
	token := setupAuthToken(t)

	req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.Code)
	}

	var payload errorResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if payload.Status != "error" {
		t.Fatalf("expected error status, got %q", payload.Status)
	}
	if payload.Error == "" {
		t.Fatalf("expected error code to be set")
	}
	if payload.Service != serviceName {
		t.Fatalf("expected service name %q, got %q", serviceName, payload.Service)
	}
	if payload.RequestID == "" {
		t.Fatalf("expected request id to be set on error response")
	}
}

func TestInvokeStreamsSSE(t *testing.T) {
	router := buildServer()
	token := setupAuthToken(t)
	body, err := json.Marshal(invokeRequest{
		Input: "hello",
	})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/invoke?stream=true", bytes.NewBuffer(body))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	contentType := resp.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/event-stream") {
		t.Fatalf("expected SSE content type, got %q", contentType)
	}

	bodyText := resp.Body.String()
	if !strings.Contains(bodyText, "data: ") {
		t.Fatalf("expected data frames, got %q", bodyText)
	}
	if !strings.Contains(bodyText, "\"object\":\"chat.completion.chunk\"") {
		t.Fatalf("expected openai chunk object, got %q", bodyText)
	}
	if !strings.Contains(bodyText, "\"content\":\"hello\"") {
		t.Fatalf("expected streamed content, got %q", bodyText)
	}
	if !strings.Contains(bodyText, "\"finish_reason\":\"stop\"") {
		t.Fatalf("expected finish_reason stop, got %q", bodyText)
	}
	if !strings.Contains(bodyText, "data: [DONE]") {
		t.Fatalf("expected done sentinel, got %q", bodyText)
	}
}

func TestInvokeRejectsMissingAuth(t *testing.T) {
	router := buildServer()
	body, err := json.Marshal(invokeRequest{
		Input: "hello",
	})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}

func TestInvokeRejectsInvalidToken(t *testing.T) {
	router := buildServer()
	token := setupInvalidToken(t)
	body, err := json.Marshal(invokeRequest{
		Input: "hello",
	})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}

func TestHealthDoesNotRequireAuth(t *testing.T) {
	router := buildServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}
