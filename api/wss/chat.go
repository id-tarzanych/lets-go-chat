package wss

import "sync"

type ChatData struct {
	Clients      map[*Client]bool
	ClientTokens map[string]*Client

	mu sync.Mutex
}

func NewChatData() *ChatData {
	return &ChatData{
		Clients:      make(map[*Client]bool),
		ClientTokens: make(map[string]*Client),
	}
}

func (c *ChatData) ClientExists(client *Client) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.Clients[client]

	if !ok {
		return ok
	}

	return value
}

func (c *ChatData) StoreClient(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Clients[client] = true
}

func (c *ChatData) DeleteClient(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Clients, client)
}

func (c *ChatData) GetAllClients() []Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	clients := make([]Client, 0)
	for client := range c.Clients {
		clients = append(clients, *client)
	}

	return clients
}

func (c *ChatData) LoadClient(token string) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, ok := c.ClientTokens[token]
	if !ok {
		return nil
	}

	return client
}

func (c *ChatData) StoreToken(token string, client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ClientTokens[token] = client
}

func (c *ChatData) DeleteToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.ClientTokens, token)
}