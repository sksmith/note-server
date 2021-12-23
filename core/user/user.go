package user

import (
	"context"
)

func NewService() Service {
	return Service{}
}

type Service struct {
}

func (u Service) Auth(ctx context.Context, username, password string) bool {
	return username == "test" && password == "test"
}
