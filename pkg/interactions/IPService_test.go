package interactions

import (
	"ifconfig/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock HTTPGateway
type MockHTTPGateway struct {
	mock.Mock
}

func (m *MockHTTPGateway) GetHeader(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockHTTPGateway) GetPort() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockHTTPGateway) GetMethod() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockHTTPGateway) GetHostname() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockHTTPGateway) GetIP() string {
	args := m.Called()
	return args.String(0)
}

// Test IPService
func TestIPService_GetIPInfo(t *testing.T) {
	mockGateway := new(MockHTTPGateway)
	ipService := NewIPService()

	// Set responses for all elements of the struct
	mockGateway.On("GetIP").Return("192.168.1.100")
	mockGateway.On("GetHeader", "User-Agent").Return("Mozilla/5.0")
	mockGateway.On("GetPort").Return("443")
	mockGateway.On("GetMethod").Return("GET")
	mockGateway.On("GetHostname").Return("example.com")
	mockGateway.On("GetHeader", "Accept-Language").Return("es-ES,es;q=0.9")
	mockGateway.On("GetHeader", "Accept-Encoding").Return("gzip, deflate, br")
	mockGateway.On("GetHeader", "Accept").Return("text/html,application/xhtml+xml")
	mockGateway.On("GetHeader", "Via").Return("1.1 proxy")
	mockGateway.On("GetHeader", "Connection").Return("keep-alive")
	mockGateway.On("GetHeader", "Keep-Alive").Return("timeout=5, max=100")
	mockGateway.On("GetHeader", "Referer").Return("https://example.com")
	mockGateway.On("GetHeader", "Cf-Ipcountry").Return("AR")
	mockGateway.On("GetHeader", "Cf-Connecting-Ip").Return("")
	mockGateway.On("GetHeader", "X-Real-Ip").Return("")
	mockGateway.On("GetHeader", "X-Forwarded-For").Return("192.168.1.1")

	// Call the function with the mock
	ipInfo := ipService.GetIPInfo(mockGateway)

	expectedIPInfo := models.IPInfo{
		IPAddr:     "192.168.1.1",
		RemoteHost: "unavailable",
		UserAgent:  "Mozilla/5.0",
		Port:       "443",
		Language:   "es-ES,es;q=0.9",
		Method:     "GET",
		Encoding:   "gzip, deflate, br",
		Mime:       "text/html,application/xhtml+xml",
		Via:        "1.1 proxy",
		Forwarded:  "192.168.1.1",
		Connection: "keep-alive",
		KeepAlive:  "timeout=5, max=100",
		Referer:    "https://example.com",
		Country:    "AR",
		Host:       "example.com",
	}

	assert.Equal(t, expectedIPInfo, ipInfo)
}

// Test for GetRealIP with different scenarios
func TestIPService_GetRealIP(t *testing.T) {
	mockGateway := new(MockHTTPGateway)
	ipService := NewIPService()

	// Case 1: Cf-Connecting-Ip has a value
	mockGateway.On("GetHeader", "Cf-Connecting-Ip").Return("203.0.113.1")
	mockGateway.On("GetHeader", "X-Real-Ip").Return("")
	mockGateway.On("GetHeader", "X-Forwarded-For").Return("")
	mockGateway.On("GetIP").Return("192.168.1.100")

	realIP := ipService.GetRealIP(mockGateway)
	assert.Equal(t, "203.0.113.1", realIP)

	// Case 2: X-Real-Ip has a value
	mockGateway = new(MockHTTPGateway) // Reset mock
	mockGateway.On("GetHeader", "Cf-Connecting-Ip").Return("")
	mockGateway.On("GetHeader", "X-Real-Ip").Return("198.51.100.2")
	mockGateway.On("GetHeader", "X-Forwarded-For").Return("")
	mockGateway.On("GetIP").Return("192.168.1.100")

	realIP = ipService.GetRealIP(mockGateway)
	assert.Equal(t, "198.51.100.2", realIP)

	// Case 3: X-Forwarded-For has multiple IPs
	mockGateway = new(MockHTTPGateway) // Reset mock
	mockGateway.On("GetHeader", "Cf-Connecting-Ip").Return("")
	mockGateway.On("GetHeader", "X-Real-Ip").Return("")
	mockGateway.On("GetHeader", "X-Forwarded-For").Return("203.0.113.3, 198.51.100.4, 192.168.1.200")
	mockGateway.On("GetIP").Return("192.168.1.100")

	realIP = ipService.GetRealIP(mockGateway)
	assert.Equal(t, "203.0.113.3", realIP)

	// Case 4: None of the headers are present, use GetIP()
	mockGateway = new(MockHTTPGateway) // Reset mock
	mockGateway.On("GetHeader", "Cf-Connecting-Ip").Return("")
	mockGateway.On("GetHeader", "X-Real-Ip").Return("")
	mockGateway.On("GetHeader", "X-Forwarded-For").Return("")
	mockGateway.On("GetIP").Return("192.168.1.100")

	realIP = ipService.GetRealIP(mockGateway)
	assert.Equal(t, "192.168.1.100", realIP)
}
