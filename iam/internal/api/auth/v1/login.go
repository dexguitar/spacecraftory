package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	switch v := req.Type.GetBy().(type) {
	case *authV1.LoginType_ByLogin:
		login := v.ByLogin.GetLogin()
		password := v.ByLogin.GetPassword()

		sessionUUID, err := a.loginByLogin(ctx, &authV1.Login{
			Login:    login,
			Password: password,
		})
		if err != nil {
			return nil, err
		}
		return &authV1.LoginResponse{
			SessionUuid: sessionUUID,
		}, nil
	case *authV1.LoginType_ByUserUuid:
		userUUID := v.ByUserUuid.GetUserUuid()

		sessionUUID, err := a.loginByUserUuid(ctx, userUUID)
		if err != nil {
			return nil, err
		}
		return &authV1.LoginResponse{
			SessionUuid: sessionUUID,
		}, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid login type")
	}
}

func (a *api) loginByLogin(ctx context.Context, loginData *authV1.Login) (string, error) {
	if loginData == nil || loginData.GetLogin() == "" || loginData.GetPassword() == "" {
		return "", status.Errorf(codes.InvalidArgument, "login and password are required")
	}
	sessionUUID, err := a.authService.Login(ctx, &authV1.LoginType{
		By: &authV1.LoginType_ByLogin{
			ByLogin: loginData,
		},
	})
	if err != nil {
		return "", err
	}
	return sessionUUID, nil
}

func (a *api) loginByUserUuid(ctx context.Context, userUUID string) (string, error) {
	if userUUID == "" {
		return "", status.Errorf(codes.InvalidArgument, "user UUID is required")
	}
	sessionUUID, err := a.authService.Login(ctx, &authV1.LoginType{
		By: &authV1.LoginType_ByUserUuid{
			ByUserUuid: &authV1.UserUuid{
				UserUuid: userUUID,
			},
		},
	})
	if err != nil {
		return "", err
	}
	return sessionUUID, nil
}
