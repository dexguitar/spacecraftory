package converter

import (
	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	commonV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/common/v1"
)

func SessionProtoToServiceModel(protoSession *commonV1.Session) *model.Session {
	return &model.Session{
		UUID:      protoSession.Uuid,
		CreatedAt: protoSession.CreatedAt.AsTime(),
		UpdatedAt: protoSession.UpdatedAt.AsTime(),
		ExpiresAt: protoSession.ExpiresAt.AsTime(),
	}
}

func UserProtoToServiceModel(protoUser *commonV1.User) *model.User {
	return &model.User{
		UUID: protoUser.GetUuid(),
		Info: *UserInfoProtoToServiceModel(protoUser.GetInfo()),
	}
}

func UserInfoProtoToServiceModel(protoUserInfo *commonV1.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Login:               protoUserInfo.GetLogin(),
		Email:               protoUserInfo.GetEmail(),
		NotificationMethods: toModelNotificationMethods(protoUserInfo.GetNotificationMethods()),
	}
}

func toModelNotificationMethods(methods []*commonV1.NotificationMethod) []model.NotificationMethod {
	result := make([]model.NotificationMethod, 0, len(methods))
	for _, m := range methods {
		result = append(result, model.NotificationMethod{
			ProviderName: m.GetProviderName(),
			Target:       m.GetTarget(),
		})
	}
	return result
}
