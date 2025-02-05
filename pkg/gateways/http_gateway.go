package gateways

import "github.com/gofiber/fiber/v2"

type Gateway interface {
	GetHeader(string) string
	GetPort() string
	GetMethod() string
	GetHostname() string
	GetIP() string
	GetReqHeaders() map[string]string
}
type FiberHTTPGateway struct {
	ctx *fiber.Ctx
}

func NewHTTPGateway(c *fiber.Ctx) Gateway {
	return &FiberHTTPGateway{ctx: c}
}

func (g *FiberHTTPGateway) GetHeader(key string) string { return g.ctx.Get(key) }
func (g *FiberHTTPGateway) GetPort() string             { return g.ctx.Port() }
func (g *FiberHTTPGateway) GetMethod() string           { return g.ctx.Method() }
func (g *FiberHTTPGateway) GetHostname() string         { return g.ctx.Hostname() }
func (g *FiberHTTPGateway) GetIP() string               { return g.ctx.IP() }
func (g *FiberHTTPGateway) GetReqHeaders() map[string]string {
	headers := make(map[string]string)

	g.ctx.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	return headers
}
