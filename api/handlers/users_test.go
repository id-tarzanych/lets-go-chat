package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/id-tarzanych/lets-go-chat/mocks"
	"github.com/id-tarzanych/lets-go-chat/models"
	"github.com/stretchr/testify/mock"
)

func TestUsers_HandleUserCreate(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getUserHandlerMocks()

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
		{"Invalid syntax", "{123]", http.StatusBadRequest, "Syntax error"},
		{"Empty username", "{\"password\": \"12345678\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty username", "{\"userName\": \"\", \"password\": \"12345678\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty password", "{\"userName\": \"testuser\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty password", "{\"userName\": \"username\", \"password\": \"\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty username and password", "{\"userName\": \"\", \"password\": \"\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty username and password", "{}", http.StatusBadRequest, "Empty username or password"},
		{"Short password", "{\"userName\": \"username\", \"password\": \"123\"}", http.StatusBadRequest, "Password should be at least 8 characters"},
		{"Username conflict", "{\"userName\": \"existingUser\", \"password\": \"12345678\"}", http.StatusBadRequest, "User with username existingUser already exists"},
		{"Storage operation error", "{\"userName\": \"storageErrorUser\", \"password\": \"12345678\"}", http.StatusBadRequest, "Could not create user storageErrorUser"},
		{"Successful user creation", "{\"userName\": \"newUser\", \"password\": \"12345678\"}", http.StatusOK, ""},
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
	loggerMock, userRepoMock, tokenRepoMock := getUserHandlerMocks()

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
		{"Invalid syntax", "{123]", http.StatusBadRequest, "Syntax error"},
		{"Empty username", "{\"password\": \"12345678\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty username", "{\"userName\": \"\", \"password\": \"12345678\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty password", "{\"userName\": \"testuser\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty password", "{\"userName\": \"username\", \"password\": \"\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty username and password", "{\"userName\": \"\", \"password\": \"\"}", http.StatusBadRequest, "Empty username or password"},
		{"Empty username and password", "{}", http.StatusBadRequest, "Empty username or password"},
		{"Non-existing user", "{\"userName\": \"newUser\", \"password\": \"12345678\"}", http.StatusBadRequest, "User newUser does not exist"},
		{"Valid User", "{\"userName\": \"existingUser\", \"password\": \"12345678\"}", http.StatusOK, ""},
		{"Incorrect password", "{\"userName\": \"existingUser\", \"password\": \"1234567890\"}", http.StatusBadRequest, "Invalid username/password"},
		{"Token Storage Error", "{\"userName\": \"tokenStorageError\", \"password\": \"12345678\"}", http.StatusInternalServerError, "Could not generate one-time token"},
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
}

func getUserHandlerMocks() (*mocks.FieldLogger, *mocks.UserRepository, *mocks.TokenRepository) {
	loggerMock := &mocks.FieldLogger{}
	userRepoMock := &mocks.UserRepository{}
	tokenRepoMock := &mocks.TokenRepository{}

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"existingUser",
	).Return(models.User{ID: "uuid", UserName: "existingUser", PasswordHash: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f"}, nil)

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"tokenStorageError",
	).Return(models.User{ID: "uuid-token-storage-error", UserName: "tokenStorageError", PasswordHash: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f"}, nil)

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"storageErrorUser",
	).Return(models.User{ID: "uuid"}, errors.New("storage error"))

	userRepoMock.On(
		"GetByUserName",
		mock.Anything,
		"newUser",
	).Return(models.User{}, errors.New("could not find user"))

	userRepoMock.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(u *models.User) bool { return u.UserName == "storageErrorUser" }),
	).Return(errors.New("storage error"))

	userRepoMock.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(u *models.User) bool { return u.UserName != "storageErrorUser" }),
	).Return(nil)

	tokenRepoMock.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(t *models.Token) bool { return t.UserId == "uuid-token-storage-error" }),
	).Return(errors.New("storage error"))

	tokenRepoMock.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	return loggerMock, userRepoMock, tokenRepoMock
}
