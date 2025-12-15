package rule_test

import (
	"testing"

	"iutime.com/utime/uf/uvalidator/rule"
)

func TestDate(t *testing.T) {
	r := rule.NewDate()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "2024-01-01", true},
		{"invalid - wrong format", "2024/01/01", false},
		{"invalid - not a date", "invalid", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Date.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDateTime(t *testing.T) {
	r := rule.NewDateTime()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "2024-01-01 12:00:00", true},
		{"invalid - wrong format", "2024/01/01 12:00:00", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("DateTime.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDateBefore(t *testing.T) {
	r := rule.NewDateBefore("2024-12-31")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "2024-06-01", true},
		{"invalid - after", "2025-01-01", false},
		{"invalid - equal", "2024-12-31", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("DateBefore.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDateAfter(t *testing.T) {
	r := rule.NewDateAfter("2024-01-01")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "2024-06-01", true},
		{"invalid - before", "2023-12-31", false},
		{"invalid - equal", "2024-01-01", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("DateAfter.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDateBetween(t *testing.T) {
	r := rule.NewDateBetween("2024-01-01", "2024-12-31")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "2024-06-15", true},
		{"valid - start boundary", "2024-01-01", true},
		{"valid - end boundary", "2024-12-31", true},
		{"invalid - before", "2023-12-31", false},
		{"invalid - after", "2025-01-01", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("DateBetween.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
