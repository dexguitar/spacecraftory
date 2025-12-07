package model

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrSessionNotFound  = errors.New("session not found")
	ErrInvalidLoginData = errors.New("invalid login data")
)
