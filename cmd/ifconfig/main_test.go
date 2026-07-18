package main

import (
	"reflect"
	"testing"
)

func TestEnvOr(t *testing.T) {
	t.Run("returns value when variable is set", func(t *testing.T) {
		t.Setenv("IFCONFIG_TEST_VAR", "custom")
		if got := envOr("IFCONFIG_TEST_VAR", "fallback"); got != "custom" {
			t.Errorf("envOr() = %q, want %q", got, "custom")
		}
	})

	t.Run("returns empty string when variable is set but empty", func(t *testing.T) {
		t.Setenv("IFCONFIG_TEST_VAR", "")
		if got := envOr("IFCONFIG_TEST_VAR", "fallback"); got != "" {
			t.Errorf("envOr() = %q, want empty string", got)
		}
	})

	t.Run("returns fallback when variable is unset", func(t *testing.T) {
		if got := envOr("IFCONFIG_TEST_UNSET_VAR", "fallback"); got != "fallback" {
			t.Errorf("envOr() = %q, want %q", got, "fallback")
		}
	})
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{name: "empty string yields nil", in: "", want: nil},
		{name: "single value", in: "127.0.0.1", want: []string{"127.0.0.1"}},
		{name: "trims surrounding spaces", in: " 127.0.0.1 , ::1 ", want: []string{"127.0.0.1", "::1"}},
		{name: "drops empty entries from trailing comma", in: "127.0.0.1,", want: []string{"127.0.0.1"}},
		{name: "drops whitespace-only entries", in: "127.0.0.1, ,::1", want: []string{"127.0.0.1", "::1"}},
		{name: "whitespace-only input yields nil", in: " , , ", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitAndTrim(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitAndTrim(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestListenAddr(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    string
		want    string
		wantErr bool
	}{
		{name: "default port on all interfaces", host: "", port: "3000", want: ":3000"},
		{name: "explicit host and port", host: "127.0.0.1", port: "8080", want: "127.0.0.1:8080"},
		{name: "non-numeric port is rejected", host: "", port: "abc", wantErr: true},
		{name: "empty port is rejected", host: "", port: "", wantErr: true},
		{name: "port zero is rejected", host: "", port: "0", wantErr: true},
		{name: "port above 65535 is rejected", host: "", port: "70000", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listenAddr(tt.host, tt.port)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("listenAddr(%q, %q) returned nil error, want error", tt.host, tt.port)
				}
				return
			}
			if err != nil {
				t.Fatalf("listenAddr(%q, %q) returned error: %v", tt.host, tt.port, err)
			}
			if got != tt.want {
				t.Errorf("listenAddr(%q, %q) = %q, want %q", tt.host, tt.port, got, tt.want)
			}
		})
	}
}
