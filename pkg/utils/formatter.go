package utils

import "ifconfig/pkg/models"

// FormatPlainTextResponse
func FormatPlainTextResponse(info models.IPInfo) string {
	return "ip_addr: " + info.IPAddr + "\n" +
		"remote_host: " + info.RemoteHost + "\n" +
		"user_agent: " + info.UserAgent + "\n" +
		"port: " + info.Port + "\n" +
		"language: " + info.Language + "\n" +
		"method: " + info.Method + "\n" +
		"encoding: " + info.Encoding + "\n" +
		"mime: " + info.Mime + "\n" +
		"via: " + info.Via + "\n" +
		"forwarded: " + info.Forwarded + "\n" +
		"country: " + info.Country + "\n" +
		"host: " + info.Host
}
