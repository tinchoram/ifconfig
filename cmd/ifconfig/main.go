package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

// defaultTrustedProxies trusts loopback only: proxy headers are ignored unless
// the operator explicitly lists their reverse proxy in TRUSTED_PROXIES.
const defaultTrustedProxies = "127.0.0.1,::1"

const shutdownTimeout = 10 * time.Second

func main() {
	app, err := newApp(appConfig{
		TrustedProxies: splitAndTrim(envOr("TRUSTED_PROXIES", defaultTrustedProxies)),
		ProxyHeader:    envOr("PROXY_HEADER", fiber.HeaderXForwardedFor),
		ViewsDir:       "./views",
		PublicDir:      "./public",
	})
	if err != nil {
		log.Fatal(err)
	}

	addr, err := listenAddr(envOr("HOST", ""), envOr("PORT", "3000"))
	if err != nil {
		log.Fatal(err)
	}

	// ready closes once the server is actually accepting connections. Gating
	// shutdown on it prevents a signal that arrives during startup from firing
	// an ineffective shutdown that the server would then outlive.
	ready := make(chan struct{})
	app.Hooks().OnListen(func(fiber.ListenData) error {
		close(ready)
		return nil
	})

	// Register the signal handler before Listen starts so a signal delivered
	// during startup is caught, not dropped to the default disposition.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() { serverErr <- app.Listen(addr) }()

	select {
	case err := <-serverErr:
		// Listen failed (startup or runtime) before any signal arrived.
		if err != nil {
			log.Fatal(err)
		}
		return
	case <-sig:
	}

	// Do not shut down until the server is confirmed up, or has already died.
	select {
	case <-ready:
	case err := <-serverErr:
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Println("shutting down")
	if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
		// Drain deadline hit with connections still in flight: report a
		// degraded shutdown with a non-zero exit so supervisors can tell it
		// apart from a clean stop.
		log.Printf("graceful shutdown incomplete: %v", err)
		os.Exit(1)
	}
	if err := <-serverErr; err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}

// listenAddr builds the host:port listen address, rejecting non-numeric or
// out-of-range ports so misconfiguration fails at startup.
func listenAddr(host, port string) (string, error) {
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return "", fmt.Errorf("invalid PORT %q: must be a number between 1 and 65535", port)
	}
	return net.JoinHostPort(host, port), nil
}

func envOr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func splitAndTrim(list string) []string {
	var out []string
	for _, item := range strings.Split(list, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
