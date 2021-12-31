package wss

type ChatData struct {
	Clients        map[*Client]bool
	ClientTokenMap map[string]*Client
}

func NewChatData() *ChatData {
	return &ChatData{
		Clients:        make(map[*Client]bool),
		ClientTokenMap: make(map[string]*Client),
	}
}
