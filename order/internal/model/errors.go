package model

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderAlreadyPaid    = errors.New("order already paid")
	ErrOrderNotPaid        = errors.New("order not paid yet")
	ErrInvalidOrderStatus  = errors.New("invalid order status")
	ErrBadRequest          = errors.New("bad request")
	ErrPartsNotFound       = errors.New("some parts were not found")
	ErrPaymentFailed       = errors.New("payment failed")
	ErrInternalServerError = errors.New("internal server error")
)
