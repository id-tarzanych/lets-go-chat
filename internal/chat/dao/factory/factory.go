package factory

import (
	"log"

	"github.com/id-tarzanych/lets-go-chat/internal/chat/dao/inmemory"
	"github.com/id-tarzanych/lets-go-chat/internal/chat/dao/interfaces"
	"github.com/id-tarzanych/lets-go-chat/internal/chat/types"
)

func MakeUserDao(t types.DaoType) interfaces.UserDao {
	switch t {
	case "inmemory":
		return inmemory.NewUserDaoInMemory()
	default:
		log.Println("Fallback to default In Memory DAO")

		return inmemory.NewUserDaoInMemory()
	}
}