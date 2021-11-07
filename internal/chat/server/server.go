package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/id-tarzanych/lets-go-chat/internal/chat/dao/interfaces"
)

type Server struct {
	router *mux.Router
	userDao *interfaces.UserDao
}

const RateLimit = 100
const TokenDuration = time.Hour

func New(dao *interfaces.UserDao) *Server  {
	s := &Server{mux.NewRouter(), dao}
	s.routes()

	return s
}

func (s *Server) WithUserDao(d *interfaces.UserDao) *Server {
	s.userDao = d

	return s
}

func (s *Server) Handle()  {
	log.Fatal(http.ListenAndServe(":8080", s.router))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
    s.router.ServeHTTP(w, r)
}