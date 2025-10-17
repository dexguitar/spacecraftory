package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/inventory/internal/converter"
	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		return nil, status.Errorf(codes.Internal, "failed to get part: %v", err)
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartServiceModelToProto(part),
	}, nil
}
