package wss

import (
	"context"
	"errors"
	"time"

	"github.com/gorilla/websocket"

	"github.com/id-tarzanych/lets-go-chat/models"
)

type Client struct {
	ctx       context.Context
	cancelCtx context.CancelFunc

	JoinedAt time.Time    `json:"joinedAt,omitempty"`
	User     *models.User `json:"-"`

	EntryToken string          `json:"-"`
	IPAddress  string          `json:"-"`
	WebSocket  *websocket.Conn `json:"-"`

	messageCh chan *models.Message
}

func NewClientObject(joinedAt time.Time, user *models.User, entryToken string, webSocket *websocket.Conn) *Client {
	client := &Client{JoinedAt: joinedAt, User: user, EntryToken: entryToken, WebSocket: webSocket}

	client.ctx, client.cancelCtx = context.WithCancel(context.Background())

	client.IPAddress = webSocket.RemoteAddr().String()
	client.messageCh = make(chan *models.Message)

	client.processIncomingMessages()

	return client
}

type ClientRequest struct {
	Message string `json:"message"`

	EntryToken string          `json:"-"`
	WebSocket  *websocket.Conn `json:"-"`
}

func (c *Client) SendMessage(message *models.Message) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("channel is closed")
		}
	}()

	c.messageCh <- message

	return err
}

func (c *Client) Stop() {
	c.cancelCtx()
}

func (c *Client) processIncomingMessages() {
	go func() {
		for {
			select {
			case message := <-c.messageCh:
				output := struct {
					Author  string    `json:"author"`
					Message string    `json:"message"`
					SentAt  time.Time `json:"sentAt"`
				}{Author: message.Author.UserName, Message: message.Message, SentAt: message.CreatedAt}

				c.WebSocket.WriteJSON(output)

			case <-c.ctx.Done():
				return
			}
		}
	}()
}
