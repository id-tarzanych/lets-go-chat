package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/id-tarzanych/lets-go-chat/internal/testserver"
	"github.com/id-tarzanych/lets-go-chat/mocks"
	"github.com/id-tarzanych/lets-go-chat/models"
)

const (
	validToken   = "validtoken"
	invalidToken = "invalidToken"
	expiredToken = "expiredToken"
)

func TestAuthMiddleware_ValidateToken(t *testing.T) {
	s := &testserver.Server{}

	tests := []struct {
		name        string
		req         *http.Request
		wantCode    int
		wantMessage string
	}{
		{
			name:        "No Token",
			req:         httptest.NewRequest(http.MethodGet, "/test", nil),
			wantCode:    http.StatusBadRequest,
			wantMessage: "Access token is required.",
		},
		{
			name:     "Valid Token",
			req:      httptest.NewRequest(http.MethodGet, "/test?token="+validToken, nil),
			wantCode: http.StatusOK,
		},
		{
			name:        "Invalid Token",
			req:         httptest.NewRequest(http.MethodGet, "/test?token="+invalidToken, nil),
			wantCode:    http.StatusBadRequest,
			wantMessage: "Access token is invalid.",
		},
		{
			name:        "Expired Token",
			req:         httptest.NewRequest(http.MethodGet, "/test?token="+expiredToken, nil),
			wantCode:    http.StatusBadRequest,
			wantMessage: "Access token expired.",
		},
	}

	tokenRepoMock := &mocks.TokenRepository{}

	tokenRepoMock.On(
		"Get",
		mock.Anything,
		validToken,
	).Return(*models.NewToken(validToken, "uuid", time.Now().Add(time.Hour*24)), nil)

	tokenRepoMock.On(
		"Get",
		mock.Anything,
		invalidToken,
	).Return(models.Token{}, errors.New("token not found"))

	tokenRepoMock.On(
		"Get",
		mock.Anything,
		expiredToken,
	).Return(*models.NewToken(expiredToken, "uuid", time.Now().Add(-24*time.Hour)), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			authMiddleware := &AuthMiddleware{
				tokenRepo: tokenRepoMock,
			}

			authMiddleware.ValidateToken(s).ServeHTTP(w, tt.req)
			result := w.Result()

			if tt.wantCode != result.StatusCode {
				t.Errorf("Incorrect status code, wanted %d, got %d.", tt.wantCode, result.StatusCode)
			}

			responseBody := strings.TrimSpace(w.Body.String())
			if tt.wantCode != http.StatusOK && tt.wantMessage != responseBody {
				t.Errorf("Incorrect error message, wanted \"%s\", got \"%s\"", tt.wantMessage, responseBody)
			}
		})
	}
}
