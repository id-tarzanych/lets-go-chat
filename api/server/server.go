package server

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/id-tarzanych/lets-go-chat/api/middlewares"
	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/configurations"
	"github.com/id-tarzanych/lets-go-chat/db/message"
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

type Server struct {
	port int

	logger logrus.FieldLogger

	logMiddleware  *middlewares.LogMiddleware
	authMiddleware *middlewares.AuthMiddleware
	router         *mux.Router

	chatData *wss.ChatData
	taskCh   chan WorkerTask

	requestUpgrader websocket.Upgrader

	userRepo    user.UserRepository
	tokenRepo   token.TokenRepository
	messageRepo message.MessageRepository
}

func New(
	cfg configurations.Configuration,
	userRepo user.UserRepository,
	tokenRepo token.TokenRepository,
	messageRepo message.MessageRepository,
	logger logrus.FieldLogger,
) *Server {
	s := &Server{
		port: cfg.Server.Port,

		logger: logger,

		logMiddleware:  middlewares.NewLogMiddleware(logger),
		authMiddleware: middlewares.NewAuthMiddleware(tokenRepo),

		chatData: wss.NewChatData(),
		taskCh:   make(chan WorkerTask),

		requestUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},

		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		messageRepo: messageRepo,
	}

	go broadcastWorker(s.taskCh, s.logger, s.userRepo)

	return s
}

func (s Server) Port() int {
	return s.port
}

func broadcastWorker(taskCh <-chan WorkerTask, logger logrus.FieldLogger, userRepo user.UserRepository) {
	var err error

	for {
		task := <-taskCh

		if err = task.Client.SendMessage(task.Message); err != nil {
			logger.Errorln(err)
			break
		}

		if err = userRepo.UpdateLastActivity(task.Context, task.Client.User, task.Message.CreatedAt); err != nil {
			logger.Errorln(err)
			break
		}
	}
}
