package auth

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

func (s *service) WhoAmI(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error) {
	session, err := s.cacheRepository.Get(ctx, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userService.GetUserByUUID(ctx, session.UserUUID)
	if err != nil {
		return nil, nil, err
	}

	return &session, user, nil
}
