package v1

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/service"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

type api struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) orderV1.Handler {
	return &api{
		orderService: orderService,
	}
}

func (a *api) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.GenericError{
			Code:    500,
			Message: err.Error(),
		},
	}
}
