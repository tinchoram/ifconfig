package interactions

import (
	"errors"
	"testing"

	"ifconfig/pkg/models"
)

// fakeGateway is a hand-rolled test double. The compile-time assertion below
// guarantees it stays in sync with the Gateway interface.
type fakeGateway struct {
	headers  map[string]string
	ip       string
	method   string
	hostname string
}

var _ Gateway = (*fakeGateway)(nil)

func (f *fakeGateway) GetHeader(key string) string { return f.headers[key] }
func (f *fakeGateway) GetMethod() string           { return f.method }
func (f *fakeGateway) GetHostname() string         { return f.hostname }
func (f *fakeGateway) GetIP() string               { return f.ip }

func TestIPService_GetIPInfo(t *testing.T) {
	gateway := &fakeGateway{
		headers: map[string]string{
			"User-Agent":      "Mozilla/5.0",
			"Accept-Language": "es-ES,es;q=0.9",
			"Accept-Encoding": "gzip, deflate, br",
			"Accept":          "text/html,application/xhtml+xml",
			"Via":             "1.1 proxy",
			"Connection":      "keep-alive",
			"Keep-Alive":      "timeout=5, max=100",
			"Referer":         "https://example.com",
			"Cf-Ipcountry":    "AR",
			"Cf-Ipcity":       "Buenos Aires",
			"Cf-Region":       "Buenos Aires",
			"Cf-Postal-Code":  "1425",
			"Cf-Iplatitude":   "-34.6037",
			"Cf-Iplongitude":  "-58.3816",
			"Cf-Timezone":     "America/Argentina/Buenos_Aires",
			"Cf-Ipcontinent":  "SA",
			"X-Forwarded-For": "192.168.1.1",
		},
		ip:       "192.168.1.100",
		method:   "GET",
		hostname: "example.com",
	}

	svc := &Service{lookupAddr: func(ip string) ([]string, error) {
		if ip != "192.168.1.100" {
			t.Fatalf("reverse DNS looked up %q, want the client IP", ip)
		}
		return []string{"host.example.com."}, nil
	}}

	got := svc.GetIPInfo(gateway)

	want := models.IPInfo{
		IPAddr:     "192.168.1.100",
		RemoteHost: "host.example.com",
		UserAgent:  "Mozilla/5.0",
		Language:   "es-ES,es;q=0.9",
		Method:     "GET",
		Encoding:   "gzip, deflate, br",
		Mime:       "text/html,application/xhtml+xml",
		Via:        "1.1 proxy",
		Forwarded:  "192.168.1.1",
		Connection: "keep-alive",
		KeepAlive:  "timeout=5, max=100",
		Referer:    "https://example.com",
		Country:    "AR",
		City:       "Buenos Aires",
		Region:     "Buenos Aires",
		PostalCode: "1425",
		Latitude:   "-34.6037",
		Longitude:  "-58.3816",
		Timezone:   "America/Argentina/Buenos_Aires",
		Continent:  "SA",
		Host:       "example.com",
	}

	if got != want {
		t.Errorf("GetIPInfo() = %+v, want %+v", got, want)
	}
}

func TestIPService_ReverseDNS(t *testing.T) {
	tests := []struct {
		name   string
		lookup func(string) ([]string, error)
		want   string
	}{
		{
			name:   "resolves and strips the trailing dot",
			lookup: func(string) ([]string, error) { return []string{"one.example.com.", "two.example.com."}, nil },
			want:   "one.example.com",
		},
		{
			name:   "empty result yields empty host",
			lookup: func(string) ([]string, error) { return nil, nil },
			want:   "",
		},
		{
			name:   "lookup error yields empty host",
			lookup: func(string) ([]string, error) { return nil, errors.New("no PTR record") },
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{lookupAddr: tt.lookup}
			got := svc.GetIPInfo(&fakeGateway{ip: "203.0.113.5", headers: map[string]string{}}).RemoteHost
			if got != tt.want {
				t.Errorf("RemoteHost = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIPService_GeoHeadersAbsentDegradeToEmpty(t *testing.T) {
	svc := &Service{lookupAddr: func(string) ([]string, error) { return nil, nil }}

	got := svc.GetIPInfo(&fakeGateway{ip: "203.0.113.5", headers: map[string]string{}})

	if got.City != "" || got.Region != "" || got.Country != "" || got.PostalCode != "" ||
		got.Latitude != "" || got.Longitude != "" || got.Timezone != "" || got.Continent != "" {
		t.Errorf("expected empty geo fields when Cloudflare headers are absent, got %+v", got)
	}
}
