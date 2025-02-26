package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"area/api"
	"area/schemas"
	"area/test"
)

type MockGithubController struct {
	mock.Mock
}

func (m *MockGithubController) RedirectToService(ctx *gin.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockGithubController) HandleServiceCallback(ctx *gin.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockGithubController) HandleServiceCallbackMobile(ctx *gin.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockGithubController) GetUserInfo(ctx *gin.Context) (schemas.UserCredentials, error) {
	args := m.Called(ctx)
	return args.Get(0).(schemas.UserCredentials), args.Error(1)
}

func TestGithubAPI(t *testing.T) {
	t.Parallel()

	mockController := new(MockGithubController)
	router := gin.Default()
	apiRoutes := router.Group("/api")
	mockUserService := new(test.MockUserService)
	api.NewGithubAPI(mockController, apiRoutes, mockUserService)

	t.Run("TestRedirectToService", func(t *testing.T) {
		t.Parallel()

		mockController.On("RedirectToService", mock.Anything).Return("http://example.com/auth", nil)

		responseRecorder := httptest.NewRecorder()
		ctx := context.Background()

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/api/github/auth", nil)
		router.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Contains(t, responseRecorder.Body.String(), "http://example.com/auth")
	})
	t.Run("TestHandleServiceCallback", func(t *testing.T) {
		t.Parallel()

		mockController.On("HandleServiceCallback", mock.Anything).Return("mock_token", nil)

		responseRecorder := httptest.NewRecorder()
		ctx := context.Background()

		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/api/github/auth/callback", nil)
		router.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Contains(t, responseRecorder.Body.String(), "mock_token")
	})
	t.Run("TestHandleServiceCallback", func(t *testing.T) {
		t.Parallel()

		mockController.On("HandleServiceCallbackMobile", mock.Anything).
			Return("mock_mobile_token", nil)

		responseRecorder := httptest.NewRecorder()
		ctx := context.Background()

		req, _ := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			"/api/github/auth/callback/mobile",
			nil,
		)
		router.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Contains(t, responseRecorder.Body.String(), "mock_mobile_token")
	})
	// t.Run("TestHandleServiceCallback", func(t *testing.T) {
	// 	mockUserInfo := schemas.UserCredentials{Username: "Test User", Email: "aze"}
	// 	mockController.On("GetUserInfo", mock.Anything).Return(mockUserInfo, nil)

	// 	responseRecorder := httptest.NewRecorder()
	// 	req, _ := http.NewRequest(http.MethodGet, "/api/github/info", nil)
	// 	router.ServeHTTP(responseRecorder, req)

	// 	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	// 	assert.Contains(t, responseRecorder.Body.String(), "Test User")
	// })
}
