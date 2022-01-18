package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/id-tarzanych/lets-go-chat/api/wss"
	"github.com/id-tarzanych/lets-go-chat/mocks"
	"github.com/id-tarzanych/lets-go-chat/models"
	"github.com/id-tarzanych/lets-go-chat/pkg/generators"
)

func TestServer_GetActiveUsers(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getChatHandlerMocks()

	tests := []struct {
		name            string
		data            *wss.ChatData
		wantActiveUsers int
	}{
		{
			name: "0 users",
			data: generateClientsData(0),
		},
		{
			name:            "5 users",
			data:            generateClientsData(5),
			wantActiveUsers: 5,
		},
		{
			name:            "10 users",
			data:            generateClientsData(10),
			wantActiveUsers: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				logger: loggerMock,
				requestUpgrader: websocket.Upgrader{
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
				},
				chatData:  tt.data,
				userRepo:  userRepoMock,
				tokenRepo: tokenRepoMock,
			}

			w := httptest.NewRecorder()

			srv.GetActiveUsers(w, httptest.NewRequest(http.MethodPost, "/test", nil))
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

func TestServer_WebsocketInitiationError(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getChatHandlerMocks()
	loggerMock.On("Error", mock.AnythingOfType("string")).Return()

	upgrader := websocket.Upgrader{}

	srv := &Server{
		logger:          loggerMock,
		requestUpgrader: upgrader,
		chatData:        nil,
		userRepo:        userRepoMock,
		tokenRepo:       tokenRepoMock,
	}

	w := httptest.NewRecorder()

	srv.WsRTMStart(w, httptest.NewRequest(http.MethodPost, "/test", nil), WsRTMStartParams{})

	loggerMock.AssertCalled(t, "Error", "Could not initiate WebSocket connection.")
}

func TestChat_HandleChatSession_ProcessValidMessage(t *testing.T) {
	loggerMock, userRepoMock, tokenRepoMock := getChatHandlerMocks()

	loggerMock.On("Error", mock.AnythingOfType("string")).Maybe().Return()
	loggerMock.On("Println", mock.Anything).Maybe().Return()
	loggerMock.On("Println", mock.Anything, mock.Anything).Return()
	loggerMock.On("Println", mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	loggerMock.On("Warningln", mock.AnythingOfType("string")).Maybe().Return()

	user := *models.NewUser("testuser", "12345678")
	tokenString := generators.RandomString(16)

	tokenRepoMock.On("Get", mock.Anything, tokenString).Return(models.Token{Token: tokenString, UserId: user.ID, Expiration: time.Now().Add(time.Hour * 24)}, nil)
	tokenRepoMock.On("Delete", mock.Anything, mock.Anything).Return(nil)

	userRepoMock.On("GetById", mock.Anything, user.ID).Return(user, nil)

	srv := &Server{
		logger: loggerMock,
		requestUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		chatData:  wss.NewChatData(),
		userRepo:  userRepoMock,
		tokenRepo: tokenRepoMock,
	}

	s := httptest.NewServer(HandlerWithOptions(srv, ChiServerOptions{}))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/chat/ws.rtm.start?token=" + tokenString

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

		assert.NoError(t, err, "json should be valid")
		assert.Equal(t, requestMessage, responseMessage, "objects should be equal")
	}

	loggerMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
	tokenRepoMock.AssertExpectations(t)
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

		client := &wss.Client{
			JoinedAt:   time.Now(),
			User:       models.NewUser(username, "password"),
			EntryToken: token,
			IPAddress:  "1.1.1.1",
			WebSocket:  nil,
		}

		data.Clients[client] = true
		data.ClientTokens[token] = client
	}

	return data
}
