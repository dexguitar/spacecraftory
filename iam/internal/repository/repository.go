package repository

import (
	"context"
	"time"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUserByUUID(ctx context.Context, userUUID string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

type CacheRepository interface {
	Get(ctx context.Context, uuid string) (model.Session, error)
	Set(ctx context.Context, uuid string, session model.Session, ttl time.Duration) error
	AddSessionToUserSet(ctx context.Context, userUUID string, session model.Session) error
	Delete(ctx context.Context, uuid string) error
}
