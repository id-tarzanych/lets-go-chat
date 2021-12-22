package models

import (
	"fmt"
	"testing"
)

func TestNewUser(t *testing.T) {
	type pair struct {
		username string
		password string
	}

	tests := []struct {
		name string
		args pair
		want pair
	}{
		{
			name: "User 1",
			args: pair{
				username: "user1",
				password: "12345678",
			},
			want: pair{
				username: "user1",
				password: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f",
			},
		},
		{"User 2", pair{"user2", "87654321"}, pair{"user2", "e24df920078c3dd4e7e8d2442f00e5c9ab2a231bb3918d65cc50906e49ecaef4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.args.username, tt.args.password)

			if got.UserName != tt.want.username || got.PasswordHash != tt.want.password {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_SetPassword(t *testing.T) {
	user := &User{UserName: "user", PasswordHash: "p@ssw0rd"}

	tests := []struct {
		name     string
		password string
		wantHash string
	}{
		{
			name:     "Password 1",
			password: "12345678",
			wantHash: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f",
		},
		{
			name:     "Password 2",
			password: "Secr3tP@s$w0rd",
			wantHash: "39a982ef00b8dff3632db20d56141060648ad3bf0dbc57d3e57386aa6b3b81d1",
		},
		{
			name:     "Password 3",
			password: "SomePass",
			wantHash: "c21c4ec6e85ba12ccb2ccf117dd405a9f5915f2722d50dce4739400b164cfb9e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user.SetPassword(tt.password)

			if got := user.PasswordHash; got != tt.wantHash {
				t.Errorf("Password = %v, want %v", got, tt.wantHash)
			}
		})
	}
}

func ExampleNewUser() {
	user := NewUser("testuser", "12345678")

	fmt.Println(user.UserName)
	fmt.Println(user.PasswordHash)

	// Output:
	// testuser
	// ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f
}

func ExampleUser_SetPassword() {
	user := NewUser("testuser", "12345678")
	user.SetPassword("87654321")

	fmt.Println(user.PasswordHash)

	// Output:
	// e24df920078c3dd4e7e8d2442f00e5c9ab2a231bb3918d65cc50906e49ecaef4
}

func BenchmarkNewUser(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewUser("testuser", "testpassword")
	}
}
