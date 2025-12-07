package v1

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/converter"
	userV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	user := converter.ToModelUser(req.Info)
	userUUID, err := a.userService.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return &userV1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
