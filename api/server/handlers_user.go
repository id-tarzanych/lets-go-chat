package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/id-tarzanych/lets-go-chat/models"
	"github.com/id-tarzanych/lets-go-chat/pkg/generators"
	"github.com/id-tarzanych/lets-go-chat/pkg/hasher"
)

const rateLimit = 100
const tokenDuration = time.Hour

func (s Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateUserJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Syntax error", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(reqBody.UserName)
	password := strings.TrimSpace(reqBody.Password)
	if username == "" || password == "" {
		http.Error(w, "Empty username or password", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "Password should be at least 8 characters", http.StatusBadRequest)
		return
	}

	if _, err := s.userRepo.GetByUserName(nil, username); err == nil {
		http.Error(w, fmt.Sprintf("User with username %s already exists", username), http.StatusBadRequest)
		return
	}

	user := models.NewUser(username, password)
	if err := s.userRepo.Create(nil, user); err != nil {
		http.Error(w, fmt.Sprintf("Could not create user %s", username), http.StatusBadRequest)
		return
	}

	userId := string(user.ID)
	respBody := CreateUserResponse{Id: &userId, UserName: &user.UserName}

	w.Header().Set("Content-Type", "application/json")
	js, _ := json.Marshal(respBody)
	_, err = w.Write(js)
	if err != nil {
		s.logger.Errorln("Could not write response")
	}
}

func (s Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	var reqBody LoginUserJSONRequestBody
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Syntax error", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(reqBody.UserName)
	password := strings.TrimSpace(reqBody.Password)
	if username == "" || password == "" {
		http.Error(w, "Empty username or password", http.StatusBadRequest)
		return
	}

	user, err = s.userRepo.GetByUserName(nil, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("User %s does not exist", username), http.StatusBadRequest)
		return
	}

	if !hasher.CheckPasswordHash(password, user.PasswordHash) {
		http.Error(w, "Invalid username/password", http.StatusBadRequest)
		return
	}

	token := models.NewToken(
		generators.RandomString(16),
		user.ID,
		time.Now().Add(tokenDuration),
	)

	err = s.tokenRepo.Create(nil, token)
	if err != nil {
		http.Error(w, "Could not generate one-time token", http.StatusInternalServerError)
		return
	}

	oneTimeUrl := netUrl.URL{
		Scheme:   "ws",
		Host:     r.Host,
		Path:     "/chat/ws.rtm.start",
		RawQuery: fmt.Sprintf("token=%s", token.Token),
	}

	respBody := LoginUserResponse{Url: oneTimeUrl.String()}

	js, _ := json.Marshal(respBody)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Rate-Limit", strconv.Itoa(rateLimit))
	w.Header().Set("X-Expires-After", token.Expiration.Format(time.RFC1123))
	_, err = w.Write(js)
	if err != nil {
		s.logger.Errorln("Could not write response")
	}
}
