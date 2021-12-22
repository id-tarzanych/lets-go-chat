package middlewares

import (
	"net/http"
)

func GetOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed.", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func PostOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed.", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
