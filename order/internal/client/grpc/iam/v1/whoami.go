package iam

import (
	"context"

	converter "github.com/dexguitar/spacecraftory/order/internal/client/converter"
	"github.com/dexguitar/spacecraftory/order/internal/model"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

func (c *iamClient) WhoAmI(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error) {
	req := &authV1.WhoAmIRequest{
		SessionUuid: sessionUUID,
	}

	resp, err := c.grpcClient.WhoAmI(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	session := converter.SessionProtoToServiceModel(resp.Session)
	user := converter.UserProtoToServiceModel(resp.User)

	return session, user, nil
}
