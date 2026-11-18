package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/config"
	"github.com/nuryanto2121/laundryos/apps/api/internal/db"
	"github.com/nuryanto2121/laundryos/apps/api/internal/handler"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// Run bootstraps the application components and manages the server lifecycle.
func Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	ctx := context.Background()

	// Initialize Database
	dbPool, err := db.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}
	defer dbPool.Close()

	// Initialize Customer module
	customerRepo := repository.NewCustomerRepository(dbPool)
	customerSvc := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerSvc)

	// Initialize Order module
	orderRepo := repository.NewOrderRepository(dbPool)
	orderSvc := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// Initialize Bot module
	serviceRepo := repository.NewServiceRepository(dbPool)
	notifSvc := service.NewNotificationService(cfg.GatewayURL, cfg.StaffPhone)
	botSvc := service.NewBotService(service.BotServiceConfig{
		CustomerService: customerSvc,
		OrderService:    orderSvc,
		ServiceRepo:     serviceRepo,
		NotificationSvc: notifSvc,
		GatewayURL:      cfg.GatewayURL,
	})
	botHandler := handler.NewBotHandler(botSvc)

	// Initialize HTTP Router
	router := handler.NewRouter(handler.RouterParams{
		AppVersion:      cfg.AppVersion,
		CustomerHandler: customerHandler,
		OrderHandler:    orderHandler,
		BotHandler:      botHandler,
	})

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server asynchronously
	go func() {
		slog.Info("laundryos-api starting", "version", cfg.AppVersion, "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("server exiting")
	return nil
}
