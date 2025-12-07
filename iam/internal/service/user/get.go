package user

import (
	"context"
	"errors"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

func (s *UserService) GetUserByUUID(ctx context.Context, userUUID string) (*model.User, error) {
	user, err := s.userRepository.GetUserByUUID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrInvalidLoginData
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByLogin(ctx context.Context, login, password string) (*model.User, error) {
	user, err := s.userRepository.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	if user.Password != password {
		return nil, model.ErrInvalidLoginData
	}
	return user, nil
}
