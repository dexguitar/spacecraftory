package user

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

func (s *UserService) GetUser(ctx context.Context, filter *model.UserFilter) (*model.User, error) {
	if filter == nil {
		return nil, model.ErrInvalidFilter
	}

	uuid := filter.UUID
	if uuid != nil {
		return s.userRepository.GetUserByUUID(ctx, *uuid)
	}

	login := filter.LoginData.Login
	password := filter.LoginData.Password

	if login != "" && password != "" {
		return s.userRepository.GetUserByLogin(ctx, login)
	}

	return nil, model.ErrUserNotFound
}
