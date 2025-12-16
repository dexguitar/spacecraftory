package converter

import (
	"time"

	serviceModel "github.com/dexguitar/spacecraftory/iam/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/iam/internal/repository/model"
)

func ToRepoSession(session *serviceModel.Session) *repoModel.Session {
	return &repoModel.Session{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}

func ToModelSession(session *repoModel.Session) *serviceModel.Session {
	return &serviceModel.Session{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}

func SessionFromRedisView(session *repoModel.SessionRedisView) *serviceModel.Session {
	return &serviceModel.Session{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: time.Unix(session.CreatedAt, 0),
		UpdatedAt: time.Unix(session.UpdatedAt, 0),
		ExpiresAt: time.Unix(session.ExpiresAt, 0),
	}
}

func SessionToRedisView(session *serviceModel.Session) *repoModel.SessionRedisView {
	return &repoModel.SessionRedisView{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt.Unix(),
		UpdatedAt: session.UpdatedAt.Unix(),
		ExpiresAt: session.ExpiresAt.Unix(),
	}
}
