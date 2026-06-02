package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/models"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	beego "github.com/beego/beego/v2/server/web"
)

func TestMain(m *testing.M) {
	utils.InitConfig()
	utils.InitLogger("dev")

	if err := os.MkdirAll("data", 0755); err != nil {
		panic("TestMain: failed to create data dir: " + err.Error())
	}
	if err := models.EnsureUsersCSV(); err != nil {
		panic("TestMain: failed to initialise users.csv: " + err.Error())
	}
	if err := models.EnsureExpensesCSV(); err != nil {
		panic("TestMain: failed to initialise expenses.csv: " + err.Error())
	}

	beego.Router("/api/v1/health", &HealthController{}, "get:Get")
	beego.Router("/api/v1/auth/register", &AuthController{}, "post:Register")
	beego.Router("/api/v1/auth/login", &AuthController{}, "post:Login")
	beego.Router("/api/v1/expenses/summary", &SummaryController{}, "get:Get")
	beego.Router("/api/v1/expenses", &ExpenseController{}, "post:Post;get:Get")
	beego.Router("/api/v1/expenses/:id", &ExpenseController{}, "get:Get;put:Put;delete:Delete")

	code := m.Run()

	os.RemoveAll("data")
	os.Exit(code)
}

// makeRequest fires a request with no extra headers.
func makeRequest(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	return makeRequestWithHeaders(t, method, path, body, nil)
}

// makeRequestWithHeaders fires a request with custom headers.
func makeRequestWithHeaders(t *testing.T, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, path, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	if err != nil {
		t.Fatalf("makeRequestWithHeaders: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

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

// registerAndLogin registers a user and returns their user_id as a string
// ready for use in the X-User-ID header.
func registerAndLogin(t *testing.T, name, email, password string) string {
	t.Helper()

	makeRequest(t, http.MethodPost, "/api/v1/auth/register",
		fmt.Sprintf(`{"name":%q,"email":%q,"password":%q}`, name, email, password),
	)

	w := makeRequest(t, http.MethodPost, "/api/v1/auth/login",
		fmt.Sprintf(`{"email":%q,"password":%q}`, email, password),
	)

	body := parseResponse(t, w)
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatal("registerAndLogin: no data in login response")
	}
	userID, ok := data["user_id"].(float64)
	if !ok {
		t.Fatal("registerAndLogin: no user_id in data")
	}
	return fmt.Sprintf("%d", int(userID))
}
