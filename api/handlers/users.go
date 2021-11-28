package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/id-tarzanych/lets-go-chat/db/token"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/id-tarzanych/lets-go-chat/db/user"
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/models"
	"github.com/id-tarzanych/lets-go-chat/pkg/generators"
	"github.com/id-tarzanych/lets-go-chat/pkg/hasher"
)

type Users struct {
	userRepo user.UserRepository
	tokenRepo token.TokenRepository
}

const rateLimit = 100
const tokenDuration = time.Hour

func NewUsers(userRepo user.UserRepository, tokenRepo token.TokenRepository) *Users {
	return &Users{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (s Users) HandleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			UserName string `json:"userName"`
			Password string `json:"password"`
		}

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

		respBody := struct {
			Id       types.Uuid `json:"id"`
			UserName string     `json:"userName"`
		}{user.ID, user.UserName}

		js, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err = w.Write(js); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s Users) HandleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			UserName string `json:"userName"`
			Password string `json:"password"`
		}
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		oneTimeUrl := netUrl.URL{
			Scheme: "ws",
			Host: r.Host,
			Path: "/chat/ws.rtm.start",
			RawQuery: fmt.Sprintf("token=%s", token.Token),
		}

		respBody := struct {
			Url string `json:"url"`
		}{oneTimeUrl.String()}

		js, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Rate-Limit", strconv.Itoa(rateLimit))
		w.Header().Set("X-Expires-After", token.Expiration.Format(time.RFC1123))

		if _, err = w.Write(js); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

