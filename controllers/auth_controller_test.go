package controllers

import (
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantSuccess bool
		wantMessage string
	}{
		{
			name:        "valid registration",
			body:        `{"name":"Sadik MR","email":"sadik@example.com","password":"secret123"}`,
			wantStatus:  http.StatusCreated,
			wantSuccess: true,
			wantMessage: "User registered successfully",
		},
		{
			name:        "duplicate email",
			body:        `{"name":"Sadik MR","email":"sadik@example.com","password":"secret123"}`,
			wantStatus:  http.StatusConflict,
			wantSuccess: false,
			wantMessage: "Email already exists",
		},
		{
			name:        "missing name",
			body:        `{"email":"noname@example.com","password":"secret123"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Name, email, and password are required",
		},
		{
			name:        "missing email",
			body:        `{"name":"Sadik MR","password":"secret123"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Name, email, and password are required",
		},
		{
			name:        "missing password",
			body:        `{"name":"Sadik MR","email":"nopwd@example.com"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Name, email, and password are required",
		},
		{
			name:        "invalid email format",
			body:        `{"name":"Sadik MR","email":"not-an-email","password":"secret123"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Invalid email format",
		},
		{
			name:        "password too short",
			body:        `{"name":"Sadik MR","email":"short@example.com","password":"abc"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Password must be at least 6 characters",
		},
		{
			name:        "empty body",
			body:        `{}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Name, email, and password are required",
		},
		{
			name:        "malformed JSON",
			body:        `{invalid}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequest(t, http.MethodPost, "/api/v1/auth/register", tt.body)
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

func TestLogin(t *testing.T) {
	// Seed a known user before running login tests
	makeRequest(t, http.MethodPost, "/api/v1/auth/register",
		`{"name":"Login User","email":"login@example.com","password":"secret123"}`,
	)

	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantSuccess bool
		wantMessage string
	}{
		{
			name:        "valid credentials",
			body:        `{"email":"login@example.com","password":"secret123"}`,
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantMessage: "Login successful",
		},
		{
			name:        "wrong password",
			body:        `{"email":"login@example.com","password":"wrongpass"}`,
			wantStatus:  http.StatusUnauthorized,
			wantSuccess: false,
			wantMessage: "Invalid email or password",
		},
		{
			name:        "non-existent email",
			body:        `{"email":"ghost@example.com","password":"secret123"}`,
			wantStatus:  http.StatusUnauthorized,
			wantSuccess: false,
			wantMessage: "Invalid email or password",
		},
		{
			name:        "missing email",
			body:        `{"password":"secret123"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Email and password are required",
		},
		{
			name:        "missing password",
			body:        `{"email":"login@example.com"}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Email and password are required",
		},
		{
			name:        "empty body",
			body:        `{}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Email and password are required",
		},
		{
			name:        "malformed JSON",
			body:        `{invalid}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequest(t, http.MethodPost, "/api/v1/auth/login", tt.body)
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

func TestLoginResponsePayload(t *testing.T) {
	// Seed user
	makeRequest(t, http.MethodPost, "/api/v1/auth/register",
		`{"name":"Payload User","email":"payload@example.com","password":"secret123"}`,
	)

	tests := []struct {
		name           string
		body           string
		wantStatus     int
		checkPayload   bool
		wantEmail      string
		wantName       string
		wantNoPassword bool
	}{
		{
			name:           "response contains user_id, name, email",
			body:           `{"email":"payload@example.com","password":"secret123"}`,
			wantStatus:     http.StatusOK,
			checkPayload:   true,
			wantEmail:      "payload@example.com",
			wantName:       "Payload User",
			wantNoPassword: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequest(t, http.MethodPost, "/api/v1/auth/login", tt.body)
			body := parseResponse(t, w)

			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}

			if !tt.checkPayload {
				return
			}

			data, ok := body["data"].(map[string]interface{})
			if !ok {
				t.Fatal("expected data field in response")
			}
			if data["user_id"] == nil {
				t.Error("expected user_id in data")
			}
			if data["email"] != tt.wantEmail {
				t.Errorf("email: want %q, got %v", tt.wantEmail, data["email"])
			}
			if data["name"] != tt.wantName {
				t.Errorf("name: want %q, got %v", tt.wantName, data["name"])
			}
			if tt.wantNoPassword && data["password"] != nil {
				t.Error("password must never be returned in login response")
			}
		})
	}
}
