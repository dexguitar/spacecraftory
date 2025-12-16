package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	authAPI "github.com/dexguitar/spacecraftory/iam/internal/api/auth/v1"
	userApi "github.com/dexguitar/spacecraftory/iam/internal/api/user/v1"
	"github.com/dexguitar/spacecraftory/iam/internal/config"
	"github.com/dexguitar/spacecraftory/iam/internal/repository"
	cacheRepository "github.com/dexguitar/spacecraftory/iam/internal/repository/session"
	userRepository "github.com/dexguitar/spacecraftory/iam/internal/repository/user"
	"github.com/dexguitar/spacecraftory/iam/internal/service"
	authService "github.com/dexguitar/spacecraftory/iam/internal/service/auth"
	userService "github.com/dexguitar/spacecraftory/iam/internal/service/user"
	"github.com/dexguitar/spacecraftory/platform/pkg/cache"
	"github.com/dexguitar/spacecraftory/platform/pkg/cache/redis"
	"github.com/dexguitar/spacecraftory/platform/pkg/closer"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
	userV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authApi         authV1.AuthServiceServer
	authService     service.AuthService
	cacheRepository repository.CacheRepository
	redisPool       *redigo.Pool
	redisClient     cache.RedisClient

	userApi        userV1.UserServiceServer
	userService    service.UserService
	userRepository repository.UserRepository
	pgPool         *pgxpool.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authApi == nil {
		d.authApi = authAPI.NewAPI(d.AuthService(ctx))
	}

	return d.authApi
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(d.CacheRepository(ctx), d.UserService(ctx), config.AppConfig().Redis.CacheTTL())
	}

	return d.authService
}

func (d *diContainer) CacheRepository(ctx context.Context) repository.CacheRepository {
	if d.cacheRepository == nil {
		d.cacheRepository = cacheRepository.NewRepository(d.RedisClient(ctx))
	}

	return d.cacheRepository
}

func (d *diContainer) RedisPool(ctx context.Context) *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(ctx), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userApi == nil {
		d.userApi = userApi.NewAPI(d.UserService(ctx))
	}

	return d.userApi
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewUserService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewUserRepository(d.PgPool(ctx))
	}

	return d.userRepository
}

func (d *diContainer) PgPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		dbURI := config.AppConfig().Postgres.Address()

		pool, err := pgxpool.New(ctx, dbURI)
		if err != nil {
			panic(fmt.Sprintf("failed to create connection pool: %s", err.Error()))
		}

		closer.AddNamed("PostgreSQL connection pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}
