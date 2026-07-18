package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testViewsDir = "../../views"

// testConnIP returns the remote IP Fiber sees on app.Test connections, probed
// with an empty trust list so proxy headers cannot influence the result.
func testConnIP(t *testing.T) string {
	t.Helper()
	app, err := newApp(appConfig{ViewsDir: testViewsDir})
	if err != nil {
		t.Fatalf("newApp failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/ip", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("probe request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading probe body: %v", err)
	}
	return string(body)
}

func TestTrustedProxyModel(t *testing.T) {
	connIP := testConnIP(t)

	tests := []struct {
		name           string
		trustedProxies []string
		forwardedFor   string
		want           string
	}{
		{
			name:           "spoofed header from untrusted peer is ignored",
			trustedProxies: []string{"192.0.2.1"},
			forwardedFor:   "203.0.113.9",
			want:           connIP,
		},
		{
			name:           "empty trust list ignores proxy headers",
			trustedProxies: nil,
			forwardedFor:   "203.0.113.9",
			want:           connIP,
		},
		{
			name:           "header honored when peer is a trusted proxy",
			trustedProxies: []string{connIP},
			forwardedFor:   "203.0.113.9",
			want:           "203.0.113.9",
		},
		{
			name:           "chain from trusted peer resolves to first valid ip",
			trustedProxies: []string{connIP},
			forwardedFor:   "203.0.113.9, 10.0.0.1",
			want:           "203.0.113.9",
		},
		{
			name:           "no header returns socket ip even when trusted",
			trustedProxies: []string{connIP},
			forwardedFor:   "",
			want:           connIP,
		},
		{
			name:           "invalid value from trusted peer falls back to socket ip",
			trustedProxies: []string{connIP},
			forwardedFor:   "not-an-ip",
			want:           connIP,
		},
		{
			name:           "invalid first entry is skipped for the next valid ip",
			trustedProxies: []string{connIP},
			forwardedFor:   "not-an-ip, 203.0.113.9",
			want:           "203.0.113.9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := newApp(appConfig{
				TrustedProxies: tt.trustedProxies,
				ProxyHeader:    "X-Forwarded-For",
				ViewsDir:       testViewsDir,
			})
			if err != nil {
				t.Fatalf("newApp failed: %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/ip", nil)
			if tt.forwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.forwardedFor)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("reading body: %v", err)
			}

			if got := string(body); got != tt.want {
				t.Errorf("GET /ip = %q, want %q", got, tt.want)
			}
		})
	}
}

// The app.Test connection peer (0.0.0.0) sits outside the loopback-only
// default trust list, so it behaves like any external client here.
func TestDefaultTrustListIgnoresSpoofedHeader(t *testing.T) {
	connIP := testConnIP(t)

	app, err := newApp(appConfig{
		TrustedProxies: splitAndTrim(defaultTrustedProxies),
		ProxyHeader:    "X-Forwarded-For",
		ViewsDir:       testViewsDir,
	})
	if err != nil {
		t.Fatalf("newApp failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/ip", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.9")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	if got := string(body); got != connIP {
		t.Errorf("GET /ip with default trust list = %q, want socket ip %q", got, connIP)
	}
}

func TestRootRoute(t *testing.T) {
	connIP := testConnIP(t)

	tests := []struct {
		name            string
		userAgent       string
		wantContentType string
		wantBodyPart    string
	}{
		{
			name:            "curl gets the bare ip in plain text",
			userAgent:       "curl/8.6.0",
			wantContentType: "text/plain",
			wantBodyPart:    connIP,
		},
		{
			name:            "browser gets the html page",
			userAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			wantContentType: "text/html",
			wantBodyPart:    "What Is My IP Address?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := newApp(appConfig{ViewsDir: testViewsDir})
			if err != nil {
				t.Fatalf("newApp failed: %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("User-Agent", tt.userAgent)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("GET / status = %d, want %d", resp.StatusCode, http.StatusOK)
			}
			if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, tt.wantContentType) {
				t.Errorf("Content-Type = %q, want it to contain %q", ct, tt.wantContentType)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("reading body: %v", err)
			}
			if !strings.Contains(string(body), tt.wantBodyPart) {
				t.Errorf("body does not contain %q; body = %q", tt.wantBodyPart, truncate(string(body), 200))
			}
		})
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func TestNewAppFailsOnMissingViewsDir(t *testing.T) {
	_, err := newApp(appConfig{ViewsDir: "./does-not-exist"})
	if err == nil {
		t.Fatal("newApp with a missing views directory returned nil error, want template load failure")
	}
}
