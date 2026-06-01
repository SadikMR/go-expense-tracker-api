package controllers

import (
	"fmt"
	"net/http"
	"testing"
)

// authHeader returns the X-User-ID header map for a given user ID string.
func authHeader(userID string) map[string]string {
	return map[string]string{"X-User-ID": userID}
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCreateExpense(t *testing.T) {
	userID := registerAndLogin(t, "Expense User", "expense@example.com", "secret123")

	tests := []struct {
		name        string
		body        string
		headers     map[string]string
		wantStatus  int
		wantSuccess bool
		wantMessage string
	}{
		{
			name:        "valid expense",
			body:        `{"title":"Lunch","amount":350.50,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusCreated,
			wantSuccess: true,
			wantMessage: "Expense created successfully",
		},
		{
			name:        "missing title",
			body:        `{"amount":350.50,"category":"Food","expense_date":"2025-06-10"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "missing category",
			body:        `{"title":"Lunch","amount":350.50,"expense_date":"2025-06-10"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "invalid category",
			body:        `{"title":"Lunch","amount":350.50,"category":"InvalidCat","expense_date":"2025-06-10"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "zero amount",
			body:        `{"title":"Lunch","amount":0,"category":"Food","expense_date":"2025-06-10"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "negative amount",
			body:        `{"title":"Lunch","amount":-10,"category":"Food","expense_date":"2025-06-10"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "invalid date format",
			body:        `{"title":"Lunch","amount":350.50,"category":"Food","expense_date":"10-06-2025"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "missing date",
			body:        `{"title":"Lunch","amount":350.50,"category":"Food"}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "no auth header",
			body:        `{"title":"Lunch","amount":350.50,"category":"Food","expense_date":"2025-06-10"}`,
			headers:     nil,
			wantStatus:  http.StatusUnauthorized,
			wantSuccess: false,
			wantMessage: "Unauthorized",
		},
		{
			name:        "malformed JSON",
			body:        `{invalid}`,
			headers:     authHeader(userID),
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequestWithHeaders(t, http.MethodPost, "/api/v1/expenses", tt.body, tt.headers)
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			if body["success"] != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, body["success"])
			}
			if tt.wantMessage != "" && body["message"] != tt.wantMessage {
				t.Errorf("message: want %q, got %q", tt.wantMessage, body["message"])
			}
		})
	}
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestListExpenses(t *testing.T) {
	userID := registerAndLogin(t, "List User", "list@example.com", "secret123")

	// Seed expenses with different dates and amounts
	seeds := []string{
		`{"title":"Expense A","amount":300,"category":"Food","expense_date":"2025-06-01"}`,
		`{"title":"Expense B","amount":100,"category":"Transport","expense_date":"2025-06-15"}`,
		`{"title":"Expense C","amount":200,"category":"Housing","expense_date":"2025-07-01"}`,
	}
	for _, e := range seeds {
		makeRequestWithHeaders(t, http.MethodPost, "/api/v1/expenses", e, authHeader(userID))
	}

	tests := []struct {
		name           string
		path           string
		wantStatus     int
		wantSuccess    bool
		wantCount      int
		wantFirstTitle string
	}{
		{
			name:        "no filters — returns all",
			path:        "/api/v1/expenses",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantCount:   3,
		},
		{
			name:        "date_from only",
			path:        "/api/v1/expenses?date_from=2025-06-15",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantCount:   2,
		},
		{
			name:        "date_to only",
			path:        "/api/v1/expenses?date_to=2025-06-15",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantCount:   2,
		},
		{
			name:        "date range",
			path:        "/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantCount:   2,
		},
		{
			name:        "date range no match",
			path:        "/api/v1/expenses?date_from=2024-01-01&date_to=2024-12-31",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantCount:   0,
		},
		{
			name:           "sort by amount asc — cheapest first",
			path:           "/api/v1/expenses?sort_by=amount&sort_order=asc",
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantCount:      3,
			wantFirstTitle: "Expense B",
		},
		{
			name:           "sort by amount desc — most expensive first",
			path:           "/api/v1/expenses?sort_by=amount&sort_order=desc",
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantCount:      3,
			wantFirstTitle: "Expense A",
		},
		{
			name:           "sort by expense_date asc — oldest first",
			path:           "/api/v1/expenses?sort_by=expense_date&sort_order=asc",
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantCount:      3,
			wantFirstTitle: "Expense A",
		},
		{
			name:           "sort by expense_date desc — newest first",
			path:           "/api/v1/expenses?sort_by=expense_date&sort_order=desc",
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantCount:      3,
			wantFirstTitle: "Expense C",
		},
		{
			name:           "default sort is expense_date desc",
			path:           "/api/v1/expenses",
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantCount:      3,
			wantFirstTitle: "Expense C",
		},
		{
			name:        "limit",
			path:        "/api/v1/expenses?limit=2",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantCount:   2,
		},
		{
			name:        "limit exceeding total",
			path:        "/api/v1/expenses?limit=100",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantCount:   3,
		},
		{
			name:           "date range + sort + limit combined",
			path:           "/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30&sort_by=amount&sort_order=asc&limit=1",
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantCount:      1,
			wantFirstTitle: "Expense B",
		},
		{
			name:        "invalid sort_by",
			path:        "/api/v1/expenses?sort_by=title",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "invalid sort_order",
			path:        "/api/v1/expenses?sort_order=random",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "invalid date_from format",
			path:        "/api/v1/expenses?date_from=01-06-2025",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "invalid date_to format",
			path:        "/api/v1/expenses?date_to=01-06-2025",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequestWithHeaders(t, http.MethodGet, tt.path, "", authHeader(userID))
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			if body["success"] != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, body["success"])
			}

			if !tt.wantSuccess {
				return
			}

			data, ok := body["data"].([]interface{})
			if !ok {
				t.Fatal("expected data array in response")
			}
			if len(data) != tt.wantCount {
				t.Errorf("count: want %d, got %d", tt.wantCount, len(data))
			}

			if tt.wantFirstTitle != "" && len(data) > 0 {
				first := data[0].(map[string]interface{})
				if first["title"] != tt.wantFirstTitle {
					t.Errorf("first title: want %q, got %q", tt.wantFirstTitle, first["title"])
				}
			}
		})
	}
}

// ── Get One ───────────────────────────────────────────────────────────────────

func TestGetExpense(t *testing.T) {
	userID := registerAndLogin(t, "GetOne User", "getone@example.com", "secret123")

	// Seed one expense and grab its ID
	w := makeRequestWithHeaders(t, http.MethodPost, "/api/v1/expenses",
		`{"title":"Solo Expense","amount":200,"category":"Transport","expense_date":"2025-06-01"}`,
		authHeader(userID),
	)
	created := parseResponse(t, w)
	data := created["data"].(map[string]interface{})
	expenseID := fmt.Sprintf("%d", int(data["id"].(float64)))

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantSuccess bool
	}{
		{
			name:        "get existing expense",
			path:        "/api/v1/expenses/" + expenseID,
			wantStatus:  http.StatusOK,
			wantSuccess: true,
		},
		{
			name:        "get non-existent expense",
			path:        "/api/v1/expenses/99999",
			wantStatus:  http.StatusNotFound,
			wantSuccess: false,
		},
		{
			name:        "invalid id",
			path:        "/api/v1/expenses/abc",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequestWithHeaders(t, http.MethodGet, tt.path, "", authHeader(userID))
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			if body["success"] != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, body["success"])
			}
		})
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdateExpense(t *testing.T) {
	userID := registerAndLogin(t, "Update User", "update@example.com", "secret123")

	// Seed one expense
	w := makeRequestWithHeaders(t, http.MethodPost, "/api/v1/expenses",
		`{"title":"Old Title","amount":100,"category":"Food","expense_date":"2025-06-01"}`,
		authHeader(userID),
	)
	created := parseResponse(t, w)
	data := created["data"].(map[string]interface{})
	expenseID := fmt.Sprintf("%d", int(data["id"].(float64)))

	tests := []struct {
		name        string
		path        string
		body        string
		wantStatus  int
		wantSuccess bool
	}{
		{
			name:        "valid update",
			path:        "/api/v1/expenses/" + expenseID,
			body:        `{"title":"New Title","amount":500,"category":"Transport","expense_date":"2025-06-15"}`,
			wantStatus:  http.StatusOK,
			wantSuccess: true,
		},
		{
			name:        "invalid category",
			path:        "/api/v1/expenses/" + expenseID,
			body:        `{"title":"New Title","amount":500,"category":"InvalidCat","expense_date":"2025-06-15"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "non-existent expense",
			path:        "/api/v1/expenses/99999",
			body:        `{"title":"New Title","amount":500,"category":"Food","expense_date":"2025-06-15"}`,
			wantStatus:  http.StatusNotFound,
			wantSuccess: false,
		},
		{
			name:        "invalid id",
			path:        "/api/v1/expenses/abc",
			body:        `{"title":"New Title","amount":500,"category":"Food","expense_date":"2025-06-15"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequestWithHeaders(t, http.MethodPut, tt.path, tt.body, authHeader(userID))
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			if body["success"] != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, body["success"])
			}
		})
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDeleteExpense(t *testing.T) {
	userID := registerAndLogin(t, "Delete User", "delete@example.com", "secret123")

	// Seed one expense
	w := makeRequestWithHeaders(t, http.MethodPost, "/api/v1/expenses",
		`{"title":"To Delete","amount":100,"category":"Food","expense_date":"2025-06-01"}`,
		authHeader(userID),
	)
	created := parseResponse(t, w)
	data := created["data"].(map[string]interface{})
	expenseID := fmt.Sprintf("%d", int(data["id"].(float64)))

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantSuccess bool
	}{
		{
			name:        "delete existing expense",
			path:        "/api/v1/expenses/" + expenseID,
			wantStatus:  http.StatusOK,
			wantSuccess: true,
		},
		{
			name:        "delete already deleted expense",
			path:        "/api/v1/expenses/" + expenseID,
			wantStatus:  http.StatusNotFound,
			wantSuccess: false,
		},
		{
			name:        "delete non-existent expense",
			path:        "/api/v1/expenses/99999",
			wantStatus:  http.StatusNotFound,
			wantSuccess: false,
		},
		{
			name:        "invalid id",
			path:        "/api/v1/expenses/abc",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequestWithHeaders(t, http.MethodDelete, tt.path, "", authHeader(userID))
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			if body["success"] != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, body["success"])
			}
		})
	}
}

// ── Ownership ─────────────────────────────────────────────────────────────────

func TestExpenseOwnership(t *testing.T) {
	// Two separate users
	userA := registerAndLogin(t, "User A", "usera@example.com", "secret123")
	userB := registerAndLogin(t, "User B", "userb@example.com", "secret123")

	// User A creates an expense
	w := makeRequestWithHeaders(t, http.MethodPost, "/api/v1/expenses",
		`{"title":"User A Expense","amount":100,"category":"Food","expense_date":"2025-06-01"}`,
		authHeader(userA),
	)
	created := parseResponse(t, w)
	data := created["data"].(map[string]interface{})
	expenseID := fmt.Sprintf("%d", int(data["id"].(float64)))

	tests := []struct {
		name        string
		method      string
		path        string
		body        string
		userID      string
		wantStatus  int
		wantSuccess bool
	}{
		{
			name:        "user B cannot get user A expense",
			method:      http.MethodGet,
			path:        "/api/v1/expenses/" + expenseID,
			userID:      userB,
			wantStatus:  http.StatusNotFound,
			wantSuccess: false,
		},
		{
			name:        "user B cannot update user A expense",
			method:      http.MethodPut,
			path:        "/api/v1/expenses/" + expenseID,
			body:        `{"title":"Stolen","amount":1,"category":"Food","expense_date":"2025-06-01"}`,
			userID:      userB,
			wantStatus:  http.StatusNotFound,
			wantSuccess: false,
		},
		{
			name:        "user B cannot delete user A expense",
			method:      http.MethodDelete,
			path:        "/api/v1/expenses/" + expenseID,
			userID:      userB,
			wantStatus:  http.StatusNotFound,
			wantSuccess: false,
		},
		{
			name:        "user A can get own expense",
			method:      http.MethodGet,
			path:        "/api/v1/expenses/" + expenseID,
			userID:      userA,
			wantStatus:  http.StatusOK,
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequestWithHeaders(t, tt.method, tt.path, tt.body, authHeader(tt.userID))
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			if body["success"] != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, body["success"])
			}
		})
	}
}
