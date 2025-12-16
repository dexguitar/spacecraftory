package service

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

type UserService interface {
	Register(ctx context.Context, user *model.User) (string, error)
	GetUser(ctx context.Context, filter *model.UserFilter) (*model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	WhoAmI(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error)
}
