package routes

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/controllers/admin"
	"github.com/dionofrizal88/dsi/auth/controllers/user"
	"github.com/dionofrizal88/dsi/auth/infrastructure/dao"
	libraryMiddleware "github.com/dionofrizal88/dsi/auth/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

type Router struct {
	config      config.Configuration
	dbService   *dao.Repositories
	redisClient *redis.Client
}

// NewRouter is a constructor will initialize Router.
func NewRouter(options ...RouterOption) *Router {
	router := &Router{}

	for _, opt := range options {
		opt(router)
	}

	return router
}

// Init is a function  to initialize Router.
func (r *Router) Init() *echo.Echo {
	e := echo.New()

	// test API
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello this is echo!")
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.Logger())

	adminController := admin.NewController(r.config, r.dbService, r.redisClient, e.Logger)
	userController := user.NewController(r.config, r.dbService, r.redisClient, e.Logger)

	// admin registration
	adminV1 := e.Group("/api/v1/internal/admin")
	adminV1.POST("/register", adminController.RegisterUserAdmin)

	adminV1.Use(libraryMiddleware.JWTMiddleware(r.config, r.redisClient))
	adminV1.GET("/recovery-request", adminController.GetRequestTokenByEmail)

	// auth API
	userV1 := e.Group("/api/v1/external/user")

	userV1.POST("/register", userController.RegisterUser)
	userV1.POST("/login", userController.Login)
	userV1.POST("/logout", userController.Logout)

	// recovery API
	userV1.POST("/recovery/request", userController.RequestRecovery)
	userV1.POST("/recovery/:token", userController.Recovery)

	return e
}
