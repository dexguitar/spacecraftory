package model

import "errors"

var (
	ErrPartNotFound = errors.New("part not found")
	ErrBadRequest   = errors.New("bad request")
)
