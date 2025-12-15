package rule_test

import (
	"testing"

	"iutime.com/utime/uf/uvalidator/rule"
)

func TestFileExtension(t *testing.T) {
	r := rule.NewFileExtension("jpg", "png", "gif")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid - jpg", "photo.jpg", true},
		{"valid - PNG (case insensitive)", "photo.PNG", true},
		{"invalid - pdf", "document.pdf", false},
		{"invalid - no extension", "file", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("FileExtension.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestMimeType(t *testing.T) {
	r := rule.NewMimeType("image/jpeg", "image/png")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid - jpeg", "image/jpeg", true},
		{"valid - png", "image/png", true},
		{"invalid - pdf", "application/pdf", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("MimeType.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestFileSize(t *testing.T) {
	r := rule.NewFileSize(1024, 1024*1024) // 1KB - 1MB

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - int", 102400, true},
		{"valid - int64", int64(102400), true},
		{"valid - min boundary", int64(1024), true},
		{"valid - max boundary", int64(1024 * 1024), true},
		{"invalid - too small", 500, false},
		{"invalid - too large", 2 * 1024 * 1024, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("FileSize.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestConfirmed(t *testing.T) {
	r := rule.NewConfirmed("password123")

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid", "password123", true},
		{"invalid", "different", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Confirmed.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDistinct(t *testing.T) {
	r := rule.NewDistinct("admin")

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid", "user", true},
		{"invalid", "admin", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Distinct.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNotIn(t *testing.T) {
	r := rule.NewNotIn("spam", "test", "demo")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "production", true},
		{"invalid", "test", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("NotIn.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestNullable(t *testing.T) {
	r := rule.NewNullable()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"any value", "anything", true},
		{"nil", nil, true},
		{"zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Nullable.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLen(t *testing.T) {
	r := rule.NewLen(5)

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "hello", true},
		{"invalid - too short", "hi", false},
		{"invalid - too long", "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Len.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
