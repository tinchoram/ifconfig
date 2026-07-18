package interactions

import (
	"testing"

	"ifconfig/pkg/models"
)

// fakeGateway is a hand-rolled test double. The compile-time assertion below
// guarantees it stays in sync with the Gateway interface.
type fakeGateway struct {
	headers  map[string]string
	ip       string
	port     string
	method   string
	hostname string
}

var _ Gateway = (*fakeGateway)(nil)

func (f *fakeGateway) GetHeader(key string) string { return f.headers[key] }
func (f *fakeGateway) GetPort() string             { return f.port }
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
			"X-Forwarded-For": "192.168.1.1",
		},
		ip:       "192.168.1.100",
		port:     "443",
		method:   "GET",
		hostname: "example.com",
	}

	got := NewService().GetIPInfo(gateway)

	want := models.IPInfo{
		IPAddr:     "192.168.1.100",
		RemoteHost: "unavailable",
		UserAgent:  "Mozilla/5.0",
		Port:       "443",
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
		Host:       "example.com",
	}

	if got != want {
		t.Errorf("GetIPInfo() = %+v, want %+v", got, want)
	}
}
