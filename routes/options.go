package routes

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/infrastructure/dao"
	"github.com/go-redis/redis/v8"
)

// RouterOption return Router with RouterOption.
type RouterOption func(*Router)

// WithConfig is a function to set config to the RouterOption.
func WithConfig(config config.Configuration) RouterOption {
	return func(r *Router) {
		r.config = config
	}
}

// WithDBService is a function to set dbService to the RouterOption.
func WithDBService(dbService *dao.Repositories) RouterOption {
	return func(r *Router) {
		r.dbService = dbService
	}
}

// WithRedisDB is a function to set redis db to the RouterOption.
func WithRedisDB(rdb *redis.Client) RouterOption {
	return func(r *Router) {
		r.redisClient = rdb
	}
}
