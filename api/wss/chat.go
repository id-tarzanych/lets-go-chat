package wss

type ChatData struct {
	Clients        map[*ClientObject]bool
	ClientTokenMap map[string]*ClientObject
}

func NewChatData() *ChatData {
	return &ChatData{
		Clients:        make(map[*ClientObject]bool),
		ClientTokenMap: make(map[string]*ClientObject),
	}
}
