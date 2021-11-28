package wss

type ChatData struct {
	Clients            map[*ClientObject]bool
	ClientTokenMap     map[string]*ClientObject
	ClientRoomActivity chan string
}

func NewChatData() *ChatData {
	return &ChatData{
		Clients: make(map[*ClientObject]bool),
		ClientTokenMap: make(map[string]*ClientObject),
		ClientRoomActivity: make(chan string),
	}
}