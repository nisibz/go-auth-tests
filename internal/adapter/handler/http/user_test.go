package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	handlerhttp "github.com/nisibz/go-auth-tests/internal/adapter/handler/http"
	"github.com/nisibz/go-auth-tests/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	args := m.Called(ctx, name, email, password)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context, limit, offset int64) ([]*domain.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error) {
	args := m.Called(ctx, id, name, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) CountUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetUserByID_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := handlerhttp.NewUserHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", handler.GetUserByID)

	expectedUser := &domain.User{ID: [12]byte{1, 2, 3}, Name: "Alice", Email: "alice@example.com"}
	mockService.On("GetUserByID", mock.Anything, "123").Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := handlerhttp.NewUserHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", handler.GetUserByID)

	mockService.On("GetUserByID", mock.Anything, "notfound").Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/users/notfound", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
