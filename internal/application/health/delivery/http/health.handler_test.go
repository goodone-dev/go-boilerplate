package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/domain/health"
	healthmock "github.com/goodone-dev/go-boilerplate/internal/domain/health/mocks"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	logger.Disabled()
	code := m.Run()

	os.Exit(code)
}

func TestNewHealthHandler(t *testing.T) {
	mockService := healthmock.NewHealthCheckerMock(t)
	handler := NewHealthHandler(mockService)

	assert.NotNil(t, handler)
}

func TestHealthHandler_LiveCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health/live", nil)

	handler.LiveCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"up"`)
}

func TestHealthHandler_ReadyCheck_AllServicesUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService1 := healthmock.NewHealthCheckerMock(t)
	mockService2 := healthmock.NewHealthCheckerMock(t)

	mockService1.On("Ping", mock.Anything).Return(nil)
	mockService2.On("Ping", mock.Anything).Return(nil)

	handler := NewHealthHandler(mockService1, mockService2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health/ready", nil)

	handler.ReadyCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"up"`)

	mockService1.AssertExpectations(t)
	mockService2.AssertExpectations(t)
}

func TestHealthHandler_ReadyCheck_ServiceDown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService1 := healthmock.NewHealthCheckerMock(t)
	mockService2 := healthmock.NewHealthCheckerMock(t)

	mockService1.On("Ping", mock.Anything).Return(nil)
	mockService2.On("Ping", mock.Anything).Return(errors.New("connection failed"))

	handler := NewHealthHandler(mockService1, mockService2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health/ready", nil)

	handler.ReadyCheck(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"down"`)

	mockService1.AssertExpectations(t)
	mockService2.AssertExpectations(t)
}

func TestHealthHandler_ReadyCheck_NoServices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health/ready", nil)

	handler.ReadyCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestParsePackageName(t *testing.T) {
	tests := []struct {
		name     string
		service  health.HealthChecker
		expected string
	}{
		{
			name:     "MockHealthService",
			service:  healthmock.NewHealthCheckerMock(t),
			expected: "http",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePackageName(tt.service)
			assert.NotEmpty(t, result)
		})
	}
}

func TestParsePackageName_EmptyResult(t *testing.T) {
	mockService := healthmock.NewHealthCheckerMock(t)
	result := parsePackageName(mockService)
	assert.NotEmpty(t, result)
}
