package models

type IPInfo struct {
	IPAddr     string `json:"ip_addr"`
	RemoteHost string `json:"remote_host"`
	UserAgent  string `json:"user_agent"`
	Port       string `json:"port"`
	Language   string `json:"language"`
	Method     string `json:"method"`
	Encoding   string `json:"encoding"`
	Mime       string `json:"mime"`
	Via        string `json:"via"`
	Forwarded  string `json:"forwarded"`
	Connection string `json:"connection"`
	KeepAlive  string `json:"keep_alive"`
	Referer    string `json:"referer"`
	Country    string `json:"country,omitempty"`
	Host       string `json:"host"`
}