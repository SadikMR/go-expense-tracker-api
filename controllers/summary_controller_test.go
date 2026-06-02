package controllers

import (
	"net/http"
	"testing"
)

func TestGetSummary(t *testing.T) {
	userID := registerAndLogin(t, "Summary User", "summary@example.com", "secret123")

	seeds := []string{
		`{"title":"Lunch","amount":100.00,"category":"Food","expense_date":"2025-06-01"}`,
		`{"title":"Dinner","amount":200.00,"category":"Food","expense_date":"2025-06-15"}`,
		`{"title":"Bus","amount":50.00,"category":"Transport","expense_date":"2025-06-15"}`,
		`{"title":"Rent","amount":300.00,"category":"Housing","expense_date":"2025-07-01"}`,
	}
	for _, e := range seeds {
		makeRequestWithHeaders(t, http.MethodPost, "/api/v1/expenses", e, authHeader(userID))
	}

	tests := []struct {
		name            string
		path            string
		wantStatus      int
		wantSuccess     bool
		wantTotalCount  int
		wantTotalAmount float64
		wantCategories  int
	}{
		{
			name:            "no date filter — all expenses",
			path:            "/api/v1/expenses/summary",
			wantStatus:      http.StatusOK,
			wantSuccess:     true,
			wantTotalCount:  4,
			wantTotalAmount: 650.00,
			wantCategories:  3,
		},
		{
			name:            "date_from only",
			path:            "/api/v1/expenses/summary?date_from=2025-06-15",
			wantStatus:      http.StatusOK,
			wantSuccess:     true,
			wantTotalCount:  3,
			wantTotalAmount: 550.00,
			wantCategories:  3,
		},
		{
			name:            "date_to only",
			path:            "/api/v1/expenses/summary?date_to=2025-06-15",
			wantStatus:      http.StatusOK,
			wantSuccess:     true,
			wantTotalCount:  3,
			wantTotalAmount: 350.00,
			wantCategories:  2,
		},
		{
			name:            "date range — june only",
			path:            "/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30",
			wantStatus:      http.StatusOK,
			wantSuccess:     true,
			wantTotalCount:  3,
			wantTotalAmount: 350.00,
			wantCategories:  2,
		},
		{
			name:            "date range — july only",
			path:            "/api/v1/expenses/summary?date_from=2025-07-01&date_to=2025-07-31",
			wantStatus:      http.StatusOK,
			wantSuccess:     true,
			wantTotalCount:  1,
			wantTotalAmount: 300.00,
			wantCategories:  1,
		},
		{
			name:            "date range no match",
			path:            "/api/v1/expenses/summary?date_from=2024-01-01&date_to=2024-12-31",
			wantStatus:      http.StatusOK,
			wantSuccess:     true,
			wantTotalCount:  0,
			wantTotalAmount: 0,
			wantCategories:  0,
		},
		{
			name:        "invalid date_from format",
			path:        "/api/v1/expenses/summary?date_from=01-06-2025",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "invalid date_to format",
			path:        "/api/v1/expenses/summary?date_to=01-06-2025",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "no auth header",
			path:        "/api/v1/expenses/summary",
			wantStatus:  http.StatusUnauthorized,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := authHeader(userID)
			if tt.name == "no auth header" {
				headers = nil
			}

			w := makeRequestWithHeaders(t, http.MethodGet, tt.path, "", headers)
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

			data, ok := body["data"].(map[string]interface{})
			if !ok {
				t.Fatal("expected data object in response")
			}

			totalCount := int(data["total_count"].(float64))
			if totalCount != tt.wantTotalCount {
				t.Errorf("total_count: want %d, got %d", tt.wantTotalCount, totalCount)
			}

			totalAmount := data["total_amount"].(float64)
			if totalAmount != tt.wantTotalAmount {
				t.Errorf("total_amount: want %.2f, got %.2f", tt.wantTotalAmount, totalAmount)
			}

			byCategory, ok := data["by_category"].([]interface{})
			if !ok {
				t.Fatal("expected by_category array in response")
			}
			if len(byCategory) != tt.wantCategories {
				t.Errorf("by_category count: want %d, got %d", tt.wantCategories, len(byCategory))
			}

			for _, item := range byCategory {
				cat := item.(map[string]interface{})
				if cat["category"] == nil {
					t.Error("by_category item missing category field")
				}
				if cat["total"] == nil {
					t.Error("by_category item missing total field")
				}
				if cat["count"] == nil {
					t.Error("by_category item missing count field")
				}
			}
		})
	}
}
