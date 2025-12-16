package converter

import (
	"github.com/dexguitar/spacecraftory/iam/internal/model"
	commonV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/common/v1"
	userV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/user/v1"
)

func ToModelUser(user *userV1.UserRegistrationInfo) *model.User {
	return &model.User{
		Info:     *ToModelUserInfo(user.GetInfo()),
		Password: user.Password,
	}
}

func ToModelUserInfo(userInfo *commonV1.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Login:               userInfo.GetLogin(),
		Email:               userInfo.GetEmail(),
		NotificationMethods: toModelNotificationMethods(userInfo.GetNotificationMethods()),
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

func ToProtoUser(user *model.User) *commonV1.User {
	return &commonV1.User{
		Uuid: user.UUID,
		Info: ToProtoUserInfo(&user.Info),
	}
}

func ToProtoUserInfo(userInfo *model.UserInfo) *commonV1.UserInfo {
	return &commonV1.UserInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		NotificationMethods: toProtoNotificationMethods(userInfo.NotificationMethods),
	}
}

func toProtoNotificationMethods(methods []model.NotificationMethod) []*commonV1.NotificationMethod {
	result := make([]*commonV1.NotificationMethod, 0, len(methods))
	for _, m := range methods {
		result = append(result, &commonV1.NotificationMethod{
			ProviderName: m.ProviderName,
			Target:       m.Target,
		})
	}
	return result
}
