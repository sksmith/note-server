package user_test

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sksmith/note-server/core/user"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}

func TestAuth(t *testing.T) {
	tests := []struct {
		ctx      context.Context
		username string
		password string
		want     bool
	}{
		{ctx: context.Background(), username: "test", password: "test", want: true},
		{ctx: context.Background(), username: "test", password: "badpass", want: false},
		{ctx: context.Background(), username: "baduser", password: "test", want: false},
		{ctx: context.Background(), username: "baduser", password: "badpass", want: false},
	}

	for _, test := range tests {
		svc := user.NewService()

		if got := svc.Auth(test.ctx, test.username, test.password); got != test.want {
			t.Errorf("got=[%v] want=[%v]", got, test.want)
		}
	}
}
