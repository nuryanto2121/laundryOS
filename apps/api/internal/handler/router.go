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
	POSHandler      *POSHandler
	ExpenseHandler  *ExpenseHandler
	ReportHandler   *ReportHandler
	MasterHandler   *MasterHandler
}

// NewRouter constructs and configures the Chi HTTP router with middlewares and route mappings.
func NewRouter(p RouterParams) http.Handler {
	r := chi.NewRouter()
	r.Use(appMiddleware.RequestLogger())
	r.Use(appMiddleware.CORS())

	// System routes
	r.Get("/health", NewHealthHandler(p.AppVersion))
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Bot routes
	if p.BotHandler != nil {
		r.Post("/bot/messages", p.BotHandler.HandleMessage)
	}

	// Customer routes
	if p.CustomerHandler != nil {
		r.Post("/customers", p.CustomerHandler.Create)
		r.Get("/customers", p.CustomerHandler.List)
		r.Get("/customers/{id}", p.CustomerHandler.Get)
		r.Put("/customers/{id}", p.CustomerHandler.Update)
		r.Patch("/customers/{id}/status", p.CustomerHandler.UpdateStatus)
		r.Delete("/customers/{id}", p.CustomerHandler.Delete)
	}

	// Legacy Order routes
	if p.OrderHandler != nil {
		r.Post("/orders", p.OrderHandler.Create)
		r.Get("/orders", p.OrderHandler.List)
		r.Get("/orders/{id}", p.OrderHandler.Get)
		r.Put("/orders/{id}", p.OrderHandler.Update)
		r.Delete("/orders/{id}", p.OrderHandler.Delete)
	}

	// Web POS routes
	if p.POSHandler != nil {
		r.Route("/pos", func(pos chi.Router) {
			pos.Post("/orders", p.POSHandler.Create)
			pos.Get("/orders", p.POSHandler.List)
			pos.Get("/orders/{id}", p.POSHandler.Get)
			pos.Patch("/orders/{id}/payment", p.POSHandler.UpdatePayment)
			pos.Patch("/orders/{id}/status", p.POSHandler.UpdateStatus)
			pos.Patch("/orders/{id}/cancel", p.POSHandler.Cancel)
			pos.Patch("/orders/{id}/rack", p.POSHandler.UpdateRack)
			pos.Post("/orders/{id}/photos", p.POSHandler.AddPhoto)
			pos.Post("/orders/{id}/resend-invoice", p.POSHandler.ResendInvoice)
		})
	}

	// Operational Expense (Petty Cash) routes
	if p.ExpenseHandler != nil {
		r.Route("/expenses", func(exp chi.Router) {
			exp.Post("/", p.ExpenseHandler.Create)
			exp.Get("/", p.ExpenseHandler.List)
			exp.Get("/summary", p.ExpenseHandler.GetSummary)
			exp.Delete("/{id}", p.ExpenseHandler.Delete)
		})
	}

	// Reports & Analytics routes
	if p.ReportHandler != nil {
		r.Route("/reports", func(rep chi.Router) {
			rep.Get("/home-metrics", p.ReportHandler.GetHomeMetrics)
			rep.Get("/financial", p.ReportHandler.GetFinancial)
			rep.Get("/unpaid", p.ReportHandler.ListUnpaid)
			rep.Get("/services", p.ReportHandler.GetServices)
			rep.Get("/top-customers", p.ReportHandler.GetTopCustomers)
			rep.Get("/customers/{id}/crm", p.ReportHandler.GetCustomerCRM)
		})
	}

	// Master Data routes
	if p.MasterHandler != nil {
		r.Route("/master", func(m chi.Router) {
			m.Get("/perfumes", p.MasterHandler.ListPerfumes)
			m.Get("/discounts", p.MasterHandler.ListDiscounts)
			m.Get("/payment-methods", p.MasterHandler.ListPaymentMethods)
			m.Get("/outlets", p.MasterHandler.ListOutlets)
			m.Get("/services", p.MasterHandler.ListServices)
		})
	}

	return r
}
