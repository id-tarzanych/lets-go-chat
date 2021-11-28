package middlewares

import (
	"net/http"
	"time"

	"github.com/id-tarzanych/lets-go-chat/db/token"
)

type AuthMiddleware struct {
	tokenRepo token.TokenRepository
}

func NewAuthMiddleware(tokenRepo token.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{tokenRepo: tokenRepo}
}

func (a AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")

		if tokenString == "" {
			http.Error(w, "Access token is required.", http.StatusBadRequest)
			return
		}

		requestToken, err := a.tokenRepo.Get(nil, tokenString)
		if err != nil {
			http.Error(w, "Access token is invalid.", http.StatusBadRequest)
			return
		}

		if requestToken.Expiration.Before(time.Now()) {
			http.Error(w, "Access token expired.", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
