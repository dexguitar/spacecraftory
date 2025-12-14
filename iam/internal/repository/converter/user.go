package converter

import (
	"github.com/dexguitar/spacecraftory/iam/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/iam/internal/repository/model"
)

func ToModelUserFromRow(row *repoModel.UserRow, methods []repoModel.NotificationMethodRow) *model.User {
	return &model.User{
		UUID:     row.ID,
		Password: row.Password,
		Info: model.UserInfo{
			Login:               row.Login,
			Email:               row.Email,
			NotificationMethods: toModelNotificationMethods(methods),
		},
	}
}

func toModelNotificationMethods(methods []repoModel.NotificationMethodRow) []model.NotificationMethod {
	result := make([]model.NotificationMethod, 0, len(methods))
	for _, m := range methods {
		result = append(result, model.NotificationMethod{
			ProviderName: m.ProviderName,
			Target:       m.Target,
		})
	}
	return result
}
