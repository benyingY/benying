package server

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthz(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
    rr := httptest.NewRecorder()

    Handler().ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", rr.Code)
    }
    if rr.Body.String() != `{"status":"ok"}` {
        t.Fatalf("unexpected body: %s", rr.Body.String())
    }
}
