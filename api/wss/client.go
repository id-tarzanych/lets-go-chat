package wss

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientObject struct {
	JoinedAt time.Time `json:joinedAt,omitempty`
	UserName string    `json:joinedAt,omitempty`

	EntryToken      string          `json:"-"`
	IPAddress       string          `json:"-"`
	ClientWebSocket *websocket.Conn `json:"-"`
	mu              sync.Mutex      `json:"-"`
}

func (c *ClientObject) SendJSON(v interface{}) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.ClientWebSocket.WriteJSON(v)
}

type ClientRequest struct {
	Message string `json:"message"`

	EntryToken string          `json:"-"`
	WebSocket  *websocket.Conn `json:"-"`
}
