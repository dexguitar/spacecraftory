package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/iam/internal/converter"
	"github.com/dexguitar/spacecraftory/iam/internal/model"
	userV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.GetUser(ctx, &model.UserFilter{
		UUID: &req.UserUuid,
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user")
	}

	return &userV1.GetUserResponse{
		User: converter.ToProtoUser(user),
	}, nil
}
