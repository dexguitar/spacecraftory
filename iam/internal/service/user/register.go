package user

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

func (s *UserService) Register(ctx context.Context, user *model.User) (string, error) {
	userUUID, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return userUUID, nil
}
