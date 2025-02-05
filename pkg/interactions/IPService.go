package interactions

import (
	"ifconfig/pkg/gateways"
	"ifconfig/pkg/models"
	"strings"
)

type IPService interface {
	GetIPInfo(gateways.Gateway) models.IPInfo
	GetRealIP(gateways.Gateway) string
}

type ipService struct{}

func NewIPService() IPService {
	return &ipService{}
}

func (s *ipService) GetIPInfo(gateway gateways.Gateway) models.IPInfo {
	return models.IPInfo{
		IPAddr:     s.GetRealIP(gateway),
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

func (s *ipService) GetRealIP(gateway gateways.Gateway) string {
	if cfIP := gateway.GetHeader("Cf-Connecting-Ip"); cfIP != "" {
		return cfIP
	}
	if realIP := gateway.GetHeader("X-Real-Ip"); realIP != "" {
		return realIP
	}
	if forwardedFor := gateway.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	return gateway.GetIP()
}
