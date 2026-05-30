package utils

import (
	"regexp"
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
