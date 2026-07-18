package models

type IPInfo struct {
	IPAddr     string `json:"ip_addr"`
	RemoteHost string `json:"remote_host,omitempty"`
	UserAgent  string `json:"user_agent"`
	Language   string `json:"language"`
	Method     string `json:"method"`
	Encoding   string `json:"encoding"`
	Mime       string `json:"mime"`
	Via        string `json:"via"`
	Forwarded  string `json:"forwarded"`
	Connection string `json:"connection"`
	KeepAlive  string `json:"keep_alive"`
	Referer    string `json:"referer"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
	Timezone   string `json:"timezone,omitempty"`
	Continent  string `json:"continent,omitempty"`
	Host       string `json:"host"`
}
