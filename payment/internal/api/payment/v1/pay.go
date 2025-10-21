package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/payment/internal/converter"
	"github.com/dexguitar/spacecraftory/payment/internal/model"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	payment, err := converter.ToModelPayment(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request details")
	}
	transactionUUID, err := a.paymentService.PayOrder(ctx, payment)
	if err != nil {
		if errors.Is(err, model.ErrBadRequest) {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid request details")
		}
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
