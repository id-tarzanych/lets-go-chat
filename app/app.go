package app

import "github.com/id-tarzanych/lets-go-chat/db/user"

type App struct {
	user.UserDao
}

func NewApp() *App {
	user.MakeUserDao()

	return &App{UserDao: userDao}
}
