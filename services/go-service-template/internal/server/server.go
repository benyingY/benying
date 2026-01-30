package server

import (
    "net/http"
)

func Handler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", healthz)
    return mux
}

func healthz(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ok"}`))
}
