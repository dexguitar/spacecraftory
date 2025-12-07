package iam

import (
	"context"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
)

type IAMClient interface {
	WhoAmI(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error)
}
