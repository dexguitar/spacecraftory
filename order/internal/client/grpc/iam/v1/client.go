package iam

import (
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
)

type iamClient struct {
	grpcClient authV1.AuthServiceClient
}

func NewIAMClient(grpcClient authV1.AuthServiceClient) *iamClient {
	return &iamClient{
		grpcClient: grpcClient,
	}
}
