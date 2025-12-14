package session

import (
	"context"
	"time"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/iam/internal/repository/converter"
	"github.com/dexguitar/spacecraftory/platform/pkg/cache"
)

func (r *repository) Set(ctx context.Context, session *model.Session, ttl time.Duration) error {
	cacheKey := r.getCacheKey(session.UUID)
	userSessionSetKey := r.getUserSessionSetKey(session.UserUUID)
	redisView := repoConverter.SessionToRedisView(session)

	return r.cache.TxPipeline(ctx, func(tx cache.TxPipeliner) error {
		if err := tx.HashSet(cacheKey, redisView); err != nil {
			return err
		}
		if err := tx.Expire(cacheKey, ttl); err != nil {
			return err
		}
		return tx.SAdd(userSessionSetKey, session.UUID)
	})
}
