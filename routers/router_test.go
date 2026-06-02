package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestHealthRouteRegistered(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestSummaryRouteProtected(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses/summary", nil)
	w := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for protected summary route, got %d", w.Code)
	}
}
