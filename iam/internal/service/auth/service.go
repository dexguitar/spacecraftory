package auth

import (
	"time"

	"github.com/dexguitar/spacecraftory/iam/internal/repository"
	userService "github.com/dexguitar/spacecraftory/iam/internal/service"
)

type service struct {
	cacheRepository repository.CacheRepository
	cacheTTL        time.Duration
	userService     userService.UserService
}

func NewService(
	cacheRepository repository.CacheRepository,
	userService userService.UserService,
	cacheTTL time.Duration,
) *service {
	return &service{
		cacheRepository: cacheRepository,
		cacheTTL:        cacheTTL,
		userService:     userService,
	}
}
