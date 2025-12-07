package v1

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/converter"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

func (a *api) WhoAmI(ctx context.Context, req *authV1.WhoAmIRequest) (*authV1.WhoAmIResponse, error) {
	session, user, err := a.authService.WhoAmI(ctx, req.SessionUuid)
	if err != nil {
		return nil, err
	}

	return &authV1.WhoAmIResponse{
		Session: converter.ToProtoSession(session),
		User:    converter.ToProtoUser(user),
	}, nil
}
