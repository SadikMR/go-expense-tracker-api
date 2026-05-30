package controllers

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		wantStatus  int
		wantSuccess bool
		wantMessage string
	}{
		{
			name:        "returns 200 when server is running",
			method:      http.MethodGet,
			path:        "/api/v1/health",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantMessage: "Server is running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequest(t, tt.method, tt.path, "")
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			if body["success"] != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, body["success"])
			}
			if body["message"] != tt.wantMessage {
				t.Errorf("message: want %q, got %q", tt.wantMessage, body["message"])
			}
		})
	}
}
