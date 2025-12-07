package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

func (s *service) Login(ctx context.Context, loginType *authV1.LoginType) (string, error) {
	if loginType == nil || loginType.GetBy() == nil {
		return "", model.ErrInvalidLoginData
	}

	var user *model.User
	var err error

	switch v := loginType.GetBy().(type) {
	case *authV1.LoginType_ByLogin:
		user, err = s.userService.GetUserByLogin(ctx, v.ByLogin.Login, v.ByLogin.Password)
		if err != nil {
			return "", err
		}
	case *authV1.LoginType_ByUserUuid:
		user, err = s.userService.GetUserByUUID(ctx, v.ByUserUuid.UserUuid)
		if err != nil {
			return "", err
		}
	default:
		return "", model.ErrInvalidLoginData
	}

	if user == nil {
		return "", model.ErrInvalidLoginData
	}

	session := model.Session{
		UUID:      uuid.New().String(),
		UserUUID:  user.UUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.cacheTTL),
	}

	err = s.cacheRepository.Set(ctx, session.UUID, session, s.cacheTTL)
	if err != nil {
		return "", err
	}
	err = s.cacheRepository.AddSessionToUserSet(ctx, user.UUID, session)
	if err != nil {
		return "", err
	}

	return session.UUID, nil
}
