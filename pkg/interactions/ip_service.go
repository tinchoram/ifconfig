package interactions

import (
	"context"
	"net"
	"strings"
	"time"

	"ifconfig/pkg/models"
)

// Gateway is the port the IP service requires from the transport layer.
type Gateway interface {
	GetHeader(key string) string
	GetMethod() string
	GetHostname() string
	GetIP() string
}

// Service assembles the request snapshot exposed by the echo endpoints.
type Service struct {
	// lookupAddr resolves an IP to its PTR hostnames. Injected so tests stay
	// hermetic and never touch the network.
	lookupAddr func(ip string) ([]string, error)
}

func NewService() *Service {
	return &Service{lookupAddr: defaultLookupAddr}
}

// defaultLookupAddr does a reverse DNS lookup bounded by a short timeout so a
// slow or missing PTR never stalls a page render.
func defaultLookupAddr(ip string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	return net.DefaultResolver.LookupAddr(ctx, ip)
}

// GetIPInfo builds the request snapshot. Client IP resolution (proxy headers,
// trust checks) is a transport concern already handled by the gateway; the
// geolocation fields come from Cloudflare's visitor-location headers and are
// empty when that managed transform is disabled.
func (s *Service) GetIPInfo(gateway Gateway) models.IPInfo {
	ip := gateway.GetIP()
	return models.IPInfo{
		IPAddr:     ip,
		RemoteHost: s.reverseDNS(ip),
		UserAgent:  gateway.GetHeader("User-Agent"),
		Language:   gateway.GetHeader("Accept-Language"),
		Method:     gateway.GetMethod(),
		Encoding:   gateway.GetHeader("Accept-Encoding"),
		Mime:       gateway.GetHeader("Accept"),
		Via:        gateway.GetHeader("Via"),
		Forwarded:  gateway.GetHeader("X-Forwarded-For"),
		Connection: gateway.GetHeader("Connection"),
		KeepAlive:  gateway.GetHeader("Keep-Alive"),
		Referer:    gateway.GetHeader("Referer"),
		City:       gateway.GetHeader("Cf-Ipcity"),
		Region:     gateway.GetHeader("Cf-Region"),
		Country:    gateway.GetHeader("Cf-Ipcountry"),
		PostalCode: gateway.GetHeader("Cf-Postal-Code"),
		Latitude:   gateway.GetHeader("Cf-Iplatitude"),
		Longitude:  gateway.GetHeader("Cf-Iplongitude"),
		Timezone:   gateway.GetHeader("Cf-Timezone"),
		Continent:  gateway.GetHeader("Cf-Ipcontinent"),
		Host:       gateway.GetHostname(),
	}
}

// reverseDNS returns the first PTR hostname for ip, or "" when it cannot be
// resolved. Unspecified and loopback addresses are skipped (no useful PTR).
func (s *Service) reverseDNS(ip string) string {
	parsed := net.ParseIP(ip)
	if s.lookupAddr == nil || parsed == nil || parsed.IsUnspecified() || parsed.IsLoopback() {
		return ""
	}

	names, err := s.lookupAddr(ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	return strings.TrimSuffix(names[0], ".")
}
