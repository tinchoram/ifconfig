# ifconfig

A modern IP information service inspired by ifconfig.me, built with Go and Fiber.

## 🚀 Features

- Get detailed information about your connection
- Clean and responsive web interface
- RESTful API endpoints
- Command-line interface support
- Docker support
- High performance using Fiber framework
- Uses **GIMU Architecture** for better scalability and maintainability

## 🛠️ Technologies

- **Go 1.26+**
- [Fiber](https://github.com/gofiber/fiber) - Web framework
- **GIMU Architecture** (Designed by [@tinchoram](https://github.com/tinchoram))
- HTML/CSS for frontend
- Docker for containerization

## 📚 GIMU Architecture

This project follows the **GIMU** architecture, a structure designed by **@tinchoram** to simplify and scale applications efficiently. The architecture consists of the following layers:

### **📌 GIMU Layers**
| Layer | Description |
|--------|------------|
| **Gateways** | Handles external interactions (e.g., HTTP requests, APIs, databases) |
| **Interactions** | Implements business logic and use cases |
| **Models** | Defines the domain entities and data structures |
| **Utils** | Provides utility functions for various tasks |

### **📂 Project Structure**
```
.
├── Dockerfile
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── ifconfig/
│       └── main.go   # Entry point of the application
├── pkg/
│   ├── gateways/ # Adapters for external interactions (HTTP)
│   │   ├── http_gateway.go
│   ├── interactions/ # Business logic and ports (IP processing)
│   │   ├── ip_service.go
│   │   ├── ip_service_test.go
│   ├── models/  # Defines the data structures
│   │   ├── ip_info.go
│   ├── utils/  # Utility functions (formatting)
│   │   ├── formatter.go
├── views/
│   └── index.html
└── public/
    └── css/
        └── styles.css
```

📌 **Advantages of GIMU:**
- **Clear separation of concerns** (each layer has a single responsibility).
- **Easier testing** (mocks can be injected into interactions).
- **Scalability** (you can add new gateways, interactions, and models without modifying the core logic).
- **Better maintainability** (logical separation makes debugging and updates easier).

## 📋 API Endpoints

| Endpoint | Description |
|----------|-------------|
| `/` | Returns IP address for curl, web interface for browsers |
| `/ip` | Returns only the IP address |
| `/ua` | Returns User Agent |
| `/lang` | Returns Accept-Language |
| `/encoding` | Returns Accept-Encoding |
| `/mime` | Returns accepted MIME types |
| `/charset` | Returns Accept-Charset |
| `/forwarded` | Returns X-Forwarded-For |
| `/headers` | Returns all request headers |
| `/all` | Returns all information in plain text |
| `/all.json` | Returns all information in JSON format |
| `/ping` | Returns request details (headers, method, IP, etc.) |
| `/details.json` | Returns extended information with timestamp |
| `/status` | Health check: returns `{"status":"OK"}` with a timestamp |

`/all` and `/all.json` include the client's **reverse DNS** (PTR) hostname and, when the service runs behind Cloudflare with the **Add visitor location headers** managed transform enabled, **geolocation** fields (city, region, country, postal code, coordinates, timezone, continent). Location fields are omitted/empty when that transform is off.

## ⚙️ Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `TRUSTED_PROXIES` | `127.0.0.1,::1` | Comma-separated IPs/CIDR ranges of reverse proxies allowed to set the proxy header |
| `PROXY_HEADER` | `X-Forwarded-For` | Header used to resolve the client IP. Behind Cloudflare, set it to `Cf-Connecting-Ip` |
| `PORT` | `3000` | Port the server listens on |
| `HOST` | *(all interfaces)* | Address the server binds to |

The server handles `SIGINT`/`SIGTERM` with a graceful shutdown (10s timeout), so in-flight requests complete when Docker stops the container.

> **Security note:** by default only loopback is trusted, so forged proxy headers from external clients are ignored and the TCP socket IP is reported. When deploying behind a reverse proxy, set `TRUSTED_PROXIES` to that proxy's IP or CIDR range (see `examples/docker-compose.yml`). Do **not** widen the trust list for standalone `docker run -p` deployments — trusting broad ranges lets any peer on those networks spoof its reported address.

## 🌍 How the data is resolved

The service surfaces three kinds of data, each obtained a different way.

### Client IP

The IP shown is the real client address, resolved through the proxy chain. When the
peer connecting to the app is a trusted proxy (`TRUSTED_PROXIES`), the app reads the
client IP from `PROXY_HEADER` instead of the raw socket address. Behind Cloudflare that
header is `Cf-Connecting-Ip`, which Cloudflare sets to the original visitor's IP.
Untrusted peers cannot spoof it — forged headers are ignored (see the security note above).

### Reverse DNS (PTR)

The app performs a reverse DNS lookup on the resolved client IP — a PTR query that asks
"which hostname points to this IP?". It is bounded by a **500 ms timeout** so a slow or
missing record never delays the page, and loopback/unspecified addresses are skipped
(they have no useful PTR). If nothing resolves, the field is left empty.

### Geolocation (provided by Cloudflare)

The app does **not** geolocate anything itself. It reads location headers that Cloudflare
injects at its edge: `Cf-Ipcity`, `Cf-Region`, `Cf-Ipcountry`, `Cf-Postal-Code`,
`Cf-Iplatitude`, `Cf-Iplongitude`, `Cf-Timezone`, `Cf-Ipcontinent`. These require enabling
the **Add visitor location headers** managed transform in Cloudflare
(*Rules → Managed Transforms*); `Cf-Ipcountry` is sent whenever IP geolocation is on. With
the transform off, the fields are simply empty.

**Where does Cloudflare's geo data come from?** It is an IP-to-location *mapping*, not a
measurement of your device:

1. **Regional Internet Registries** (LACNIC, ARIN, RIPE, APNIC, AFRINIC) record which
   organization owns each IP block and in which region — this makes *country* highly reliable.
2. **ISPs** announce their ranges via BGP and may publish per-range location hints
   (geofeeds, RFC 8805).
3. **Commercial geo databases** (e.g. MaxMind) plus **Cloudflare's own network** (300+ edge
   cities, latency data) refine that to *city*-level estimates.

Accuracy therefore varies: country is solid, while city/coordinates are estimates
(coordinates are typically the centroid of the city or postal area) and can be off for
VPN, mobile, or corporate IPs.

## 🚦 Getting Started

### Prerequisites

- **Go 1.26+**
- Docker (optional)

### Local Setup

1. Clone the repository:
```bash
git clone https://github.com/tinchoram/ifconfig.git
cd ifconfig
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run ./cmd/ifconfig
```

The service will be available at `http://localhost:3000`

### Docker Setup

1. Build the Docker image:
```bash
docker build -t ifconfig .
```

2. Run the container:
```bash
docker run -p 3000:3000 ifconfig
```

## 📝 Usage Examples

Get your IP address:
```bash
curl localhost:3000
```

Get all information in JSON format:
```bash
curl localhost:3000/all.json
```

Get all headers:
```bash
curl localhost:3000/headers
```

Get extended details:
```bash
curl localhost:3000/details.json
```

## 🧗‍♂️ Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m '[Module] Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🐜 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👤 Author

**[@tinchoram](https://github.com/tinchoram)**

## ⭐️ Show your support

Give a ⭐️ if this project helped you!

## 📝 Notes

- Make sure to handle CORS and security considerations in production.
- The service is designed to be lightweight and fast.
- Contributions and suggestions are welcome.

---
Made with ❤️ by [@tinchoram](https://github.com/tinchoram)

