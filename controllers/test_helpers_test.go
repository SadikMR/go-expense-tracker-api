package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/models"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	beego "github.com/beego/beego/v2/server/web"
)

// TestMain runs once before all tests in the controllers package.
func TestMain(m *testing.M) {
	utils.InitConfig()
	utils.InitLogger("dev")

	if err := os.MkdirAll("data", 0755); err != nil {
		panic("TestMain: failed to create data dir: " + err.Error())
	}

	if err := models.EnsureUsersCSV(); err != nil {
		panic("TestMain: failed to initialise users.csv: " + err.Error())
	}

	beego.Router("/api/v1/health", &HealthController{}, "get:Get")
	beego.Router("/api/v1/auth/register", &AuthController{}, "post:Register")
	beego.Router("/api/v1/auth/login", &AuthController{}, "post:Login")

	code := m.Run()

	os.RemoveAll("data")
	os.Exit(code)
}

// makeRequest fires an HTTP request through the Beego router
// and returns the recorded response.
func makeRequest(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, path, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	if err != nil {
		t.Fatalf("makeRequest: failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, req)
	return w
}

// parseResponse unmarshals the response body into a map for assertions.
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("parseResponse: failed to parse %q: %v", w.Body.String(), err)
	}
	return result
}
