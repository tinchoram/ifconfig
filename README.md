# ifconfig

A modern IP information service inspired by ifconfig.me, built with Go and Fiber.

## 🚀 Features

- Get detailed information about your connection
- Clean and responsive web interface
- RESTful API endpoints
- Command-line interface support
- Docker support
- High performance using Fiber framework

## 🛠️ Technologies

- Go 1.21+
- [Fiber](https://github.com/gofiber/fiber) - Web framework
- HTML/CSS for frontend
- Docker for containerization

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
| `/all` | Returns all information in plain text |
| `/all.json` | Returns all information in JSON format |

## 🚦 Getting Started

### Prerequisites

- Go 1.21 or higher
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
go run main.go
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

Get User Agent:
```bash
curl localhost:3000/ua
```

## 🏗️ Project Structure

```
.
├── Dockerfile
├── README.md
├── main.go
├── go.mod
├── go.sum
├── views/
│   └── index.html
└── public/
    └── css/
        └── styles.css
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m '[Module] Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👤 Author

**@tinchoram**

* GitHub: [@tinchoram](https://github.com/tinchoram)

## ⭐️ Show your support

Give a ⭐️ if this project helped you!

## 📝 Notes

- Make sure to handle CORS and security considerations in production
- The service is designed to be lightweight and fast
- Contributions and suggestions are welcome

---
Made with ❤️ by @tinchoram
