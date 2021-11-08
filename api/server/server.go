package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/id-tarzanych/lets-go-chat/api/handlers"
	"github.com/id-tarzanych/lets-go-chat/api/middlewares"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

type Server struct {
	router  *mux.Router
	userDao *inmemory.UserDao
}

const RateLimit = 100
const TokenDuration = time.Hour

func New(dao *inmemory.UserDao) *Server {
	s := &Server{mux.NewRouter(), dao}
	s.routes()

	return s
}

func (s *Server) WithUserDao(d *inmemory.UserDao) *Server {
	s.userDao = d

	return s
}

func (s *Server) Handle() {
	log.Fatal(http.ListenAndServe(":8080", s.router))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	u := handlers.Users{dao}

	s.router.HandleFunc("/user", middlewares.GetOnly(u.HandleUserCreate()))
	s.router.HandleFunc("/user/login", middlewares.PostOnly(u.HandleUserLogin())
}