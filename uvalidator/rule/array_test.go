package rule_test

import (
	"testing"

	"github.com/whosafe/uf/uvalidator/rule"
)

func TestUnique(t *testing.T) {
	r := rule.NewUnique()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - string array", []string{"a", "b", "c"}, true},
		{"invalid - duplicate strings", []string{"a", "b", "a"}, false},
		{"valid - int array", []int{1, 2, 3}, true},
		{"invalid - duplicate ints", []int{1, 2, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Unique.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestArrayMin(t *testing.T) {
	r := rule.NewArrayMin(2)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid", []string{"a", "b", "c"}, true},
		{"valid - exactly min", []string{"a", "b"}, true},
		{"invalid", []string{"a"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("ArrayMin.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestArrayMax(t *testing.T) {
	r := rule.NewArrayMax(3)

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid", []string{"a", "b"}, true},
		{"valid - exactly max", []string{"a", "b", "c"}, true},
		{"invalid", []string{"a", "b", "c", "d"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("ArrayMax.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestArrayContains(t *testing.T) {
	r := rule.NewArrayContains("admin")

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid", []string{"admin", "user"}, true},
		{"invalid", []string{"user", "guest"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("ArrayContains.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	r := rule.NewOneOf("red", "green", "blue")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "red", true},
		{"invalid", "yellow", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("OneOf.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestOneOfInt(t *testing.T) {
	r := rule.NewOneOfInt(1, 2, 3)

	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"valid", 2, true},
		{"invalid", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("OneOfInt.Validate(%d) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
