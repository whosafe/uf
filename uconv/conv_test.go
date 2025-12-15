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

func TestToTime(t *testing.T) {
	// 创建测试用的时间
	testTime := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	testTimeStr := "2024-01-15 10:30:45"
	testDateStr := "2024-01-15"
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		input   any
		want    time.Time
		wantErr bool
	}{
		{"time.Time", testTime, testTime, false},
		{"time_ptr", &testTime, testTime, false},
		{"unix_timestamp_int64", int64(1705315845), time.Unix(1705315845, 0), false},
		{"unix_timestamp_int", int(1705315845), time.Unix(1705315845, 0), false},
		{"rfc3339", "2024-01-15T10:30:45Z", time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC), false},
		{"datetime", testTimeStr, testTime, false},
		{"date", testDateStr, testDate, false},
		{"empty_string", "", time.Time{}, false},
		{"nil", nil, time.Time{}, true},
		{"invalid_string", "invalid", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToTimeDef(t *testing.T) {
	defaultTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	validTime := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	tests := []struct {
		name  string
		input any
		def   time.Time
		want  time.Time
	}{
		{"valid_string", "2024-01-15 10:30:45", defaultTime, validTime},
		{"invalid_string", "invalid", defaultTime, defaultTime},
		{"nil", nil, defaultTime, defaultTime},
		{"valid_time", validTime, defaultTime, validTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToTimeDef(tt.input, tt.def)
			if !got.Equal(tt.want) {
				t.Errorf("ToTimeDef() = %v, want %v", got, tt.want)
			}
		})
	}
}
