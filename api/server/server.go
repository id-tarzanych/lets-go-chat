package server

import (
	"fmt"
	"github.com/id-tarzanych/lets-go-chat/configurations"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/id-tarzanych/lets-go-chat/api/handlers"
	"github.com/id-tarzanych/lets-go-chat/api/middlewares"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

type Server struct {
	port int

	router *mux.Router
	userRepo user.UserRepository
}

func New(cfg configurations.Configuration, userRepo user.UserRepository) *Server {
	s := &Server{port: cfg.Server.Port, userRepo: userRepo, router: mux.NewRouter()}
	s.routes()

	return s
}

func (s *Server) Handle() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	u := handlers.NewUsers(s.userRepo)

	s.router.HandleFunc("/user", middlewares.PostOnly(u.HandleUserCreate()))
	s.router.HandleFunc("/user/login", middlewares.PostOnly(u.HandleUserLogin()))
}
