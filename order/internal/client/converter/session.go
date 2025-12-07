package converter

import (
	"github.com/dexguitar/spacecraftory/order/internal/model"
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
		NotificationMethods: protoUserInfo.GetNotificationMethods(),
	}
}
