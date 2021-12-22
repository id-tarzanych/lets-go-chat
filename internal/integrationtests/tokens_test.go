package integrationtests

import (
	"testing"
	"time"

	"github.com/id-tarzanych/lets-go-chat/internal/testdb"
	"github.com/id-tarzanych/lets-go-chat/models"
)

func Test_GetAllTokens(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	expectedTokensMap, err := testdb.SeedTokens(a.DB())
	if err != nil {
		t.Error("could not seed tokens")
	}

	tokens, err := a.TokenRepo().GetAll(nil)
	if err != nil {
		t.Error("could not load tokens from database")
	}

	gotTokensMap := make(map[string]models.Token)
	for i := range tokens {
		gotTokensMap[tokens[i].Token] = tokens[i]
	}

	if len(expectedTokensMap) != len(gotTokensMap) {
		t.Error("expected token amount does not match received tokens amount")
	}

	for k, e := range expectedTokensMap {
		tk, ok := gotTokensMap[k]
		if !ok {
			t.Errorf("token %s was not returned", e.Token)
		}

		if match := compareTokens(tk, e); !match {
			t.Errorf("properties for token %s do not match expected ones. Expected %v, got %v", e.Token, e, tk)
		}
	}
}

func Test_GetToken(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	expectedTokenMap, err := testdb.SeedTokens(a.DB())
	if err != nil {
		t.Error("could not seed tokens")
	}

	for _, e := range expectedTokenMap {
		tk, err := a.TokenRepo().Get(nil, e.Token)
		if err != nil {
			t.Errorf("token %s was not returned", e.Token)
		}

		if match := compareTokens(tk, e); !match {
			t.Errorf("properties for token %s do not match expected ones. Expected %v, got %v", tk.Token, e, tk)
		}
	}
}

func Test_GetTokenByUserId(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	_, err := testdb.SeedTokens(a.DB())
	if err != nil {
		t.Error("could not seed tokens")
	}

	tokens, err := a.TokenRepo().GetByUserId(nil, "6b2db94c-6fce-4673-a1ce-d24ff6bd4d35")
	if err != nil {
		t.Error("tokens were not returned")
	}

	gotTokensMap := make(map[string]models.Token)
	for i := range tokens {
		gotTokensMap[tokens[i].Token] = tokens[i]
	}

	expectedTokensMap := map[string]models.Token{
		"18sqhpLyANr7ypoK": {
			Token:      "18sqhpLyANr7ypoK",
			UserId:     "6b2db94c-6fce-4673-a1ce-d24ff6bd4d35",
			Expiration: time.Now().Add(time.Hour * 24),
		},
		"8n9hKwlT9l037PZb": {
			Token:      "8n9hKwlT9l037PZb",
			UserId:     "6b2db94c-6fce-4673-a1ce-d24ff6bd4d35",
			Expiration: time.Now().Add(time.Hour * 24),
		},
	}

	for k, e := range expectedTokensMap {
		tk, ok := gotTokensMap[k]
		if !ok {
			t.Errorf("token %s was not returned", e.Token)
		}

		if match := compareTokens(tk, e); !match {
			t.Errorf("properties for token %s do not match expected ones. Expected %v, got %v", e.Token, e, tk)
		}
	}

}

func Test_CreateToken(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	newToken := models.NewToken("testtoken", "6b2db94c-6fce-4673-a1ce-d24ff6bd4d35", time.Now().Add(time.Hour*48))
	if err := a.TokenRepo().Create(nil, newToken); err != nil {
		t.Errorf("could not create token %s", newToken.Token)
	}

	tokenInDb := models.Token{}
	result := a.DB().Where("token = ?", newToken.Token).First(&tokenInDb)
	if err := result.Error; err != nil {
		t.Errorf("token %s is missing in db", newToken.Token)
	}

	if match := compareTokens(*newToken, tokenInDb); !match {
		t.Errorf("properties for user %s do not match expected ones. Expected %v, got %v", newToken.Token, newToken, tokenInDb)
	}
}

func Test_DeleteToken(t *testing.T) {
	defer func() {
		if err := testdb.Truncate(a.DB()); err != nil {
			t.Error("error truncating test database tables")
		}
	}()

	expectedTokensMap, err := testdb.SeedTokens(a.DB())
	if err != nil {
		t.Error("could not seed tokens")
	}

	if err = a.TokenRepo().Delete(nil, "8n9hKwlT9l037PZb"); err != nil {
		t.Errorf("token %s is could not be deleted", "8n9hKwlT9l037PZb")
	}

	delete(expectedTokensMap, "8n9hKwlT9l037PZb")

	var tokens []models.Token
	if result := a.DB().Find(&tokens); result.Error != nil {
		t.Error("could not get tokens list")
	}

	gotTokensMap := make(map[string]models.Token)
	for i := range tokens {
		gotTokensMap[tokens[i].Token] = tokens[i]
	}

	if len(expectedTokensMap) != len(gotTokensMap) {
		t.Error("expected tokens amount does not match received tokens amount")
	}

	for k, e := range expectedTokensMap {
		tk, ok := gotTokensMap[k]
		if !ok {
			t.Errorf("token %s was not returned", e.Token)
		}

		if match := compareTokens(tk, e); !match {
			t.Errorf("properties for token %s do not match expected ones. Expected %v, got %v", e.Token, e, tk)
		}
	}
}

func compareTokens(token1, token2 models.Token) bool {
	return token1.Token == token2.Token && token1.UserId == token2.UserId && token1.Expiration.Unix() == token2.Expiration.Unix()
}
