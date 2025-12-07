package v1

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/converter"
	userV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.GetUserByUUID(ctx, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &userV1.GetUserResponse{
		User: converter.ToProtoUser(user),
	}, nil
}
