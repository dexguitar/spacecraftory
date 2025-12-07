package session

import (
	"context"
	"time"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/iam/internal/repository/converter"
)

func (r *repository) Set(ctx context.Context, uuid string, session model.Session, ttl time.Duration) error {
	cacheKey := r.getCacheKey(uuid)

	redisView := repoConverter.SessionToRedisView(&session)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}
