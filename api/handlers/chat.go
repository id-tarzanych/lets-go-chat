package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"github.com/id-tarzanych/lets-go-chat/db/user"
	"github.com/sirupsen/logrus"
)

type Chat struct {
	logger logrus.FieldLogger

	upgrader websocket.Upgrader
	data     *wss.ChatData

	userRepo  user.UserRepository
	tokenRepo token.TokenRepository
}

func NewChat(logger logrus.FieldLogger, upgrader websocket.Upgrader, data *wss.ChatData, userRepo user.UserRepository, tokenRepo token.TokenRepository) *Chat {
	return &Chat{
		logger: logger,

		upgrader: upgrader,
		data:     data,

		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
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
	return func(w http.ResponseWriter, r *http.Request) {
		c.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		ws, err := c.upgrader.Upgrade(w, r, nil)
		if err != nil {
			c.logger.Error("Could not initiate WebSocket connection.")

			return
		}

		var preListen *wss.ClientObject
		defer func() {
			if preListen == nil {
				return
			}

			if err := preListen.ClientWebSocket.Close(); err != nil {
				c.logger.Println(err)
			}

			delete(c.data.Clients, preListen)
		}()

		for {
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

			retrievedClient, errHandle := c.handleClientMessage(newElement)
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

			// Pong original message.
			preListen.ClientWebSocket.WriteJSON(newElement)
		}
	}
}

func (c *Chat) chuckClient(object *wss.ClientObject) {
	delete(c.data.Clients, object)
	delete(c.data.ClientTokenMap, object.EntryToken)
}

func (c *Chat) handleClientMessage(clientData wss.ClientRequest) (*wss.ClientObject, error) {
	c.logger.Println("Entry Token is : ", clientData.EntryToken)

	if clientObj, found := c.data.ClientTokenMap[clientData.EntryToken]; found == true {
		// Update mapped client's web socket.
		delete(c.data.Clients, clientObj)

		clientObj.ClientWebSocket = clientData.WebSocket
		c.data.Clients[clientObj] = true

		return clientObj, nil
	}

	t, err := c.tokenRepo.Get(nil, clientData.EntryToken)
	if err != nil {
		return nil, err
	}

	u, err := c.userRepo.GetById(nil, t.UserId)
	if err != nil {
		return nil, err
	}

	clientObject := &wss.ClientObject{
		UserName:        u.UserName,
		ClientWebSocket: clientData.WebSocket,
		IPAddress:       clientData.WebSocket.RemoteAddr().String(),
		EntryToken:      clientData.EntryToken,
		JoinedAt:        time.Now(),
	}

	// Invalidate token.
	err = c.tokenRepo.Delete(nil, t.Token)
	if err != nil {
		return nil, err
	}

	// Map entryToken to client object
	c.data.ClientTokenMap[clientData.EntryToken] = clientObject

	// Map clientObject to a boolean true for easy broadcast
	c.data.Clients[clientObject] = true

	return clientObject, nil
}
