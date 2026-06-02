package controllers

// RegisterRequest is the payload used to create a new user account.
// Password must be at least 6 characters and email must use a valid format.
type RegisterRequest struct {
	// Full name of the user.
	// required: true
	// example: Jane Doe
	Name string `json:"name"`

	// Login email address.
	// required: true
	// format: email
	// example: jane.doe@example.com
	Email string `json:"email"`

	// Account password.
	// required: true
	// minLength: 6
	// example: secret123
	Password string `json:"password"`
}

// LoginRequest is the payload used for user authentication.
type LoginRequest struct {
	// Login email address.
	// required: true
	// format: email
	// example: jane.doe@example.com
	Email string `json:"email"`

	// Account password.
	// required: true
	// example: secret123
	Password string `json:"password"`
}

// LoginResponseData contains the authenticated user information returned after login.
type LoginResponseData struct {
	// Unique user identifier.
	// example: 123
	UserID int `json:"user_id"`

	// User display name.
	// example: Jane Doe
	Name string `json:"name"`

	// Registered user email.
	// example: jane.doe@example.com
	Email string `json:"email"`
}

// ResponseEnvelope is the standard API response envelope.
type ResponseEnvelope struct {
	// Indicates whether the request succeeded.
	// example: true
	Success bool `json:"success"`

	// HTTP status code returned by the API.
	// example: 200
	Code int `json:"code"`

	// Human-readable description of the result.
	// example: OK
	Message string `json:"message"`
}

// LoginResponse is returned after a successful login.
type LoginResponse struct {
	ResponseEnvelope
	Data LoginResponseData `json:"data"`
}

// RegisterResponse is returned after successful user registration.
type RegisterResponse struct {
	ResponseEnvelope
}

// ExpenseCreateRequest is the payload for creating a new expense.
type ExpenseCreateRequest struct {
	// Expense title.
	// required: true
	// example: Team lunch
	Title string `json:"title"`

	// Expense amount. Must be greater than zero.
	// required: true
	// example: 350.50
	Amount float64 `json:"amount"`

	// Expense category. Must be one of the allowed categories.
	// required: true
	// enum: Food,Transport,Housing,Entertainment,Shopping,Healthcare,Education,Utilities,Other
	// example: Food
	Category string `json:"category"`

	// Optional note for the expense.
	// example: Lunch with the product team
	Note string `json:"note,omitempty"`

	// Expense date in YYYY-MM-DD format.
	// required: true
	// example: 2025-06-10
	ExpenseDate string `json:"expense_date"`
}

// ExpenseUpdateRequest is the payload for partial expense updates.
// Only provided fields are updated; omitted fields remain unchanged.
type ExpenseUpdateRequest struct {
	// Optional updated expense title.
	// example: Office snacks
	Title *string `json:"title,omitempty"`

	// Optional updated expense amount.
	// example: 42.99
	Amount *float64 `json:"amount,omitempty"`

	// Optional updated expense category.
	// enum: Food,Transport,Housing,Entertainment,Shopping,Healthcare,Education,Utilities,Other
	// example: Food
	Category *string `json:"category,omitempty"`

	// Optional updated note.
	// example: Updated note for this expense
	Note *string `json:"note,omitempty"`

	// Optional updated expense date in YYYY-MM-DD format.
	// example: 2025-06-15
	ExpenseDate *string `json:"expense_date,omitempty"`
}

// ExpenseResponse is the detailed representation of an expense record.
type ExpenseResponse struct {
	// Expense resource identifier.
	// example: 42
	ID int `json:"id"`

	// User identifier who owns the expense.
	// example: 123
	UserID int `json:"user_id"`

	// Expense title.
	// example: Team lunch
	Title string `json:"title"`

	// Expense amount.
	// example: 350.50
	Amount float64 `json:"amount"`

	// Expense category.
	// example: Food
	Category string `json:"category"`

	// Optional note.
	// example: Lunch with the product team
	Note string `json:"note,omitempty"`

	// Expense date in YYYY-MM-DD format.
	// example: 2025-06-10
	ExpenseDate string `json:"expense_date"`

	// Creation timestamp in RFC3339 format.
	// example: 2025-06-10T14:45:00Z
	CreatedAt string `json:"created_at"`
}

// ExpenseResponseWrapper wraps a single expense result.
type ExpenseResponseWrapper struct {
	ResponseEnvelope
	Data ExpenseResponse `json:"data"`
}

// ExpenseListResponse wraps an expense list result.
type ExpenseListResponse struct {
	ResponseEnvelope
	Data []ExpenseResponse `json:"data"`
}

// SummaryCategoryResponse represents a single category summary item.
type SummaryCategoryResponse struct {
	// Category name.
	// example: Food
	Category string `json:"category"`

	// Total amount for this category.
	// example: 350.50
	Total float64 `json:"total"`

	// Number of expenses in this category.
	// example: 3
	Count int `json:"count"`
}

// SummaryResponseData contains the returned summary data.
type SummaryResponseData struct {
	// Total amount for the selected date range.
	// example: 650.00
	TotalAmount float64 `json:"total_amount"`

	// Total count of matching expenses.
	// example: 5
	TotalCount int `json:"total_count"`

	// Applied start date filter.
	// example: 2025-06-01
	DateFrom string `json:"date_from,omitempty"`

	// Applied end date filter.
	// example: 2025-06-30
	DateTo string `json:"date_to,omitempty"`

	// Aggregated totals grouped by category.
	ByCategory []SummaryCategoryResponse `json:"by_category"`
}

// SummaryResponse wraps the returned summary result.
type SummaryResponse struct {
	ResponseEnvelope
	Data SummaryResponseData `json:"data"`
}

// StandardResponse is used for responses without specialized data.
type StandardResponse struct {
	ResponseEnvelope
	Data interface{} `json:"data,omitempty"`
}

// ErrorResponse is used for all error responses.
type ErrorResponse struct {
	ResponseEnvelope
	Data interface{} `json:"data,omitempty"`
}
