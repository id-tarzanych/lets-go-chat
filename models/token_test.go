package models

import (
	"github.com/id-tarzanych/lets-go-chat/internal/types"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	now := time.Now()

	type args struct {
		token      string
		userId     types.Uuid
		expiration time.Time
	}

	tests := []struct {
		name string
		args args
		want *Token
	}{
		{"Token 1",
			args{
				"glf1LdUMtwLssv48",
				"ee80103a-7f8e-4d45-8613-452f5c695c5a",
				now.Add(time.Hour * 24),
			},
			&Token{
				Token:      "glf1LdUMtwLssv48",
				UserId:     "ee80103a-7f8e-4d45-8613-452f5c695c5a",
				Expiration: now.Add(time.Hour * 24),
			},
		},
		{
			"Token 2",
			args{
				"6C7hqgNLkbkRVLlU",
				"e38ac9d5-6d6d-4e4c-803d-ca1869feccdb",
				now.Add(time.Hour * 96)},
			&Token{
				Token:      "6C7hqgNLkbkRVLlU",
				UserId:     "e38ac9d5-6d6d-4e4c-803d-ca1869feccdb",
				Expiration: now.Add(time.Hour * 96),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewToken(tt.args.token, tt.args.userId, tt.args.expiration)

			if  tt.want.UserId != got.UserId || tt.want.Token != got.Token || tt.want.Expiration != got.Expiration {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
