package rule_test

import (
	"testing"

	"iutime.com/utime/uf/uvalidator/rule"
)

func TestUUID(t *testing.T) {
	r := rule.NewUUID()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", true},
		{"invalid - wrong format", "not-a-uuid", false},
		{"invalid - wrong version", "550e8400-e29b-31d4-a716-446655440000", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("UUID.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	r := rule.NewJSON()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid object", `{"name":"test"}`, true},
		{"valid array", `[1,2,3]`, true},
		{"invalid", `{invalid}`, false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("JSON.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestBase64(t *testing.T) {
	r := rule.NewBase64()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "SGVsbG8gV29ybGQ=", true},
		{"invalid", "not-base64!@#", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Base64.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLowercase(t *testing.T) {
	r := rule.NewLowercase()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "hello", true},
		{"invalid - uppercase", "Hello", false},
		{"invalid - mixed", "HeLLo", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Lowercase.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestUppercase(t *testing.T) {
	r := rule.NewUppercase()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "HELLO", true},
		{"invalid - lowercase", "Hello", false},
		{"invalid - mixed", "HeLLo", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Uppercase.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestASCII(t *testing.T) {
	r := rule.NewASCII()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "hello123", true},
		{"invalid - chinese", "你好", false},
		{"invalid - emoji", "hello😀", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("ASCII.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNotBlank(t *testing.T) {
	r := rule.NewNotBlank()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "hello", true},
		{"invalid - empty", "", false},
		{"invalid - spaces", "   ", false},
		{"invalid - tabs", "\t\t", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("NotBlank.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	r := rule.NewContains("world")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "hello world", true},
		{"invalid", "hello", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Contains.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	r := rule.NewStartsWith("hello")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "hello world", true},
		{"invalid", "world hello", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("StartsWith.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	r := rule.NewEndsWith("world")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "hello world", true},
		{"invalid", "world hello", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("EndsWith.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestRegex(t *testing.T) {
	r := rule.NewRegex(`^\d{3}-\d{4}$`)

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "123-4567", true},
		{"invalid", "1234567", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Regex.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
