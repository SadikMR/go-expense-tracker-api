package utils

import "testing"

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{name: "valid email", arg: "user@example.com", want: true},
		{name: "valid with subdomain", arg: "user@mail.example.co.uk", want: true},
		{name: "valid with plus", arg: "user+tag@example.com", want: true},
		{name: "invalid no @", arg: "userexample.com", want: false},
		{name: "invalid no domain", arg: "user@", want: false},
		{name: "invalid no local", arg: "@example.com", want: false},
		{name: "invalid double @", arg: "user@@example.com", want: false},
		{name: "empty string", arg: "", want: false},
		{name: "invalid format", arg: "user@.com", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.arg); got != tt.want {
				t.Fatalf("ValidateEmail(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestValidateMinLength(t *testing.T) {
	tests := []struct {
		name string
		str  string
		min  int
		want bool
	}{
		{name: "string meets minimum", str: "hello", min: 5, want: true},
		{name: "string exceeds minimum", str: "hello", min: 3, want: true},
		{name: "string below minimum", str: "hi", min: 3, want: false},
		{name: "minimum zero", str: "", min: 0, want: true},
		{name: "whitespace only fails", str: "   ", min: 5, want: false},
		{name: "whitespace trimmed", str: "   ab   ", min: 2, want: true},
		{name: "empty string", str: "", min: 1, want: false},
		{name: "negative minimum", str: "test", min: -1, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateMinLength(tt.str, tt.min); got != tt.want {
				t.Fatalf("ValidateMinLength(%q, %d) = %v, want %v", tt.str, tt.min, got, tt.want)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "all present", args: []string{"a", "b"}, want: true},
		{name: "one empty", args: []string{"a", ""}, want: false},
		{name: "all whitespace", args: []string{"  ", "x"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateRequired(tt.args...)
			if got != tt.want {
				t.Fatalf("ValidateRequired(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{name: "valid date", arg: "2025-06-01", want: true},
		{name: "invalid format", arg: "06/01/2025", want: false},
		{name: "missing parts", arg: "2025-06", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDate(tt.arg); got != tt.want {
				t.Fatalf("ValidateDate(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestValidatePositiveAmount(t *testing.T) {
	tests := []struct {
		name string
		arg  float64
		want bool
	}{
		{name: "positive", arg: 12.5, want: true},
		{name: "zero", arg: 0, want: false},
		{name: "negative", arg: -3.25, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePositiveAmount(tt.arg); got != tt.want {
				t.Fatalf("ValidatePositiveAmount(%v) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestValidateCategory(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{name: "allowed category", arg: "Food", want: true},
		{name: "unknown category", arg: "Unknown", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateCategory(tt.arg); got != tt.want {
				t.Fatalf("ValidateCategory(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		want   int
		wantOk bool
	}{
		{name: "valid id", arg: "42", want: 42, wantOk: true},
		{name: "zero id", arg: "0", wantOk: false},
		{name: "negative id", arg: "-1", wantOk: false},
		{name: "not a number", arg: "abc", wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseID(tt.arg)
			if got != tt.want || ok != tt.wantOk {
				t.Fatalf("ParseID(%q) = (%d, %v), want (%d, %v)", tt.arg, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestParseLimit(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want int
	}{
		{name: "valid limit", arg: "10", want: 10},
		{name: "missing limit", arg: "", want: 0},
		{name: "invalid limit", arg: "abc", want: 0},
		{name: "negative limit", arg: "-5", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLimit(tt.arg, 0); got != tt.want {
				t.Fatalf("ParseLimit(%q, 0) = %d, want %d", tt.arg, got, tt.want)
			}
		})
	}
}
