package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/nuryanto2121/laundryos/apps/api/docs"
	appMiddleware "github.com/nuryanto2121/laundryos/apps/api/internal/handler/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// RouterParams carries all dependencies required to build the HTTP router.
type RouterParams struct {
	AppVersion      string
	CustomerHandler *CustomerHandler
	OrderHandler    *OrderHandler
	BotHandler      *BotHandler
}

// NewRouter constructs and configures the Chi HTTP router with middlewares and route mappings.
func NewRouter(p RouterParams) http.Handler {
	r := chi.NewRouter()
	r.Use(appMiddleware.RequestLogger())

	// System routes
	r.Get("/health", NewHealthHandler(p.AppVersion))
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Bot routes
	if p.BotHandler != nil {
		r.Post("/bot/messages", p.BotHandler.HandleMessage)
	}

	// Customer routes
	r.Post("/customers", p.CustomerHandler.Create)
	r.Get("/customers", p.CustomerHandler.List)
	r.Get("/customers/{id}", p.CustomerHandler.Get)
	r.Put("/customers/{id}", p.CustomerHandler.Update)
	r.Patch("/customers/{id}/status", p.CustomerHandler.UpdateStatus)
	r.Delete("/customers/{id}", p.CustomerHandler.Delete)

	// Order routes
	r.Post("/orders", p.OrderHandler.Create)
	r.Get("/orders", p.OrderHandler.List)
	r.Get("/orders/{id}", p.OrderHandler.Get)
	r.Put("/orders/{id}", p.OrderHandler.Update)
	r.Delete("/orders/{id}", p.OrderHandler.Delete)

	return r
}
