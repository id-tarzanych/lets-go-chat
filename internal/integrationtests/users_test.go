package integrationtests

import (
	"testing"

	"github.com/id-tarzanych/lets-go-chat/internal/testdb"
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"github.com/id-tarzanych/lets-go-chat/models"
)

func Test_GetAllUsers(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	expectedUsersMap, err := testdb.SeedUsers(a.DB())
	if err != nil {
		t.Error("could not seed users")
	}

	users, err := a.UserRepo().GetAll(nil)
	if err != nil {
		t.Error("could not load users from database")
	}

	gotUsersMap := make(map[types.Uuid]models.User)
	for i := range users {
		gotUsersMap[users[i].ID] = users[i]
	}

	if len(expectedUsersMap) != len(gotUsersMap) {
		t.Error("expected users amount does not match received users amount")
	}

	for k, e := range expectedUsersMap {
		u, ok := gotUsersMap[k]
		if !ok {
			t.Errorf("user %s was not returned", e.ID)
		}

		if match := compareUsers(u, e); !match {
			t.Errorf("properties for user %s do not match expected ones. Expected %v, got %v", e.ID, e, u)
		}
	}
}

func Test_GetUserById(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	expectedUsersMap, err := testdb.SeedUsers(a.DB())
	if err != nil {
		t.Error("could not seed users")
	}

	for _, e := range expectedUsersMap {
		u, err := a.UserRepo().GetById(nil, e.ID)
		if err != nil {
			t.Errorf("user %s was not returned", e.ID)
		}

		if match := compareUsers(u, e); !match {
			t.Errorf("properties for user %s do not match expected ones. Expected %v, got %v", e.ID, e, u)
		}
	}
}

func Test_GetUserByUsername(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	expectedUsersMap, err := testdb.SeedUsers(a.DB())
	if err != nil {
		t.Error("could not seed users")
	}

	for _, e := range expectedUsersMap {
		u, err := a.UserRepo().GetByUserName(nil, e.UserName)
		if err != nil {
			t.Errorf("user %s was not returned", e.ID)
		}

		if match := compareUsers(u, e); !match {
			t.Errorf("properties for user %s do not match expected ones. Expected %v, got %v", e.ID, e, u)
		}
	}
}

func Test_CreateUser(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	newUser := models.NewUser("testuser", "testpassword")
	if err := a.UserRepo().Create(nil, newUser); err != nil {
		t.Errorf("could not create user %s", newUser.UserName)
	}

	userInDb := models.User{}
	result := a.DB().Where("username = ?", newUser.UserName).First(&userInDb)
	if err := result.Error; err != nil {
		t.Errorf("user %s is missing in db", newUser.UserName)
	}

	if match := compareUsers(*newUser, userInDb); !match {
		t.Errorf("properties for user %s do not match expected ones. Expected %v, got %v", newUser.UserName, newUser, userInDb)
	}
}

func Test_DeleteUsers(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	expectedUsersMap, err := testdb.SeedUsers(a.DB())
	if err != nil {
		t.Error("could not seed users")
	}

	if err = a.UserRepo().Delete(nil, "95a62e6c-e0e7-46ee-8bc3-6cca62b4cb09"); err != nil {
		t.Errorf("user %s is could not be deleted", "95a62e6c-e0e7-46ee-8bc3-6cca62b4cb09")
	}

	delete(expectedUsersMap, "95a62e6c-e0e7-46ee-8bc3-6cca62b4cb09")

	var users []models.User
	if result := a.DB().Find(&users); result.Error != nil {
		t.Error("could not get users list")
	}

	gotUsersMap := make(map[types.Uuid]models.User)
	for i := range users {
		gotUsersMap[users[i].ID] = users[i]
	}

	if len(expectedUsersMap) != len(gotUsersMap) {
		t.Error("expected users amount does not match received users amount")
	}

	for k, e := range expectedUsersMap {
		u, ok := gotUsersMap[k]
		if !ok {
			t.Errorf("user %s was not returned", e.ID)
		}

		if match := compareUsers(u, e); !match {
			t.Errorf("properties for user %s do not match expected ones. Expected %v, got %v", e.ID, e, u)
		}
	}
}

func compareUsers(user1, user2 models.User) bool {
	return user1.ID == user2.ID && user1.UserName == user2.UserName && user1.PasswordHash == user2.PasswordHash
}
