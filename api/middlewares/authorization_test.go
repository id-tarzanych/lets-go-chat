package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/id-tarzanych/lets-go-chat/internal/testserver"
	"github.com/id-tarzanych/lets-go-chat/mocks"
	"github.com/id-tarzanych/lets-go-chat/models"
	"github.com/stretchr/testify/mock"
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
		{"No Token", httptest.NewRequest(http.MethodGet, "/test", nil), http.StatusBadRequest, "Access token is required."},
		{"Valid Token", httptest.NewRequest(http.MethodGet, "/test?token="+validToken, nil), http.StatusOK, ""},
		{"Invalid Token", httptest.NewRequest(http.MethodGet, "/test?token="+invalidToken, nil), http.StatusBadRequest, "Access token is invalid."},
		{"Expired Token", httptest.NewRequest(http.MethodGet, "/test?token="+expiredToken, nil), http.StatusBadRequest, "Access token expired."},
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
