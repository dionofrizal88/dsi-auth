package middleware

import (
	"context"
	"fmt"
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"net/http"
)

// JWTMiddleware is JWT middleware.
func JWTMiddleware(conf config.Configuration, rdb *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()
			token := c.Request().Header.Get("Authorization")
			jwt := security.NewJWT(conf)

			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
			}

			if !jwt.ValidateJWT(token) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			tokenDecode, err := jwt.DecodeJWT(token)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Invalid token"})
			}

			val, err := rdb.Get(ctx, fmt.Sprintf("user:%s:token", tokenDecode["user_id"])).Result()
			if err != nil || val == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired or invalid"})
			}

			return next(c)
		}
	}
}
