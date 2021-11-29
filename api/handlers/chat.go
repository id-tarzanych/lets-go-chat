package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"github.com/id-tarzanych/lets-go-chat/db/user"
)

type Chat struct {
	logger *logrus.Logger

	upgrader websocket.Upgrader
	data     *wss.ChatData

	userRepo  user.UserRepository
	tokenRepo token.TokenRepository
}

func NewChat(logger *logrus.Logger, upgrader websocket.Upgrader, data *wss.ChatData, userRepo user.UserRepository, tokenRepo token.TokenRepository) *Chat {
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
		}{len(c.data.Clients)}

		js, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err = w.Write(js); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (c *Chat) HandleChatSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		ws, err := c.upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
			return
		}

		var preListen *wss.ClientObject
		defer func() {
			var e error

			if preListen != nil {
				e = preListen.ClientWebSocket.Close()
				delete(c.data.Clients, preListen)
			}

			fmt.Println(e)
		}()

		for {
			var newElement wss.ClientRequest

			_, p, err := ws.ReadMessage()
			if err != nil {
				c.logger.Println("Client Disconnected: ", err, preListen.EntryToken)

				break
			}

			err = json.Unmarshal(p, &newElement)
			newElement.EntryToken = r.URL.Query().Get("token")
			newElement.WebSocket = ws

			err2, retrievedClient := c.HandleClientMessage(newElement)
			preListen = retrievedClient

			if err2 != nil {
				if retrievedClient != nil {
					c.chuckClient(retrievedClient)
				}
				break
			} else {
				// Pong original message.
				preListen.ClientWebSocket.WriteJSON(newElement)
			}
		}
	}
}

func (c *Chat) chuckClient(object *wss.ClientObject) {
	delete(c.data.Clients, object)
	delete(c.data.ClientTokenMap, object.EntryToken)
}

func (c *Chat) HandleClientMessage(clientData wss.ClientRequest) (error, *wss.ClientObject) {
	c.logger.Println("Entry Token is : ", clientData.EntryToken)

	clientObj, found := c.data.ClientTokenMap[clientData.EntryToken]
	if found == true {
		// Update mapped client's web socket.
		delete(c.data.Clients, clientObj)

		clientObj.ClientWebSocket = clientData.WebSocket
		c.data.Clients[clientObj] = true

		return nil, clientObj
	} else {
		t, err := c.tokenRepo.Get(nil, clientData.EntryToken)
		if err != nil {
			return err, nil
		}

		u, err := c.userRepo.GetById(nil, t.UserId)
		if err != nil {
			return err, nil
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
			return err, nil
		}

		// Map entryToken to client object
		c.data.ClientTokenMap[clientData.EntryToken] = clientObject

		// Map clientObject to a boolean true for easy broadcast
		c.data.Clients[clientObject] = true

		return nil, clientObject
	}
}
