package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	"ifconfig/pkg/gateways"
	"ifconfig/pkg/interactions"
	"ifconfig/pkg/utils"
)

// appConfig holds the runtime settings injected into newApp.
type appConfig struct {
	// TrustedProxies lists the peer IPs/CIDR ranges allowed to set ProxyHeader.
	// Requests from any other peer report the socket IP, ignoring the header.
	TrustedProxies []string
	ProxyHeader    string
	ViewsDir       string
	PublicDir      string
}

func newApp(cfg appConfig) (*fiber.App, error) {
	engine := html.New(cfg.ViewsDir, ".html")
	// Load templates eagerly: Fiber only warn-logs a lazy load failure, which
	// would leave a healthy-looking server that 500s on every render.
	if err := engine.Load(); err != nil {
		return nil, fmt.Errorf("loading templates from %s: %w", cfg.ViewsDir, err)
	}

	app := fiber.New(fiber.Config{
		Views: engine,
		// EnableIPValidation makes c.IP() return the first valid IP of a
		// forwarded chain instead of the raw header value.
		EnableTrustedProxyCheck: true,
		TrustedProxies:          cfg.TrustedProxies,
		ProxyHeader:             cfg.ProxyHeader,
		EnableIPValidation:      true,
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            10 * time.Second,
		IdleTimeout:             60 * time.Second,
		// 1 MB is plenty for an echo service; /ping reflects the request body.
		BodyLimit: 1 * 1024 * 1024,
	})

	// Fiber does not recover from handler panics by default.
	app.Use(recover.New())

	ipService := interactions.NewService()

	if cfg.PublicDir != "" {
		app.Static("/", cfg.PublicDir)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)

		userAgent := gateway.GetHeader("User-Agent")
		isCurl := strings.HasPrefix(strings.ToLower(userAgent), "curl")
		if isCurl {
			return c.SendString(gateway.GetIP())
		}

		info := ipService.GetIPInfo(gateway)
		return c.Render("index", fiber.Map{
			"info": info,
			"year": time.Now().Year(),
		})
	})

	// Routes API
	app.Get("/ip", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.SendString(gateway.GetIP())
	})

	app.Get("/ua", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.SendString(gateway.GetHeader("User-Agent"))
	})

	app.Get("/lang", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.SendString(gateway.GetHeader("Accept-Language"))
	})

	app.Get("/encoding", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.SendString(gateway.GetHeader("Accept-Encoding"))
	})

	app.Get("/mime", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.SendString(gateway.GetHeader("Accept"))
	})

	app.Get("/charset", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.SendString(gateway.GetHeader("Accept-Charset"))
	})

	app.Get("/forwarded", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.SendString(gateway.GetHeader("X-Forwarded-For"))
	})

	app.Get("/all", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		info := ipService.GetIPInfo(gateway)
		response := utils.FormatPlainTextResponse(info)
		return c.SendString(response)
	})

	app.Get("/all.json", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		info := ipService.GetIPInfo(gateway)
		return c.JSON(info)
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)

		headers := c.GetReqHeaders()
		method := gateway.GetMethod()
		uri := c.OriginalURL()
		queryParams := c.Queries()

		requestInfo := map[string]interface{}{
			"headers":      headers,
			"method":       method,
			"uri":          uri,
			"query_params": queryParams,
			"body":         string(c.Body()),
			"real_ip":      gateway.GetIP(),
			"country":      gateway.GetHeader("Cf-Ipcountry"),
		}

		return c.JSON(requestInfo)
	})

	app.Get("/details.json", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		info := ipService.GetIPInfo(gateway)

		detailedInfo := map[string]interface{}{
			"ip_info":    info,
			"user_agent": gateway.GetHeader("User-Agent"),
			"timestamp":  time.Now().Format(time.RFC3339),
		}

		return c.JSON(detailedInfo)
	})

	app.Get("/headers", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)
		return c.JSON(gateway.GetReqHeaders())
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK", "timestamp": time.Now().Format(time.RFC3339)})
	})

	return app, nil
}
