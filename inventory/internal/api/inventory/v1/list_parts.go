package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/inventory/internal/converter"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := converter.ToModelPartsFilter(req.GetFilter())

	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	protoParts := converter.ToProtoParts(parts)

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
