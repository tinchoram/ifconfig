package interactions

import (
	"ifconfig/pkg/models"
)

// Gateway is the port the IP service requires from the transport layer.
type Gateway interface {
	GetHeader(key string) string
	GetPort() string
	GetMethod() string
	GetHostname() string
	GetIP() string
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// GetIPInfo assembles the request snapshot exposed by the echo endpoints.
// Client IP resolution (proxy headers, trust checks) is a transport concern
// handled by the gateway; GetIP already returns the resolved client IP.
func (s *Service) GetIPInfo(gateway Gateway) models.IPInfo {
	return models.IPInfo{
		IPAddr:     gateway.GetIP(),
		RemoteHost: "unavailable",
		UserAgent:  gateway.GetHeader("User-Agent"),
		Port:       gateway.GetPort(),
		Language:   gateway.GetHeader("Accept-Language"),
		Method:     gateway.GetMethod(),
		Encoding:   gateway.GetHeader("Accept-Encoding"),
		Mime:       gateway.GetHeader("Accept"),
		Via:        gateway.GetHeader("Via"),
		Forwarded:  gateway.GetHeader("X-Forwarded-For"),
		Connection: gateway.GetHeader("Connection"),
		KeepAlive:  gateway.GetHeader("Keep-Alive"),
		Referer:    gateway.GetHeader("Referer"),
		Country:    gateway.GetHeader("Cf-Ipcountry"),
		Host:       gateway.GetHostname(),
	}
}
