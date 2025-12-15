package rule_test

import (
	"testing"

	"github.com/whosafe/uf/uvalidator/rule"
)

func TestRequired(t *testing.T) {
	r := rule.NewRequired()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"empty string", "", false},
		{"non-empty string", "hello", true},
		{"zero int", 0, false},
		{"non-zero int", 42, true},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Required.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	r := rule.NewEmail()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"invalid - no @", "testexample.com", false},
		{"invalid - no domain", "test@", false},
		{"invalid - no local", "@example.com", false},
		{"empty string", "", true}, // empty is valid (use Required for non-empty)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Email.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestPhone(t *testing.T) {
	r := rule.NewPhone()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid phone", "13812345678", true},
		{"valid phone - 15x", "15987654321", true},
		{"invalid - too short", "1381234567", false},
		{"invalid - wrong prefix", "12812345678", false},
		{"invalid - not digits", "138abcd5678", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Phone.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	r := rule.NewMin(3)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"string - valid", "hello", true},
		{"string - invalid", "hi", false},
		{"int - valid", 5, true},
		{"int - invalid", 2, false},
		{"int - equal", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Min.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	r := rule.NewMax(10)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"string - valid", "hello", true},
		{"string - invalid", "hello world!", false},
		{"int - valid", 5, true},
		{"int - invalid", 15, false},
		{"int - equal", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Max.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestAlpha(t *testing.T) {
	r := rule.NewAlpha()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid - lowercase", "hello", true},
		{"valid - uppercase", "HELLO", true},
		{"valid - mixed", "HelloWorld", true},
		{"invalid - with digits", "hello123", false},
		{"invalid - with spaces", "hello world", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Alpha.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestAlphanum(t *testing.T) {
	r := rule.NewAlphanum()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid - letters only", "hello", true},
		{"valid - digits only", "12345", true},
		{"valid - mixed", "hello123", true},
		{"invalid - with spaces", "hello 123", false},
		{"invalid - with symbols", "hello@123", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Alphanum.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNumeric(t *testing.T) {
	r := rule.NewNumeric()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid - digits", "12345", true},
		{"invalid - letters", "abc", false},
		{"invalid - mixed", "123abc", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Numeric.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
