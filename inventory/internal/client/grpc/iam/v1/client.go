package iam

import (
	"context"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

// IAMClient interface for IAM service operations
type IAMClient interface {
	WhoAmI(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error)
}

type iamClient struct {
	grpcClient authV1.AuthServiceClient
}

func NewIAMClient(grpcClient authV1.AuthServiceClient) IAMClient {
	return &iamClient{
		grpcClient: grpcClient,
	}
}
