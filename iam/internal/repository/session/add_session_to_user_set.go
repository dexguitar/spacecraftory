package session

import (
	"context"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

func (r *repository) AddSessionToUserSet(ctx context.Context, userUUID string, session model.Session) error {
	userSessionSetKey := r.getUserSessionSetKey(userUUID)

	err := r.cache.SAdd(ctx, userSessionSetKey, session.UUID)
	if err != nil {
		return err
	}
	return nil
}
