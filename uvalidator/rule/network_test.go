package rule_test

import (
	"testing"

	"github.com/whosafe/uf/uvalidator/rule"
)

func TestIP(t *testing.T) {
	r := rule.NewIP()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid IPv4", "192.168.1.1", true},
		{"valid IPv6", "::1", true},
		{"invalid", "999.999.999.999", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("IP.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIPv4(t *testing.T) {
	r := rule.NewIPv4()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "192.168.1.1", true},
		{"invalid - IPv6", "::1", false},
		{"invalid - malformed", "999.999.999.999", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("IPv4.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIPv6(t *testing.T) {
	r := rule.NewIPv6()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "::1", true},
		{"valid - full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"invalid - IPv4", "192.168.1.1", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("IPv6.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestMAC(t *testing.T) {
	r := rule.NewMAC()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid - colon", "00:1B:44:11:3A:B7", true},
		{"valid - hyphen", "00-1B-44-11-3A-B7", true},
		{"invalid", "invalid", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("MAC.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDomain(t *testing.T) {
	r := rule.NewDomain()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid", "example.com", true},
		{"valid - subdomain", "www.example.com", true},
		{"invalid - no TLD", "example", false},
		{"invalid - starts with hyphen", "-example.com", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Domain.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestPort(t *testing.T) {
	r := rule.NewPort()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid - int", 8080, true},
		{"valid - string", "8080", true},
		{"valid - min", 1, true},
		{"valid - max", 65535, true},
		{"invalid - zero", 0, false},
		{"invalid - too large", 99999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("Port.Validate(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestURL(t *testing.T) {
	r := rule.NewURL()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid http", "http://example.com", true},
		{"valid https", "https://example.com", true},
		{"invalid - no protocol", "example.com", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Validate(tt.value); got != tt.want {
				t.Errorf("URL.Validate(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
