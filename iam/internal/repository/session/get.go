package session

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/iam/internal/repository/converter"
	repoModel "github.com/dexguitar/spacecraftory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Session, error) {
	cacheKey := r.getCacheKey(uuid)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrSessionNotFound
	}

	var sessionRedisView repoModel.SessionRedisView
	err = redigo.ScanStruct(values, &sessionRedisView)
	if err != nil {
		return nil, err
	}

	return repoConverter.SessionFromRedisView(&sessionRedisView), nil
}
