package testserver

import (
	"net/http"
)

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/badRequest":
		http.Error(w, "Test Bad Request", http.StatusBadRequest)
	case "/internalServerError":
		http.Error(w, "Test Internal Server Error", http.StatusInternalServerError)
	case "/panic":
		panic("Test panic!")
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("It works!"))
	}
}
