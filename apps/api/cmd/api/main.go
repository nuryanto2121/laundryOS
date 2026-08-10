package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/config"
	"github.com/nuryanto2121/laundryos/apps/api/internal/db"
	"github.com/nuryanto2121/laundryos/apps/api/internal/handler"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

func main() {
	// Initialize slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	ctx := context.Background()

	// Initialize Database
	dbPool, err := db.New(ctx, cfg)
	if err != nil {
		slog.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Initialize Customer module
	customerRepo := repository.NewCustomerRepository(dbPool)
	customerSvc := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerSvc)

	// Initialize HTTP Router
	router := chi.NewRouter()

	router.Get("/health", handler.NewHealthHandler(cfg.AppVersion))

	// Customer routes
	router.Post("/customers", customerHandler.Create)
	router.Get("/customers", customerHandler.List)
	router.Get("/customers/{id}", customerHandler.Get)
	router.Put("/customers/{id}", customerHandler.Update)
	router.Patch("/customers/{id}/status", customerHandler.UpdateStatus)
	router.Delete("/customers/{id}", customerHandler.Delete)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server asynchronously
	go func() {
		slog.Info("laundryos-api starting", "version", cfg.AppVersion, "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exiting")
}
