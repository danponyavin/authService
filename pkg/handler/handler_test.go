package handler

import (
	"AuthService/pkg/models"
	"AuthService/pkg/service"
	mock_service "AuthService/pkg/service/mocks"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetTokens(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, auth models.AuthModel)

	tests := []struct {
		name                 string
		input                models.AuthModel
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: models.AuthModel{
				UserID:   uuid.New(),
				ClientIP: "192.0.2.1",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, auth models.AuthModel) {
				s.EXPECT().GetTokens(auth).Return(models.Tokens{AccessToken: "token1", RefreshToken: "token2"}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"access_token":"token1","refresh_token":"token2"}`,
		},
		{
			name: "Service error",
			input: models.AuthModel{
				UserID:   uuid.New(),
				ClientIP: "192.0.2.1",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, auth models.AuthModel) {
				s.EXPECT().GetTokens(auth).Return(models.Tokens{}, errors.New("New service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"Service error"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mock_service.NewMockAuthorization(c)
			testCase.mockBehaviour(authService, testCase.input)

			services := &service.Service{Authorization: authService}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/tokens/:user_id", handler.GetTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/tokens/%s", testCase.input.UserID.String()), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_RefreshTokens(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, inputData models.RefreshTokenRequest)

	tests := []struct {
		name                 string
		inputBody            string
		inputData            models.RefreshTokenRequest
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"refresh_token": "refreshToken"}`,
			inputData: models.RefreshTokenRequest{
				RefreshToken: "refreshToken",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, inputData models.RefreshTokenRequest) {
				s.EXPECT().RefreshTokens(models.RefreshModel{
					RefreshToken: inputData.RefreshToken,
					ClientIP:     "192.0.2.1",
				}).Return(models.Tokens{AccessToken: "token1", RefreshToken: "token2"}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"access_token":"token1","refresh_token":"token2"}`,
		},
		{
			name:      "Invalid Refresh Token",
			inputBody: `{"refresh_token": "token1"}`,
			inputData: models.RefreshTokenRequest{
				RefreshToken: "token1",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, inputData models.RefreshTokenRequest) {
				s.EXPECT().RefreshTokens(models.RefreshModel{
					RefreshToken: inputData.RefreshToken,
					ClientIP:     "192.0.2.1",
				}).Return(models.Tokens{}, service.InvalidRefreshTokenError)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Invalid Refresh Token"}`,
		},
		{
			name:      "Token Expired",
			inputBody: `{"refresh_token": "refreshToken"}`,
			inputData: models.RefreshTokenRequest{
				RefreshToken: "refreshToken",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, inputData models.RefreshTokenRequest) {
				s.EXPECT().RefreshTokens(models.RefreshModel{
					RefreshToken: inputData.RefreshToken,
					ClientIP:     "192.0.2.1",
				}).Return(models.Tokens{}, service.TokenExpiredError)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Token is expired"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"refresh_token": "refreshToken"}`,
			inputData: models.RefreshTokenRequest{
				RefreshToken: "refreshToken",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, inputData models.RefreshTokenRequest) {
				s.EXPECT().RefreshTokens(models.RefreshModel{
					RefreshToken: inputData.RefreshToken,
					ClientIP:     "192.0.2.1",
				}).Return(models.Tokens{}, errors.New("New service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"Service error"}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mock_service.NewMockAuthorization(c)
			testCase.mockBehaviour(authService, testCase.inputData)

			services := &service.Service{Authorization: authService}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/auth/refresh", handler.RefreshTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
