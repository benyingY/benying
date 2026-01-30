package main

import (
    "log"
    "net/http"

    "MODULE_PLACEHOLDER/internal/server"
)

func main() {
    addr := ":8080"
    log.Printf("listening on %s", addr)
    if err := http.ListenAndServe(addr, server.Handler()); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
