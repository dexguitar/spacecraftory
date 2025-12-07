package v1

import (
	"github.com/dexguitar/spacecraftory/iam/internal/service"
	userV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
