package inventory

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/client/converter"
	"github.com/dexguitar/spacecraftory/order/internal/model"
	authGrpc "github.com/dexguitar/spacecraftory/platform/pkg/middleware/grpc"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	req := &inventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	}

	ctx = authGrpc.ForwardSessionUUIDToGRPC(ctx)

	resp, err := c.grpcClient.ListParts(ctx, req)
	if err != nil {
		return nil, err
	}

	parts := make([]model.Part, 0, len(resp.Parts))
	for _, protoPart := range resp.Parts {
		parts = append(parts, converter.PartProtoToServiceModel(protoPart))
	}

	return parts, nil
}
