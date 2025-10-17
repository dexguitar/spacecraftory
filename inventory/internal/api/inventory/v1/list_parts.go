package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/inventory/internal/converter"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := converter.PartsFilterProtoToServiceModel(req.GetFilter())

	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list parts: %v", err)
	}


	protoParts := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, converter.PartServiceModelToProto(part))
	}

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
