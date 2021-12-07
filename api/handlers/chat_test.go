package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/mocks"
	"github.com/id-tarzanych/lets-go-chat/models"
	"github.com/id-tarzanych/lets-go-chat/pkg/generators"
	"github.com/stretchr/testify/mock"
)

func TestChat_HandleActiveUsers(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getChatHandlerMocks()

	tests := []struct {
		name            string
		data            *wss.ChatData
		wantActiveUsers int
	}{
		{"0 users", generateClientsData(0), 0},
		{"5 users", generateClientsData(5), 5},
		{"10 users", generateClientsData(10), 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := &Chat{
				logger: loggerMock,
				upgrader: websocket.Upgrader{
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
				},
				data:      tt.data,
				userRepo:  userRepoMock,
				tokenRepo: tokenRepoMock,
			}

			w := httptest.NewRecorder()

			handlers.HandleActiveUsers().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/test", nil))
			response := &struct {
				Count int `json:"count"`
			}{}

			json.Unmarshal([]byte(w.Body.String()), response)

			if tt.wantActiveUsers != response.Count {
				t.Errorf("Invalid active users count, expected %d, got %d", tt.wantActiveUsers, response.Count)
			}
		})
	}
}

func TestChat_HandleChatSession_WebsocketInitiationError(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getChatHandlerMocks()
	loggerMock.On("Error", mock.AnythingOfType("string")).Return()

	upgrader := websocket.Upgrader{}

	handlers := &Chat{
		logger:    loggerMock,
		upgrader:  upgrader,
		data:      nil,
		userRepo:  userRepoMock,
		tokenRepo: tokenRepoMock,
	}

	w := httptest.NewRecorder()

	handlers.HandleChatSession().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/test", nil))

	loggerMock.AssertCalled(t, "Error", "Could not initiate WebSocket connection.")
}

func TestChat_HandleChatSession_ProcessValidMessage(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getChatHandlerMocks()
	loggerMock.On("Error", mock.AnythingOfType("string")).Return()
	loggerMock.On("Println", mock.Anything).Return()
	loggerMock.On("Println", mock.Anything, mock.Anything).Return()
	loggerMock.On("Println", mock.Anything, mock.Anything, mock.Anything).Return()
	loggerMock.On("Warningln", mock.AnythingOfType("string")).Return()

	user := *models.NewUser("testuser", "12345678")
	tokenString := generators.RandomString(16)

	tokenRepoMock.On(
		"Get",
		mock.Anything,
		tokenString,
	).Return(models.Token{Token: tokenString, UserId: user.ID, Expiration: time.Now().Add(time.Hour * 24)}, nil)

	tokenRepoMock.On("Delete", mock.Anything, mock.Anything).Return(nil)

	userRepoMock.On("GetById", mock.Anything, user.ID).Return(user, nil)

	handlers := &Chat{
		logger: loggerMock,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		data:      wss.NewChatData(),
		userRepo:  userRepoMock,
		tokenRepo: tokenRepoMock,
	}

	s := httptest.NewServer(handlers.HandleChatSession())
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http") + "?token=" + tokenString

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("{\"message\": \"%s\"}", generators.RandomString(16))

		if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			t.Fatalf("%v", err)
		}

		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		requestMessage := &wss.ClientRequest{}
		responseMessage := &wss.ClientRequest{}

		json.Unmarshal([]byte(message), requestMessage)
		err = json.Unmarshal(p, responseMessage)

		if !reflect.DeepEqual(requestMessage, responseMessage) {
			t.Fatalf("bad message")
		}
	}
}

func getChatHandlerMocks() (*mocks.FieldLogger, *mocks.UserRepository, *mocks.TokenRepository) {
	loggerMock := &mocks.FieldLogger{}
	userRepoMock := &mocks.UserRepository{}
	tokenRepoMock := &mocks.TokenRepository{}

	return loggerMock, userRepoMock, tokenRepoMock
}

func generateClientsData(count int) *wss.ChatData {
	data := wss.NewChatData()

	for i := 0; i < count; i++ {
		username := generators.RandomString(8)
		token := generators.RandomString(16)

		client := &wss.ClientObject{
			JoinedAt:        time.Now(),
			UserName:        username,
			EntryToken:      token,
			IPAddress:       "1.1.1.1",
			ClientWebSocket: nil,
		}

		data.Clients[client] = true
		data.ClientTokenMap[token] = client
	}

	return data
}