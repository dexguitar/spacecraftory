package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	login := req.GetLogin()
	password := req.GetPassword()

	if login == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "login and password are required, login: %s, password: %s", login, password)
	}

	sessionUUID, err := a.authService.Login(ctx, login, password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidLoginData) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid login data")
		}
		return nil, status.Errorf(codes.Internal, "failed to login")
	}

	return &authV1.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
