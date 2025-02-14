package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

// RegisterUserAdmin is a function to save user data.
func (co *Controller) RegisterUserAdmin(c echo.Context) error {
	ctx := c.Request().Context()
	var req Request

	err := c.Bind(&req)
	if err != nil {
		co.logger.Errorf("error binding request, err %v", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"message": "error binding request"})
	}

	result, errQuery := co.repository.User.FindUserByEmail(ctx, req.Email)
	if errQuery != nil {
		if !errors.Is(errQuery, gorm.ErrRecordNotFound) {
			co.logger.Errorf("error while find user by email, err %v", errQuery)

			return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
		}
	}

	if result != nil {
		co.logger.Error("error email already exists")

		return c.JSON(http.StatusConflict, map[string]string{"message": "User email already exists"})
	}

	_, err = co.repository.User.CreateUser(ctx, &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		IsAdmin:  true,
	}, co.config.AppSecret+req.Email)
	if err != nil {
		co.logger.Errorf("error while create user, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, co.transformToResponse(ctx, "Success register user admin", nil))
}

// GetRequestTokenByEmail is a function to find request token for admin.
func (co *Controller) GetRequestTokenByEmail(c echo.Context) error {
	ctx := context.Background()
	var req Request

	err := c.Bind(&req)
	if err != nil {
		co.logger.Errorf("error binding request, err %v", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"message": "error binding request"})
	}

	token := c.Request().Header.Get("Authorization")

	if token == "" {
		co.logger.Error("error token is missing")

		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing token"})
	}

	jwtGenerator := security.NewJWT(co.config)
	tokenDecode, err := jwtGenerator.DecodeJWT(token)
	if err != nil {
		co.logger.Errorf("error while decode jwt token, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Invalid token"})
	}

	result, errQuery := co.repository.User.FindUser(ctx, tokenDecode["user_id"].(string))
	if errQuery != nil {
		co.logger.Errorf("error while find user by email, err %v", errQuery)

		if errors.Is(errQuery, sql.ErrNoRows) {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Failed request token"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
	}

	if !result.IsAdmin {
		co.logger.Errorf("error user account is not an admin")

		return c.JSON(http.StatusForbidden, map[string]string{"message": "Forbidden access"})
	}

	result, errQuery = co.repository.User.FindUserByEmail(ctx, req.Email)
	if errQuery != nil {
		co.logger.Errorf("error while find user by email, err %v", errQuery)

		if errors.Is(errQuery, sql.ErrNoRows) {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Failed request recovery"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
	}

	if !result.IsRecovery {
		co.logger.Errorf("error user account is doesn't under recovery")

		return c.JSON(http.StatusBadRequest, map[string]string{"message": "user account is doesn't under recovery"})
	}

	val, err := co.redis.Get(ctx, fmt.Sprintf("recovery-user:%s:token", result.ID)).Result()
	if err != nil || val == "" {
		co.logger.Errorf("error while find recovery token from redis, err %v", err)

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired or invalid"})
	}

	ttl, err := co.redis.TTL(ctx, fmt.Sprintf("recovery-user:%s:token", result.ID)).Result()
	if err != nil {
		co.logger.Errorf("error while getting TTL for recovery token, err %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve token expiry time"})
	}

	minutes := int(ttl.Seconds()) / 60
	seconds := int(ttl.Seconds()) % 60
	validTokenTime := fmt.Sprintf("%d minutes and %d seconds", minutes, seconds)

	var recoveryLink string

	switch co.config.AppEnv {
	case "production":
		recoveryLink = fmt.Sprintf("production:8081/api/v1/external/user/recovery/%s", req.Email)

	case "staging":
		recoveryLink = fmt.Sprintf("staging:8081/api/v1/external/user/recovery/%s", req.Email)

	case "development":
		recoveryLink = fmt.Sprintf("development:8081/api/v1/external/user/recovery/%s", req.Email)

	default:
		recoveryLink = fmt.Sprintf("localhost:8081/api/v1/external/user/recovery/%s", val)

	}

	return c.JSON(http.StatusOK, Response{
		Message: "Success request recovery",
		Data: map[string]string{
			"send_to":       result.Email,
			"recovery_link": recoveryLink,
			"subject":       "Request Reset Request - Digital Sekuriti Indonesia",
			"message":       fmt.Sprintf(RecoveryRequestEmailTemplate, req.Name, validTokenTime, recoveryLink),
		},
		Meta: nil,
	})
}
