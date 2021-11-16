package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	repo user.UserRepository
}

const rateLimit = 100
const tokenDuration = time.Hour

func NewUsers(repo user.UserRepository) *Users {
	return &Users{repo: repo}
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

		if _, err := s.repo.GetByUserName(nil, username); err == nil {
			http.Error(w, fmt.Sprintf("User with username %s already exists", username), http.StatusBadRequest)
			return
		}

		user := models.NewUser(username, password)
		if err := s.repo.Create(nil, user); err != nil {
			http.Error(w, fmt.Sprintf("Could not create user %s", username), http.StatusBadRequest)
			return
		}

		respBody := struct {
			Id       types.Uuid `json:"id"`
			UserName string     `json:"userName"`
		}{user.Id(), user.UserName()}

		js, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
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

		user, err = s.repo.GetByUserName(nil, username)
		if err != nil {
			http.Error(w, fmt.Sprintf("User %s does not exist", username), http.StatusBadRequest)
			return
		}

		if !hasher.CheckPasswordHash(password, user.PasswordHash()) {
			http.Error(w, "Invalid username/password", http.StatusBadRequest)
			return
		}

		respBody := struct {
			Url string `json:"url"`
		}{fmt.Sprintf("ws://fancy-chat.io/ws&token=%s", generators.RandomString(16))}

		js, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Rate-Limit", strconv.Itoa(rateLimit))
		w.Header().Set("X-Expires-After", time.Now().Add(tokenDuration).Format(time.RFC1123))
		w.Write(js)
	}
}
