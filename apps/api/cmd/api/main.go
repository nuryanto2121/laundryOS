// @title LaundryOS API
// @version 0.1.0
// @description Backend REST API for LaundryOS.
// @host localhost:8080
// @BasePath /
package main

import (
	"log/slog"
	"os"

	"github.com/nuryanto2121/laundryos/apps/api/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("application fatal error", "error", err)
		os.Exit(1)
	}
}
