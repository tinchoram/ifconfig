package utils

import "ifconfig/pkg/models"

// FormatPlainTextResponse renders the IP info as plain "key: value" lines for /all.
func FormatPlainTextResponse(info models.IPInfo) string {
	return "ip_addr: " + info.IPAddr + "\n" +
		"remote_host: " + info.RemoteHost + "\n" +
		"user_agent: " + info.UserAgent + "\n" +
		"language: " + info.Language + "\n" +
		"method: " + info.Method + "\n" +
		"encoding: " + info.Encoding + "\n" +
		"mime: " + info.Mime + "\n" +
		"via: " + info.Via + "\n" +
		"forwarded: " + info.Forwarded + "\n" +
		"city: " + info.City + "\n" +
		"region: " + info.Region + "\n" +
		"country: " + info.Country + "\n" +
		"postal_code: " + info.PostalCode + "\n" +
		"latitude: " + info.Latitude + "\n" +
		"longitude: " + info.Longitude + "\n" +
		"timezone: " + info.Timezone + "\n" +
		"continent: " + info.Continent + "\n" +
		"host: " + info.Host
}
