package rule_test

import (
	"testing"

	"github.com/whosafe/uf/uvalidator/rule"
)

func TestIDCard(t *testing.T) {
	r := rule.NewIDCard()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"invalid - wrong format", "123456789012345678", false},
		{"invalid - too short", "12345", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("IDCard.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestBankCard(t *testing.T) {
	r := rule.NewBankCard()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"invalid - wrong length", "123456", false},
		{"invalid - not digits", "abcd1234567890", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("BankCard.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestUnifiedSocialCreditCode(t *testing.T) {
	r := rule.NewUnifiedSocialCreditCode()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "91110000600037341L", true},
		{"invalid - wrong length", "1234567890", false},
		{"invalid - wrong format", "I1110000600037341L", false}, // I is not allowed
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("UnifiedSocialCreditCode.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestPostalCode(t *testing.T) {
	r := rule.NewPostalCode()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "100000", true},
		{"invalid - too short", "10000", false},
		{"invalid - too long", "1000000", false},
		{"invalid - not digits", "10000a", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("PostalCode.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestChineseName(t *testing.T) {
	r := rule.NewChineseName()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "张三", true},
		{"valid - with dot", "欧阳·娜娜", true},
		{"invalid - too short", "张", false},
		{"invalid - English", "John", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("ChineseName.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
