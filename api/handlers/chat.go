package handlers

import (
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

const WORKERS = 50

type Chat struct {
	logger logrus.FieldLogger

	upgrader   websocket.Upgrader
	data       *wss.ChatData
	activityCh chan *models.User

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

		activityCh: make(chan *models.User),
	}

	go func() {
		for {
			updatedUser := <-chat.activityCh
			userRepo.Update(nil, updatedUser)
		}
	}()

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
	// Initialize bworkers pool
	taskCh := make(chan WorkerTask)
	for i := 1; i <= WORKERS; i++ {
		go broadcastWorker(taskCh, c.activityCh)
	}

	return func(w http.ResponseWriter, r *http.Request) {
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

			delete(c.data.Clients, preListen)
		}()

		for {
			retrievedClient, errHandle := c.retrieveClient(token, ws)
			preListen = retrievedClient

			if errHandle != nil {
				if retrievedClient != nil {
					c.chuckClient(retrievedClient)
				} else {
					c.logger.Error("Could not initiate WebSocket connection.")

					return
				}

				break
			}

			var newElement wss.ClientRequest

			_, p, err := ws.ReadMessage()
			if err != nil {
				c.logger.Println("Client Disconnected: ", err, preListen.EntryToken)

				break
			}

			err = json.Unmarshal(p, &newElement)
			if err != nil {
				c.logger.Warningln("Invalid request. ", err, p)

				break
			}

			newElement.EntryToken = r.URL.Query().Get("token")
			newElement.WebSocket = ws

			m := &models.Message{Author: *retrievedClient.User, Message: newElement.Message}
			c.messageRepo.Create(nil, m)

			// Broadcast message.
			c.broadcastMessage(taskCh, m)
		}
	}
}

func (c *Chat) chuckClient(client *wss.Client) {
	client.Stop()

	delete(c.data.Clients, client)
	delete(c.data.ClientTokenMap, client.EntryToken)
}

func (c *Chat) retrieveClient(token string, ws *websocket.Conn) (*wss.Client, error) {
	c.logger.Println("Entry Token is : ", token)

	if clientObj, found := c.data.ClientTokenMap[token]; found == true {
		// Update mapped client's web socket.
		delete(c.data.Clients, clientObj)

		clientObj.WebSocket = ws
		c.data.Clients[clientObj] = true

		return clientObj, nil
	}

	t, err := c.tokenRepo.Get(nil, token)
	if err != nil {
		return nil, err
	}

	u, err := c.userRepo.GetById(nil, t.UserId)
	if err != nil {
		return nil, err
	}

	clientObject := wss.NewClientObject(time.Now(), &u, token, ws)

	// Invalidate token.
	if err = c.tokenRepo.Delete(nil, t.Token); err != nil {
		return nil, err
	}

	// Retrieve missed messages.
	var missedMessages []models.Message
	if clientObject.User.LastActivity.IsZero() {
		missedMessages, _ = c.messageRepo.GetAll(nil)
	} else {
		missedMessages, _ = c.messageRepo.GetNewerThan(nil, clientObject.User.LastActivity)
	}

	var lastMessage models.Message
	for i := range missedMessages {
		lastMessage = missedMessages[i]
		clientObject.SendMessage(&lastMessage)
	}

	if len(missedMessages) > 0 {
		clientObject.User.LastActivity = lastMessage.CreatedAt
		c.activityCh <- clientObject.User
	}

	// Map entryToken to client object
	c.data.ClientTokenMap[token] = clientObject

	// Map clientObject to a boolean true for easy broadcast
	c.data.Clients[clientObject] = true

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

func broadcastWorker(taskCh <-chan WorkerTask, activityCh chan *models.User) {
	for {
		task := <-taskCh

		task.Client.SendMessage(task.Message)
		task.Client.User.LastActivity = task.Message.CreatedAt
		activityCh <- task.Client.User
	}
}
