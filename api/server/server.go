package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/id-tarzanych/lets-go-chat/api/handlers"
	"github.com/id-tarzanych/lets-go-chat/api/middlewares"
	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/configurations"
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"github.com/id-tarzanych/lets-go-chat/db/user"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	port int

	logger logrus.FieldLogger

	logMiddleware  *middlewares.LogMiddleware
	authMiddleware *middlewares.AuthMiddleware
	router         *mux.Router

	chatData *wss.ChatData

	requestUpgrader websocket.Upgrader

	userRepo  user.UserRepository
	tokenRepo token.TokenRepository
}

func New(cfg configurations.Configuration, userRepo user.UserRepository, tokenRepo token.TokenRepository, logger logrus.FieldLogger) *Server {
	s := &Server{
		port: cfg.Server.Port,

		logger: logger,

		logMiddleware:  middlewares.NewLogMiddleware(logger),
		authMiddleware: middlewares.NewAuthMiddleware(tokenRepo),
		router:         mux.NewRouter(),

		chatData: wss.NewChatData(),

		requestUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},

		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}

	s.routes()

	return s
}

func (s *Server) Handle() {
	s.logger.Error(http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router))
}

func (s *Server) routes() {
	userHandlers := handlers.NewUsers(s.logger, s.userRepo, s.tokenRepo)
	chatHandlers := handlers.NewChat(s.logger, s.requestUpgrader, s.chatData, s.userRepo, s.tokenRepo)

	commonMW := alice.New(s.logMiddleware.LogError, s.logMiddleware.LogRequest, s.logMiddleware.LogPanicRecovery)

	s.router.Handle("/user", commonMW.Append(middlewares.PostOnly).ThenFunc(userHandlers.HandleUserCreate()))
	s.router.Handle("/user/login", commonMW.Append(middlewares.PostOnly).ThenFunc(userHandlers.HandleUserLogin()))

	s.router.Handle("/user/active", commonMW.Append(middlewares.GetOnly).ThenFunc(chatHandlers.HandleActiveUsers()))
	s.router.Handle("/chat/ws.rtm.start", commonMW.Append(s.authMiddleware.ValidateToken).Then(chatHandlers.HandleChatSession()))
}
