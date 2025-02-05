package main

import (
	"github.com/gofiber/fiber/v2/log"
	"ifconfig/pkg/gateways"
	"ifconfig/pkg/interactions"
	"ifconfig/pkg/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Engine templates
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Start service
	ipService := interactions.NewIPService()

	// Static files
	app.Static("/", "./public")

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		gateway := gateways.NewHTTPGateway(c)

		userAgent := gateway.GetHeader("User-Agent")
		isCurl := strings.HasPrefix(strings.ToLower(userAgent), "curl")
		if isCurl {
			return c.SendString(ipService.GetRealIP(gateway))
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
		return c.SendString(ipService.GetRealIP(gateway))
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
			"real_ip":      ipService.GetRealIP(gateway),
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

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
		return
	}
}
