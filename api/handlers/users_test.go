package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/id-tarzanych/lets-go-chat/mocks"
	"github.com/id-tarzanych/lets-go-chat/models"
)

func TestUsers_HandleUserCreate(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getUserHandlerMocks(t)

	handlers := Users{
		logger:    loggerMock,
		userRepo:  userRepoMock,
		tokenRepo: tokenRepoMock,
	}

	tests := []struct {
		name        string
		requestJSON string
		wantCode    int
		wantMessage string
	}{
		{
			name:        "Invalid syntax",
			requestJSON: "{123]",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Syntax error",
		},
		{
			name:        "Empty username",
			requestJSON: "{\"password\": \"12345678\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty username",
			requestJSON: "{\"userName\": \"\", \"password\": \"12345678\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty password",
			requestJSON: "{\"userName\": \"testuser\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty password",
			requestJSON: "{\"userName\": \"username\", \"password\": \"\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty username and password",
			requestJSON: "{\"userName\": \"\", \"password\": \"\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty username and password",
			requestJSON: "{}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Short password",
			requestJSON: "{\"userName\": \"username\", \"password\": \"123\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Password should be at least 8 characters",
		},
		{
			name:        "Username conflict",
			requestJSON: "{\"userName\": \"existingUser\", \"password\": \"12345678\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "User with username existingUser already exists",
		},
		{
			name:        "Storage operation error",
			requestJSON: "{\"userName\": \"storageErrorUser\", \"password\": \"12345678\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Could not create user storageErrorUser",
		},
		{
			name:        "Successful user creation",
			requestJSON: "{\"userName\": \"newUser\", \"password\": \"12345678\"}",
			wantCode:    http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			handlers.HandleUserCreate().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.requestJSON)))
			response := w.Result()

			if response.StatusCode != tt.wantCode {
				t.Errorf("Incorrect status code, wanted %d, got %d.", tt.wantCode, response.StatusCode)
			}

			responseBody := strings.TrimSpace(w.Body.String())
			if tt.wantCode != http.StatusOK && tt.wantMessage != responseBody {
				t.Errorf("Incorrect error message, wanted \"%s\", got \"%s\"", tt.wantMessage, responseBody)
			}
		})
	}
}

func TestUsers_HandleUserLogin(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getUserHandlerMocks(t)

	handlers := Users{
		logger:    loggerMock,
		userRepo:  userRepoMock,
		tokenRepo: tokenRepoMock,
	}

	tests := []struct {
		name        string
		requestJSON string
		wantCode    int
		wantMessage string
	}{
		{
			name:        "Invalid syntax",
			requestJSON: "{123]",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Syntax error",
		},
		{
			name:        "Empty username",
			requestJSON: "{\"password\": \"12345678\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty username",
			requestJSON: "{\"userName\": \"\", \"password\": \"12345678\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty password",
			requestJSON: "{\"userName\": \"testuser\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty password",
			requestJSON: "{\"userName\": \"username\", \"password\": \"\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty username and password",
			requestJSON: "{\"userName\": \"\", \"password\": \"\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Empty username and password",
			requestJSON: "{}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Empty username or password",
		},
		{
			name:        "Non-existing user",
			requestJSON: "{\"userName\": \"newUser\", \"password\": \"12345678\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "User newUser does not exist",
		},
		{
			name:        "Valid User",
			requestJSON: "{\"userName\": \"existingUser\", \"password\": \"12345678\"}",
			wantCode:    http.StatusOK,
		},
		{
			name:        "Incorrect password",
			requestJSON: "{\"userName\": \"existingUser\", \"password\": \"1234567890\"}",
			wantCode:    http.StatusBadRequest,
			wantMessage: "Invalid username/password",
		},
		{
			name:        "Token Storage Error",
			requestJSON: "{\"userName\": \"tokenStorageError\", \"password\": \"12345678\"}",
			wantCode:    http.StatusInternalServerError,
			wantMessage: "Could not generate one-time token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			handlers.HandleUserLogin().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.requestJSON)))
			response := w.Result()

			if response.StatusCode != tt.wantCode {
				t.Errorf("Incorrect status code, wanted %d, got %d.", tt.wantCode, response.StatusCode)
			}

			responseBody := strings.TrimSpace(w.Body.String())
			if tt.wantCode != http.StatusOK && tt.wantMessage != responseBody {
				t.Errorf("Incorrect error message, wanted \"%s\", got \"%s\"", tt.wantMessage, responseBody)
			}
		})
	}

	loggerMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
	tokenRepoMock.AssertExpectations(t)
}

func getUserHandlerMocks(t *testing.T) (*mocks.FieldLogger, *mocks.UserRepository, *mocks.TokenRepository) {
	loggerMock := &mocks.FieldLogger{}
	userRepoMock := &mocks.UserRepository{}
	tokenRepoMock := &mocks.TokenRepository{}

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"existingUser",
	).Maybe().Return(models.User{ID: "uuid", UserName: "existingUser", PasswordHash: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f"}, nil)

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"tokenStorageError",
	).Maybe().Return(models.User{ID: "uuid-token-storage-error", UserName: "tokenStorageError", PasswordHash: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f"}, nil)

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"storageErrorUser",
	).Maybe().Return(models.User{ID: "uuid"}, errors.New("storage error"))

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"newUser",
	).Maybe().Return(models.User{}, errors.New("could not find user"))

	userRepoMock.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(u *models.User) bool { return u.UserName == "storageErrorUser" }),
	).Maybe().Return(errors.New("storage error"))

	userRepoMock.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(u *models.User) bool { return u.UserName != "storageErrorUser" }),
	).Maybe().Return(nil)

	tokenRepoMock.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(t *models.Token) bool { return t.UserId == "uuid-token-storage-error" }),
	).Maybe().Return(errors.New("storage error"))

	tokenRepoMock.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Maybe().Return(nil)

	return loggerMock, userRepoMock, tokenRepoMock
}
