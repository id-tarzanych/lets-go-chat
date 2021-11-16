package middlewares

import (
	"net/http"
)

func PostOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed.", http.StatusBadRequest)
			return
		}

		h(w, r)
	}
}
