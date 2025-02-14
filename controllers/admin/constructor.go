package admin

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/infrastructure/dao"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
)

// Controller struct is used when get configuration.
type Controller struct {
	config     config.Configuration
	repository *dao.Repositories
	redis      *redis.Client
	logger     echo.Logger
}

// NewController will initialize Controller.
func NewController(config config.Configuration, repository *dao.Repositories, redis *redis.Client, logger echo.Logger) *Controller {
	return &Controller{
		config:     config,
		repository: repository,
		redis:      redis,
		logger:     logger,
	}
}
