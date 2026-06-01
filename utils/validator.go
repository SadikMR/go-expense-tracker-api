package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ValidateRequired returns false if any provided field is empty or whitespace only.
func ValidateRequired(fields ...string) bool {
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			return false
		}
	}
	return true
}

// ValidateEmail returns true if the email matches a standard format.
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ValidateMinLength returns true if the string meets the minimum character length.
func ValidateMinLength(s string, min int) bool {
	return len(strings.TrimSpace(s)) >= min
}

// ValidateDate returns true if the string matches YYYY-MM-DD format.
func ValidateDate(s string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(s)
}

// ValidatePositiveAmount returns true if the amount is greater than zero.
func ValidatePositiveAmount(amount float64) bool {
	return amount > 0
}

// ValidateCategory returns true if the category is in the allowed list.
func ValidateCategory(cat string) bool {
	for _, allowed := range AllowedCategories {
		if allowed == cat {
			return true
		}
	}
	return false
}

// AllowedCategories is the canonical list of valid expense categories.
var AllowedCategories = []string{
	"Food", "Transport", "Housing", "Entertainment",
	"Shopping", "Healthcare", "Education", "Utilities", "Other",
}

// ParseLimit safely parses a limit query param with a default fallback.
func ParseLimit(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}

// ParseID safely parses a string to a positive integer ID.
// Returns 0, false when the string is not a valid positive integer.
func ParseID(s string) (int, bool) {
	id, err := strconv.Atoi(s)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
