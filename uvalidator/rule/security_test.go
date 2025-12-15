package rule_test

import (
	"testing"

	"github.com/whosafe/uf/uvalidator/rule"
)

func TestStrongPassword(t *testing.T) {
	r := rule.NewStrongPassword()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "Abc123!@#", true},
		{"invalid - no uppercase", "abc123!@#", false},
		{"invalid - no lowercase", "ABC123!@#", false},
		{"invalid - no digit", "Abcdef!@#", false},
		{"invalid - no special", "Abc123456", false},
		{"invalid - too short", "Abc1!@", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("StrongPassword.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNoHTML(t *testing.T) {
	r := rule.NewNoHTML()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "plain text", true},
		{"invalid - script tag", "<script>alert('xss')</script>", false},
		{"invalid - div tag", "<div>content</div>", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("NoHTML.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNoSQL(t *testing.T) {
	r := rule.NewNoSQL()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "safe text", true},
		{"invalid - SELECT", "SELECT * FROM users", false},
		{"invalid - DROP", "DROP TABLE users", false},
		{"invalid - mixed case", "SeLeCt * FROM users", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("NoSQL.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNoXSS(t *testing.T) {
	r := rule.NewNoXSS()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "normal text", true},
		{"invalid - javascript:", "javascript:alert(1)", false},
		{"invalid - script tag", "<script>alert(1)</script>", false},
		{"invalid - onerror", "<img onerror=alert(1)>", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("NoXSS.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
