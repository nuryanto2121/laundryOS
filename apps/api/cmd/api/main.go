package main

import (
	"log"
	"net/http"

	"github.com/nuryanto2121/laundryos/apps/api/internal/config"
	"github.com/nuryanto2121/laundryos/apps/api/internal/handler"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.NewHealthHandler(cfg.AppVersion))

	addr := ":" + cfg.Port
	log.Printf("laundryos-api v%s starting on %s", cfg.AppVersion, addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
