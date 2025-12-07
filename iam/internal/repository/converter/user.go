package converter

import (
	"github.com/dexguitar/spacecraftory/iam/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/iam/internal/repository/model"
)

func ToRepoUser(user *model.User) *repoModel.User {
	return &repoModel.User{
		UUID:     user.UUID,
		Password: user.Password,
		Info:     *ToRepoUserInfo(&user.Info),
	}
}

func ToModelUser(user *repoModel.User) *model.User {
	return &model.User{
		UUID:     user.UUID,
		Info:     *ToModelUserInfo(&user.Info),
		Password: user.Password,
	}
}

// ToModelUserFromRow converts a flat UserRow (from DB) to the nested domain User
func ToModelUserFromRow(row *repoModel.UserRow) *model.User {
	return &model.User{
		UUID:     row.ID,
		Password: row.Password,
		Info: model.UserInfo{
			Login:               row.Login,
			Email:               row.Email,
			NotificationMethods: row.NotificationMethods,
		},
	}
}

func ToModelUserInfo(userInfo *repoModel.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		NotificationMethods: userInfo.NotificationMethods,
	}
}

func ToRepoUserInfo(userInfo *model.UserInfo) *repoModel.UserInfo {
	return &repoModel.UserInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		NotificationMethods: userInfo.NotificationMethods,
	}
}
