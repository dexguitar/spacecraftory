package session

import (
	"fmt"

	"github.com/dexguitar/spacecraftory/platform/pkg/cache"
)

const (
	cacheKeyPrefix       = "iam:session:"
	userSessionSetPrefix = "iam:user:session:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

// api -> service -> cache repo -> redis client (обертка наша) -> redis
func (r *repository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}

func (r *repository) getUserSessionSetKey(userUUID string) string {
	return fmt.Sprintf("%s%s", userSessionSetPrefix, userUUID)
}
