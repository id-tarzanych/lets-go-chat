package user

import (
	"log"
)

type DaoType string

func MakeUserDao(t DaoType) UserDao {
	switch t {
	case "inmemory":
		return inmemory.NewUserDaoInMemory()
	default:
		log.Println("Fallback to default In Memory DAO")

		return inmemory.NewUserDaoInMemory()
	}
}
