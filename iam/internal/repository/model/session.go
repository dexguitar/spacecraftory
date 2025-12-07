package model

import (
	"time"
)

type Session struct {
	UUID      string
	UserUUID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type SessionRedisView struct {
	UUID      string `redis:"uuid"`
	UserUUID  string `redis:"user_uuid"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt int64  `redis:"updated_at"`
	ExpiresAt int64  `redis:"expires_at"`
}
