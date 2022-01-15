package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/db/message"
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"github.com/id-tarzanych/lets-go-chat/db/user"
	"github.com/id-tarzanych/lets-go-chat/models"
)

type Chat struct {
	logger logrus.FieldLogger

	upgrader websocket.Upgrader
	data     *wss.ChatData

	userRepo    user.UserRepository
	tokenRepo   token.TokenRepository
	messageRepo message.MessageRepository
}

type WorkerTask struct {
	Client  *wss.Client
	Message *models.Message
}

func NewChat(
	logger logrus.FieldLogger,
	upgrader websocket.Upgrader,
	data *wss.ChatData,
	userRepo user.UserRepository,
	tokenRepo token.TokenRepository,
	messageRepo message.MessageRepository,
) *Chat {
	chat := &Chat{
		logger: logger,

		upgrader: upgrader,
		data:     data,

		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		messageRepo: messageRepo,
	}

	return chat
}

func (c *Chat) HandleActiveUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respBody := struct {
			Count int `json:"count"`
		}{Count: len(c.data.Clients)}

		js, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Chat) HandleChatSession() http.HandlerFunc {
	ctxChat := context.TODO()

	taskCh := make(chan WorkerTask)
	go broadcastWorker(ctxChat, taskCh, c.logger, c.userRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()

		c.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		ws, err := c.upgrader.Upgrade(w, r, nil)
		if err != nil {
			c.logger.Error("Could not initiate WebSocket connection.")

			return
		}

		token := r.URL.Query().Get("token")

		var preListen *wss.Client
		defer func() {
			if preListen == nil {
				return
			}

			if err := preListen.WebSocket.Close(); err != nil {
				c.logger.Println(err)
			}

			c.data.DeleteClient(preListen)
		}()

		for {
			retrievedClient, err := c.retrieveClient(ctx, token, ws)

			if err != nil {
				if retrievedClient != nil {
					c.chuckClient(retrievedClient)

					break
				}

				c.logger.Error("Could not initiate WebSocket connection.")

				return
			}

			// Update pre-listen object with valid client.
			preListen = retrievedClient

			var newElement wss.ClientRequest

			_, p, err := ws.ReadMessage()
			if err != nil {
				c.logger.Println("Client Disconnected: ", err, preListen.EntryToken)

				break
			}

			if err = json.Unmarshal(p, &newElement); err != nil {
				c.logger.Warningln("Invalid request. ", err, p)

				break
			}

			newElement.EntryToken = r.URL.Query().Get("token")
			newElement.WebSocket = ws

			m := &models.Message{Author: *preListen.User, Message: newElement.Message}
			if err := c.messageRepo.Create(ctx, m); err != nil {
				return
			}

			// Broadcast message.
			c.broadcastMessage(taskCh, m)
		}
	}
}

func (c *Chat) chuckClient(client *wss.Client) {
	client.Stop()

	c.data.DeleteClient(client)
	c.data.DeleteToken(client.EntryToken)
}

func (c *Chat) retrieveClient(ctx context.Context, token string, ws *websocket.Conn) (*wss.Client, error) {
	c.logger.Println("Entry Token is : ", token)

	if clientObj := c.data.LoadClient(token); clientObj != nil {
		// Update mapped client's web socket.
		c.data.DeleteClient(clientObj)

		clientObj.WebSocket = ws
		c.data.StoreClient(clientObj)

		return clientObj, nil
	}

	t, err := c.tokenRepo.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	u, err := c.userRepo.GetById(ctx, t.UserId)
	if err != nil {
		return nil, err
	}

	clientObject := wss.NewClientObject(time.Now(), &u, token, ws)

	// Invalidate token.
	if err = c.tokenRepo.Delete(ctx, t.Token); err != nil {
		return nil, err
	}

	// Retrieve missed messages.
	var missedMessages []models.Message
	if clientObject.User.LastActivity.IsZero() {
		missedMessages, err = c.messageRepo.GetAll(ctx)
	} else {
		missedMessages, err = c.messageRepo.GetNewerThan(ctx, clientObject.User.LastActivity)
	}

	if err != nil {
		return nil, err
	}

	var lastMessage models.Message
	for i := range missedMessages {
		lastMessage = missedMessages[i]
		clientObject.SendMessage(&lastMessage)
	}

	if len(missedMessages) > 0 {
		if err := c.userRepo.UpdateLastActivity(ctx, clientObject.User, lastMessage.CreatedAt); err != nil {
			return nil, err
		}
	}

	// Map entryToken to client object
	c.data.StoreToken(token, clientObject)

	// Map clientObject to a boolean true for easy broadcast
	c.data.StoreClient(clientObject)

	return clientObject, nil
}

func (c *Chat) broadcastMessage(tasksCh chan WorkerTask, m *models.Message) {
	for client := range c.data.Clients {
		go func(c *wss.Client) {
			tasksCh <- WorkerTask{
				Client:  c,
				Message: m,
			}
		}(client)
	}
}

func broadcastWorker(ctx context.Context, taskCh <-chan WorkerTask, logger logrus.FieldLogger, userRepo user.UserRepository) {
	var err error

	for {
		task := <-taskCh

		if err = task.Client.SendMessage(task.Message); err != nil {
			logger.Errorln(err)
			break
		}

		if err = userRepo.UpdateLastActivity(ctx, task.Client.User, task.Message.CreatedAt); err != nil {
			logger.Errorln(err)
			break
		}
	}
}
