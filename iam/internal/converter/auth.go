package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	commonV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/common/v1"
)

func ToModelSession(session *commonV1.Session) *model.Session {
	return &model.Session{
		UUID:      session.GetUuid(),
		CreatedAt: session.GetCreatedAt().AsTime(),
		UpdatedAt: session.GetUpdatedAt().AsTime(),
		ExpiresAt: session.GetExpiresAt().AsTime(),
	}
}

func ToProtoSession(session *model.Session) *commonV1.Session {
	return &commonV1.Session{
		Uuid:      session.UUID,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
