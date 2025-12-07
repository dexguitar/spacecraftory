package service

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

type UserService interface {
	Register(ctx context.Context, user *model.User) (string, error)
	GetUserByUUID(ctx context.Context, userUUID string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login, password string) (*model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, loginType *authV1.LoginType) (string, error)
	WhoAmI(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error)
}
