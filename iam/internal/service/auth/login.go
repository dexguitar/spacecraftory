package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	filter := &model.UserFilter{
		LoginData: &model.LoginData{
			Login:    login,
			Password: password,
		},
	}
	user, err := s.userService.GetUser(ctx, filter)
	if err != nil {
		return "", err
	}

	session := &model.Session{
		UUID:      uuid.New().String(),
		UserUUID:  user.UUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.cacheTTL),
	}

	if err = s.cacheRepository.Set(ctx, session, s.cacheTTL); err != nil {
		return "", err
	}

	return session.UUID, nil
}
