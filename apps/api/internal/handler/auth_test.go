package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nuryanto2121/laundryos/apps/api/internal/handler"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

func TestAuthHandler_Validation(t *testing.T) {
	masterRepo := repository.NewMasterRepository(nil)
	svcRepo := repository.NewServiceRepository(nil)
	svc := service.NewMasterService(masterRepo, svcRepo)
	h := handler.NewAuthHandler(svc)

	t.Run("invalid json body returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("{invalid-json"))
		rec := httptest.NewRecorder()

		h.Login(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("missing username or password returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"username":""}`))
		rec := httptest.NewRecorder()

		h.Login(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})
}
