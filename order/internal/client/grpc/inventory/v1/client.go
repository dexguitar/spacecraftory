package inventory

import (
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

type inventoryClient struct {
	grpcClient inventoryV1.InventoryServiceClient
}

func NewInventoryClient(grpcClient inventoryV1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{
		grpcClient: grpcClient,
	}
}
