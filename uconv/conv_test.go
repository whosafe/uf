package uconv

import (
	"testing"
	"time"
)

func TestToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int
		wantErr bool
	}{
		{"int", 123, 123, false},
		{"string", "123", 123, false},
		{"string_zero", "0", 0, false},
		{"float", 123.45, 123, false},
		{"nil", nil, 0, true},
		{"invalid", "abc", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"string", "hello", "hello"},
		{"int", 123, "123"},
		{"nil", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.input); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    bool
		wantErr bool
	}{
		{"bool_true", true, true, false},
		{"bool_false", false, false, false},
		{"int_1", 1, true, false},
		{"int_0", 0, false, false},
		{"str_true", "true", true, false},
		{"str_false", "false", false, false},
		{"str_1", "1", true, false},
		{"str_0", "0", false, false},
		{"nil", nil, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    time.Duration
		wantErr bool
	}{
		{"duration", time.Second, time.Second, false},
		{"int", int(time.Second), time.Second, false}, // int is interpreted as ns if calling Duration(int) but wait logic in ToDuration
		// ToDuration implementation: case int: return time.Duration(val), nil. And time.Duration is int64 nanoseconds.
		// So passing 1 means 1 nanosecond.
		{"int_ns", 1, 1 * time.Nanosecond, false},
		{"string", "10s", 10 * time.Second, false},
		{"string_ms", "100ms", 100 * time.Millisecond, false},
		{"invalid", "abc", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		{"float", 123.456, 123.456, false},
		{"int", 123, 123.0, false},
		{"string", "123.456", 123.456, false},
		{"invalid", "abc", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
