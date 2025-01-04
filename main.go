package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"strings"
	"time"
)

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
}

// getRealIP from Cloudflare
func getRealIP(c *fiber.Ctx) string {
	// Get Cloudflare IP
	if cfIP := c.Get("Cf-Connecting-Ip"); cfIP != "" {
		return cfIP
	}

	// Get X-Real-IP
	if realIP := c.Get("X-Real-Ip"); realIP != "" {
		return realIP
	}

	// Get X-Forwarded-For
	if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	return c.IP()
}

func main() {
	// Engine templates
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Static
	app.Static("/", "./public")

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		userAgent := c.Get("User-Agent")
		isCurl := strings.HasPrefix(strings.ToLower(userAgent), "curl")
		if isCurl {
			return c.SendString(getRealIP(c))
		}

		info := getIPInfo(c)
		return c.Render("index", fiber.Map{
			"info": info,
			"year": time.Now().Year(),
		})
	})

	// Routes API
	app.Get("/ip", func(c *fiber.Ctx) error {
		return c.SendString(getRealIP(c))
	})

	app.Get("/ua", func(c *fiber.Ctx) error {
		return c.SendString(c.Get("User-Agent"))
	})

	app.Get("/lang", func(c *fiber.Ctx) error {
		return c.SendString(c.Get("Accept-Language"))
	})

	app.Get("/encoding", func(c *fiber.Ctx) error {
		return c.SendString(c.Get("Accept-Encoding"))
	})

	app.Get("/mime", func(c *fiber.Ctx) error {
		return c.SendString(c.Get("Accept"))
	})

	app.Get("/charset", func(c *fiber.Ctx) error {
		return c.SendString(c.Get("Accept-Charset"))
	})

	app.Get("/forwarded", func(c *fiber.Ctx) error {
		return c.SendString(c.Get("X-Forwarded-For"))
	})

	app.Get("/all", func(c *fiber.Ctx) error {
		info := getIPInfo(c)
		response := formatPlainTextResponse(info)
		return c.SendString(response)
	})

	app.Get("/all.json", func(c *fiber.Ctx) error {
		info := getIPInfo(c)
		return c.JSON(info)
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		method := string(c.Request().Header.Method())
		uri := string(c.Request().RequestURI())
		queryParams := c.Request().URI().QueryArgs().String()

		requestInfo := map[string]interface{}{
			"headers":      headers,
			"method":       method,
			"uri":          uri,
			"query_params": queryParams,
			"body":         string(c.Body()),
			"real_ip":      getRealIP(c),
			"country":      c.Get("Cf-Ipcountry"),
		}

		return c.JSON(requestInfo)
	})

	app.Listen(":3000")
}

func getIPInfo(c *fiber.Ctx) IPInfo {
	return IPInfo{
		IPAddr:     getRealIP(c),
		RemoteHost: "unavailable",
		UserAgent:  c.Get("User-Agent"),
		Port:       c.Port(),
		Language:   c.Get("Accept-Language"),
		Method:     c.Method(),
		Encoding:   c.Get("Accept-Encoding"),
		Mime:       c.Get("Accept"),
		Via:        c.Get("Via"),
		Forwarded:  c.Get("X-Forwarded-For"),
		Connection: c.Get("Connection"),
		KeepAlive:  c.Get("Keep-Alive"),
		Referer:    c.Get("Referer"),
		Country:    c.Get("Cf-Ipcountry"),
	}
}

func formatPlainTextResponse(info IPInfo) string {
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
		"country: " + info.Country
}
