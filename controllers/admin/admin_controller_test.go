package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dionofrizal88/dsi/auth/controllers/admin"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/dionofrizal88/dsi/auth/pkg/tests"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestControllerAdmin(t *testing.T) {
	testSuite := tests.InitTestSuite()

	t.Cleanup(func() {
		testSuite.Clean()
	})

	reqBody := admin.Request{
		Name:     "Test Admin 1",
		Email:    "testAdmin1@gmail.com",
		Password: "testAdmin1234",
	}

	jsonData, _ := json.Marshal(reqBody)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminController := admin.NewController(testSuite.Config, testSuite.DBService, testSuite.RedisClient, e.Logger)
	errRegisterUserAdmin := adminController.RegisterUserAdmin(c)

	t.Run("positive case to test admin controller, expected no error", func(t *testing.T) {
		t.Run("positive case while send register admin request, expected no error", func(t *testing.T) {

			assert.NoError(t, errRegisterUserAdmin)
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Contains(t, rec.Body.String(), "Success")
		})
	})

	t.Run("negative case to test admin controller, expected no error", func(t *testing.T) {
		t.Run("negative case while email already exist, expected no error", func(t *testing.T) {
			req1 := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
			req1.Header.Set("Content-Type", "application/json")
			rec1 := httptest.NewRecorder()
			c1 := e.NewContext(req1, rec1)
			errRegisterUserAdmin1 := adminController.RegisterUserAdmin(c1)

			assert.NoError(t, errRegisterUserAdmin1)
			assert.Equal(t, http.StatusConflict, rec1.Code)
			assert.Contains(t, rec1.Body.String(), "already exists")
		})
	})

	t.Run("positive case to test admin controller, expected no error", func(t *testing.T) {
		t.Run("positive case while send get request token by email, expected no error", func(t *testing.T) {
			ctx := context.Background()
			result, err := testSuite.DBService.User.FindUserByEmail(ctx, reqBody.Email)

			require.NoError(t, err)

			target := &entity.User{
				ID: result.ID,
			}

			updateRequest := map[string]interface{}{
				"is_recovery": true,
			}

			err = testSuite.DBService.User.UpdateUser(ctx, target, updateRequest)

			require.NoError(t, err)

			jwtGenerator := security.NewJWT(testSuite.Config)
			token, err := jwtGenerator.GenerateJWT(result.ID, result.Email)

			require.NoError(t, err)

			err = testSuite.RedisClient.Set(ctx, fmt.Sprintf("recovery-user:%s:token", result.ID), token, 1*time.Minute).Err()

			require.NoError(t, err)

			req2 := httptest.NewRequest(http.MethodPost, "/token-email", bytes.NewReader(jsonData))
			req2.Header.Set("Content-Type", "application/json")
			req2.Header.Set("Authorization", token)

			rec2 := httptest.NewRecorder()
			c2 := e.NewContext(req2, rec2)

			errGetRequestTokenByEmail := adminController.GetRequestTokenByEmail(c2)

			assert.NoError(t, errGetRequestTokenByEmail)
			assert.Equal(t, http.StatusOK, rec2.Code)
			assert.Contains(t, rec2.Body.String(), "Success")
		})
	})

	t.Run("negative case to test admin controller, expected no error", func(t *testing.T) {
		t.Run("negative case account under recovery, expected no error", func(t *testing.T) {
			ctx := context.Background()
			result, err := testSuite.DBService.User.FindUserByEmail(ctx, reqBody.Email)

			require.NoError(t, err)

			target := &entity.User{
				ID: result.ID,
			}

			updateRequest := map[string]interface{}{
				"is_recovery": false,
			}

			err = testSuite.DBService.User.UpdateUser(ctx, target, updateRequest)

			require.NoError(t, err)

			jwtGenerator := security.NewJWT(testSuite.Config)
			token, err := jwtGenerator.GenerateJWT(result.ID, result.Email)
			require.NoError(t, err)

			err = testSuite.RedisClient.Set(ctx, fmt.Sprintf("user:%s:token", result.ID), token, 1*time.Minute).Err()
			require.NoError(t, err)

			req3 := httptest.NewRequest(http.MethodPost, "/token-email", bytes.NewReader(jsonData))
			req3.Header.Set("Authorization", token)
			req3.Header.Set("Content-Type", "application/json")

			rec3 := httptest.NewRecorder()
			c3 := e.NewContext(req3, rec3)

			adminController1 := admin.NewController(testSuite.Config, testSuite.DBService, testSuite.RedisClient, e.Logger)
			errGetRequestTokenByEmail1 := adminController1.GetRequestTokenByEmail(c3)

			assert.NoError(t, errGetRequestTokenByEmail1)
			assert.Equal(t, http.StatusBadRequest, rec3.Code)
			assert.Contains(t, rec3.Body.String(), "under recovery")
		})
	})
}
