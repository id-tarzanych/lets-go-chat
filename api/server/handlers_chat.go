package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/models"
)

type WorkerTask struct {
	Context context.Context
	Client  *wss.Client
	Message *models.Message
}

func (s Server) WsRTMStart(w http.ResponseWriter, r *http.Request, params WsRTMStartParams) {
	token := params.Token
	ctx := context.WithValue(r.Context(), "token", token)

	s.requestUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := s.requestUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("Could not initiate WebSocket connection.")

		return
	}

	var preListen *wss.Client
	defer func() {
		if preListen == nil {
			return
		}

		if err := preListen.WebSocket.Close(); err != nil {
			s.logger.Println(err)
		}

		s.chatData.DeleteClient(preListen)
	}()

	for {
		retrievedClient, err := s.retrieveClient(ctx, token, ws)

		if err != nil {
			if retrievedClient != nil {
				s.chuckClient(retrievedClient)

				break
			}

			s.logger.Error("Could not initiate WebSocket connection.")

			return
		}

		// Update pre-listen object with valid client.
		preListen = retrievedClient

		var newElement wss.ClientRequest

		_, p, err := ws.ReadMessage()
		if err != nil {
			s.logger.Println("Client Disconnected: ", err, preListen.EntryToken)

			break
		}

		if err = json.Unmarshal(p, &newElement); err != nil {
			s.logger.Warningln("Invalid request. ", err, p)

			break
		}

		newElement.EntryToken = token
		newElement.WebSocket = ws

		m := &models.Message{Author: *preListen.User, Message: newElement.Message}
		if err := s.messageRepo.Create(ctx, m); err != nil {
			return
		}

		// Broadcast message.
		s.broadcastMessage(s.taskCh, ctx, m)
	}
}

func (s Server) GetActiveUsers(w http.ResponseWriter, r *http.Request) {
	respBody := ActiveUsersResponse{Count: len(s.chatData.Clients)}

	js, _ := json.Marshal(respBody)

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(js)
	if err != nil {
		s.logger.Errorln("Could not write response")
	}
}

func (s Server) chuckClient(client *wss.Client) {
	client.Stop()

	s.chatData.DeleteClient(client)
	s.chatData.DeleteToken(client.EntryToken)
}

func (s Server) retrieveClient(ctx context.Context, token string, ws *websocket.Conn) (*wss.Client, error) {
	if clientObj := s.chatData.LoadClient(token); clientObj != nil {
		// Update mapped client's web socket.
		s.chatData.DeleteClient(clientObj)

		clientObj.WebSocket = ws
		s.chatData.StoreClient(clientObj)

		return clientObj, nil
	}

	t, err := s.tokenRepo.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetById(ctx, t.UserId)
	if err != nil {
		return nil, err
	}

	clientObject := wss.NewClientObject(time.Now(), &u, token, ws)

	// Invalidate token.
	if err = s.tokenRepo.Delete(ctx, t.Token); err != nil {
		return nil, err
	}

	// Retrieve missed messages.
	var missedMessages []models.Message
	if clientObject.User.LastActivity.IsZero() {
		missedMessages, err = s.messageRepo.GetAll(ctx)
	} else {
		missedMessages, err = s.messageRepo.GetNewerThan(ctx, clientObject.User.LastActivity)
	}

	if err != nil {
		return nil, err
	}

	var lastMessage models.Message
	for i := range missedMessages {
		lastMessage = missedMessages[i]
		err := clientObject.SendMessage(&lastMessage)
		if err != nil {
			return nil, err
		}
	}

	if len(missedMessages) > 0 {
		if err := s.userRepo.UpdateLastActivity(ctx, clientObject.User, lastMessage.CreatedAt); err != nil {
			return nil, err
		}
	}

	// Map entryToken to client object
	s.chatData.StoreToken(token, clientObject)

	// Map clientObject to a boolean true for easy broadcast
	s.chatData.StoreClient(clientObject)

	return clientObject, nil
}

func (s Server) broadcastMessage(tasksCh chan WorkerTask, ctx context.Context, m *models.Message) {
	for client := range s.chatData.Clients {
		go func(c *wss.Client) {
			tasksCh <- WorkerTask{
				Context: ctx,
				Client:  c,
				Message: m,
			}
		}(client)
	}
}
