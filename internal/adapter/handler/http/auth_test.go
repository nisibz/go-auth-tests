package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	handlerhttp "github.com/nisibz/go-auth-tests/internal/adapter/handler/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(name, email, password string) (string, error) {
	args := m.Called(name, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handlerhttp.NewAuthHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handler.Register)

	mockService.On("Register", "John", "john@example.com", "password123").Return("mocked_token", nil)

	body := `{"name": "John", "email": "john@example.com", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestRegister_ValidationError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handlerhttp.NewAuthHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handler.Register)

	body := `{"name": "", "email": "bad-email", "password": "123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestLogin_Unauthorized(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handlerhttp.NewAuthHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", handler.Login)

	mockService.On("Login", "invalid@example.com", "wrongpass").Return("", errors.New("invalid credentials"))

	body := `{"email": "invalid@example.com", "password": "wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
