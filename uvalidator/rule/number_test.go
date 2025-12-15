package rule_test

import (
	"testing"

	"iutime.com/utime/uf/uvalidator/rule"
)

func TestBetween(t *testing.T) {
	r := rule.NewBetween(10, 100)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - int", 50, true},
		{"valid - min boundary", 10, true},
		{"valid - max boundary", 100, true},
		{"invalid - too small", 5, false},
		{"invalid - too large", 150, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Between.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestPositive(t *testing.T) {
	r := rule.NewPositive()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - positive int", 10, true},
		{"valid - positive float", 3.14, true},
		{"invalid - zero", 0, false},
		{"invalid - negative", -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Positive.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNegative(t *testing.T) {
	r := rule.NewNegative()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - negative int", -10, true},
		{"valid - negative float", -3.14, true},
		{"invalid - zero", 0, false},
		{"invalid - positive", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Negative.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestInteger(t *testing.T) {
	r := rule.NewInteger()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - int", 42, true},
		{"valid - int64", int64(42), true},
		{"valid - float without decimal", 42.0, true},
		{"invalid - float with decimal", 3.14, false},
		{"valid - string int", "42", true},
		{"invalid - string float", "3.14", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Integer.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDecimal(t *testing.T) {
	r := rule.NewDecimal(2)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - 2 decimals", 3.14, true},
		{"valid - 1 decimal", 3.1, true},
		{"valid - string", "3.14", true},
		{"invalid - 3 decimals", "3.141", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Decimal.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestGt(t *testing.T) {
	r := rule.NewGt(10)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid", 15, true},
		{"invalid - equal", 10, false},
		{"invalid - less", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Gt.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestGte(t *testing.T) {
	r := rule.NewGte(10)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - greater", 15, true},
		{"valid - equal", 10, true},
		{"invalid - less", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Gte.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLt(t *testing.T) {
	r := rule.NewLt(10)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid", 5, true},
		{"invalid - equal", 10, false},
		{"invalid - greater", 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Lt.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLte(t *testing.T) {
	r := rule.NewLte(10)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - less", 5, true},
		{"valid - equal", 10, true},
		{"invalid - greater", 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Lte.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
