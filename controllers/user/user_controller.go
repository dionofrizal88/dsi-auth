package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

// RegisterUser is a function to save user data.
func (co *Controller) RegisterUser(c echo.Context) error {
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
	}, co.config.AppSecret+req.Email)
	if err != nil {
		co.logger.Errorf("error while create user, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, co.transformToResponse(ctx, "Success register user", nil))
}

// Login is a function to auth user data.
func (co *Controller) Login(c echo.Context) error {
	ctx := context.Background()
	var req Request

	err := c.Bind(&req)
	if err != nil {
		co.logger.Errorf("error binding request, err %v", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"message": "error binding request"})
	}

	result, errQuery := co.repository.User.FindUserByEmail(ctx, req.Email)
	if errQuery != nil {
		co.logger.Errorf("error while find user by email, err %v", errQuery)

		if errors.Is(errQuery, sql.ErrNoRows) {

			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Email or Password is incorrect"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
	}

	if result.IsRecovery {
		co.logger.Error("error user account under recovery")

		return c.JSON(http.StatusForbidden, map[string]string{"message": "Your account is currently under recovery. Please complete the recovery process before logging in"})
	}

	hashedPassword := security.HashPasswordWithSalt(req.Password, co.config.AppSecret+req.Email)

	if hashedPassword != result.Password {
		co.logger.Error("error invalid password")

		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Email or Password is incorrect"})
	}

	jwtGenerator := security.NewJWT(co.config)
	token, err := jwtGenerator.GenerateJWT(result.ID, req.Email)
	if err != nil {
		co.logger.Errorf("error while generate jwt token, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	// Store token in Redis
	err = co.redis.Set(ctx, fmt.Sprintf("user:%s:token", result.ID), token, 15*time.Minute).Err()
	if err != nil {
		co.logger.Errorf("error while save token into redis, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save token"})
	}

	return c.JSON(http.StatusOK, co.transformToAuthResponse(ctx, "Success authenticate user", token, result))
}

// Logout is a function to logout user.
func (co *Controller) Logout(c echo.Context) error {
	ctx := context.Background()
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

	val, err := co.redis.Get(ctx, fmt.Sprintf("user:%s:token", tokenDecode["user_id"])).Result()
	if err != nil || val == "" || val != token {
		co.logger.Errorf("error while get token from redis, err %v", err)

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired or invalid"})
	}

	err = co.redis.Del(ctx, fmt.Sprintf("user:%s:token", tokenDecode["user_id"])).Err()
	if err != nil {
		co.logger.Errorf("error while delete token on redis, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to log out"})
	}

	return c.JSON(http.StatusOK, Response{
		Message: "Success logout",
		Data:    nil,
		Meta:    nil,
	})
}

// RequestRecovery is a function to request recovery user account.
func (co *Controller) RequestRecovery(c echo.Context) error {
	ctx := context.Background()
	var req Request

	err := c.Bind(&req)
	if err != nil {
		co.logger.Errorf("error binding request, err %v", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"message": "error binding request"})
	}
	result, errQuery := co.repository.User.FindUserByEmail(ctx, req.Email)
	if errQuery != nil {
		co.logger.Errorf("error while find user by email, err %v", errQuery)

		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Failed request recovery"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
	}

	if !result.IsRecovery {
		target := &entity.User{
			ID: result.ID,
		}

		updateRequest := map[string]interface{}{
			"is_recovery": true,
		}

		errQuery = co.repository.User.UpdateUser(ctx, target, updateRequest)
		if errQuery != nil {
			co.logger.Errorf("error while update user data, err %v", errQuery)

			return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
		}
	}

	err = co.redis.Del(ctx, fmt.Sprintf("user:%s:token", result.ID)).Err()
	if err != nil {
		co.logger.Errorf("error while delete token from redis, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed request recovery"})
	}

	randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	token := security.GenerateRecoveryToken(fmt.Sprintf("%s.%d", result.Email, randomizer.Int()), co.config.AppSecret+req.Email)

	err = co.redis.Set(ctx, fmt.Sprintf("recovery-user:%s:token", result.ID), token, 15*time.Minute).Err()
	if err != nil {
		co.logger.Errorf("error while save recovery token into redis, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save token"})
	}

	// send access token into email

	return c.JSON(http.StatusAccepted, Response{
		Message: "Accepted request recovery",
		Data:    nil,
		Meta:    nil,
	})
}

// Recovery is a function to recovery user account.
func (co *Controller) Recovery(c echo.Context) error {
	ctx := context.Background()
	var req RecoveryRequest

	err := c.Bind(&req)
	if err != nil {
		co.logger.Errorf("error binding request, err %v", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"message": "error binding request"})
	}

	req.Token = c.Param("token")

	if req.Token == "" {
		co.logger.Errorf("error token is missing")

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
	}

	result, errQuery := co.repository.User.FindUserByEmail(ctx, req.Email)
	if errQuery != nil {
		co.logger.Errorf("error while find user by email, err %v", errQuery)

		if errors.Is(errQuery, sql.ErrNoRows) {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Failed request recovery"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
	}

	val, err := co.redis.Get(ctx, fmt.Sprintf("recovery-user:%s:token", result.ID)).Result()
	if err != nil || val == "" || val != req.Token {
		co.logger.Errorf("error while find recovery token from redis, err %v", err)

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired or invalid"})
	}

	target := &entity.User{
		ID: result.ID,
	}

	// Hash the password with the salt
	hashedPassword := security.HashPasswordWithSalt(req.NewPassword, co.config.AppSecret+req.Email)
	if result.Password == hashedPassword {
		co.logger.Error("error the new password is same as the old password")

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "New password cannot be the same as the old one. Please choose a different password"})
	}

	updateRequest := map[string]interface{}{
		"password":    hashedPassword,
		"is_recovery": false,
	}

	errQuery = co.repository.User.UpdateUser(ctx, target, updateRequest)
	if errQuery != nil {
		co.logger.Errorf("error while update user data, err %v", errQuery)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errQuery.Error()})
	}

	err = co.redis.Del(ctx, fmt.Sprintf(fmt.Sprintf("recovery-user:%s:token", result.ID))).Err()
	if err != nil {
		co.logger.Errorf("error while delete recovery token from redis, err %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed request recovery"})
	}

	return c.JSON(http.StatusOK, Response{
		Message: "Success recovery user account",
		Data:    nil,
		Meta:    nil,
	})
}
