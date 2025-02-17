package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dionofrizal88/dsi/auth/controllers/user"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/dionofrizal88/dsi/auth/pkg/tests"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestControllerUser(t *testing.T) {
	testSuite := tests.InitTestSuite()

	t.Cleanup(func() {
		testSuite.Clean()
	})

	reqBody := user.Request{
		Name:     "Test User 1",
		Email:    "testUser1@gmail.com",
		Password: "testUser1234",
	}

	reqBodyRecovery := user.RecoveryRequest{
		Email:       "testUser1@gmail.com",
		NewPassword: "testUser12344",
	}

	jsonData, _ := json.Marshal(reqBody)
	jsonDataRecovery, _ := json.Marshal(reqBodyRecovery)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userController := user.NewController(testSuite.Config, testSuite.DBService, testSuite.RedisClient, e.Logger)
	errRegisterUser := userController.RegisterUser(c)

	t.Run("positive case to test user controller, expected no error", func(t *testing.T) {
		t.Run("positive case while send register user request, expected no error", func(t *testing.T) {

			assert.NoError(t, errRegisterUser)
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Contains(t, rec.Body.String(), "Success")
		})
	})

	t.Run("negative case to test user controller, expected no error", func(t *testing.T) {
		t.Run("negative case while send register user already exists, expected no error", func(t *testing.T) {
			req1 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonData))
			req1.Header.Set("Content-Type", "application/json")

			rec1 := httptest.NewRecorder()
			c1 := e.NewContext(req1, rec1)

			errRegisterUser1 := userController.RegisterUser(c1)

			assert.NoError(t, errRegisterUser1)
			assert.Equal(t, http.StatusConflict, rec1.Code)
			assert.Contains(t, rec1.Body.String(), "exists")
		})
	})

	t.Run("positive case to test user controller, expected no error", func(t *testing.T) {
		t.Run("positive case while login user, expected no error", func(t *testing.T) {
			req1 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonData))
			req1.Header.Set("Content-Type", "application/json")
			rec1 := httptest.NewRecorder()
			c1 := e.NewContext(req1, rec1)
			errLoginUser1 := userController.Login(c1)

			assert.NoError(t, errLoginUser1)
			assert.Equal(t, http.StatusOK, rec1.Code)
			assert.Contains(t, rec1.Body.String(), "Success")
		})
	})

	t.Run("negative case to test user controller, expected no error", func(t *testing.T) {
		t.Run("negative case while login user incorrect password or email, expected no error", func(t *testing.T) {
			reqBodyWrongEmailPassword := user.Request{
				Name:     "Test User 1",
				Email:    "testUser1@gmail.com",
				Password: "wrongPassword",
			}
			jsonReqBodyWrongEmailPassword, _ := json.Marshal(reqBodyWrongEmailPassword)

			req1 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonReqBodyWrongEmailPassword))
			req1.Header.Set("Content-Type", "application/json")
			rec1 := httptest.NewRecorder()
			c1 := e.NewContext(req1, rec1)
			errLoginUser1 := userController.Login(c1)

			assert.NoError(t, errLoginUser1)
			assert.Equal(t, http.StatusUnauthorized, rec1.Code)
			assert.Contains(t, rec1.Body.String(), "incorrect")
		})

		t.Run("negative case while login user account under recovery, expected no error", func(t *testing.T) {
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

			req1 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonData))
			req1.Header.Set("Content-Type", "application/json")
			rec1 := httptest.NewRecorder()
			c1 := e.NewContext(req1, rec1)
			errLoginUser1 := userController.Login(c1)

			assert.NoError(t, errLoginUser1)
			assert.Equal(t, http.StatusForbidden, rec1.Code)
			assert.Contains(t, rec1.Body.String(), "under recovery")
		})
	})

	t.Run("positive case to test user controller, expected no error", func(t *testing.T) {
		t.Run("positive case while logout user, expected no error", func(t *testing.T) {
			ctx := context.Background()
			result, err := testSuite.DBService.User.FindUserByEmail(ctx, reqBody.Email)

			require.NoError(t, err)

			jwtGenerator := security.NewJWT(testSuite.Config)
			token, err := jwtGenerator.GenerateJWT(result.ID, result.Email)
			require.NoError(t, err)

			err = testSuite.RedisClient.Set(ctx, fmt.Sprintf("user:%s:token", result.ID), token, 1*time.Minute).Err()
			require.NoError(t, err)

			req2 := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewReader(jsonData))
			req2.Header.Set("Content-Type", "application/json")
			req2.Header.Set("authorization", token)

			rec2 := httptest.NewRecorder()
			c2 := e.NewContext(req2, rec2)
			errLogoutUser1 := userController.Logout(c2)

			assert.NoError(t, errLogoutUser1)
			assert.Equal(t, http.StatusOK, rec2.Code)
			assert.Contains(t, rec2.Body.String(), "Success")
		})
	})

	t.Run("negative case to test user controller, expected no error", func(t *testing.T) {
		t.Run("negative case while logout user missing token, expected no error", func(t *testing.T) {
			req2 := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewReader(jsonData))
			req2.Header.Set("Content-Type", "application/json")

			rec2 := httptest.NewRecorder()
			c2 := e.NewContext(req2, rec2)
			errLogoutUser1 := userController.Logout(c2)

			assert.NoError(t, errLogoutUser1)
			assert.Equal(t, http.StatusBadRequest, rec2.Code)
			assert.Contains(t, rec2.Body.String(), "Missing")
		})
	})

	t.Run("positive case to test user controller, expected no error", func(t *testing.T) {
		t.Run("positive case while request recovery user, expected no error", func(t *testing.T) {

			req3 := httptest.NewRequest(http.MethodPost, "/request-recovery", bytes.NewReader(jsonData))
			req3.Header.Set("Content-Type", "application/json")

			rec3 := httptest.NewRecorder()
			c3 := e.NewContext(req3, rec3)
			errRequestRecoveryUser1 := userController.RequestRecovery(c3)

			assert.NoError(t, errRequestRecoveryUser1)
			assert.Equal(t, http.StatusAccepted, rec3.Code)
			assert.Contains(t, rec3.Body.String(), "Accepted")
		})
	})

	t.Run("negative case to test user controller, expected no error", func(t *testing.T) {
		t.Run("negative case while failed request recovery user wrong email, expected no error", func(t *testing.T) {
			reqBodyWrongEmail := user.RecoveryRequest{
				Email: "testWrongUser1@gmail.com",
			}
			jsonReqBodyWrongEmail, _ := json.Marshal(reqBodyWrongEmail)

			req3 := httptest.NewRequest(http.MethodPost, "/request-recovery", bytes.NewReader(jsonReqBodyWrongEmail))
			req3.Header.Set("Content-Type", "application/json")

			rec3 := httptest.NewRecorder()
			c3 := e.NewContext(req3, rec3)
			errRequestRecoveryUser1 := userController.RequestRecovery(c3)

			assert.NoError(t, errRequestRecoveryUser1)
			assert.Equal(t, http.StatusConflict, rec3.Code)
			assert.Contains(t, rec3.Body.String(), "Failed")
		})
	})

	t.Run("positive case to test user controller, expected no error", func(t *testing.T) {
		t.Run("positive case while recovery user, expected no error", func(t *testing.T) {
			ctx := context.Background()
			result, err := testSuite.DBService.User.FindUserByEmail(ctx, reqBody.Email)

			require.NoError(t, err)

			randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
			token := security.GenerateRecoveryToken(fmt.Sprintf("%s.%d", result.Email, randomizer.Int()), testSuite.Config.AppSecret+result.Email)

			err = testSuite.RedisClient.Set(ctx, fmt.Sprintf("recovery-user:%s:token", result.ID), token, 1*time.Minute).Err()
			require.NoError(t, err)

			req4 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/recovery/%s", token), bytes.NewReader(jsonDataRecovery))
			req4.Header.Set("Content-Type", "application/json")

			rec4 := httptest.NewRecorder()
			c4 := e.NewContext(req4, rec4)
			c4.SetParamNames("token")
			c4.SetParamValues(token)
			errRecoveryUser1 := userController.Recovery(c4)

			assert.NoError(t, errRecoveryUser1)
			assert.Equal(t, http.StatusOK, rec4.Code)
			assert.Contains(t, rec4.Body.String(), "Success")
		})
	})

	t.Run("negative case to test user controller, expected no error", func(t *testing.T) {
		t.Run("negative case while recovery user token is missing, expected no error", func(t *testing.T) {
			req4 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/recovery/%s", ""), bytes.NewReader(jsonDataRecovery))
			req4.Header.Set("Content-Type", "application/json")

			rec4 := httptest.NewRecorder()
			c4 := e.NewContext(req4, rec4)

			errRecoveryUser1 := userController.Recovery(c4)

			assert.NoError(t, errRecoveryUser1)
			assert.Equal(t, http.StatusUnauthorized, rec4.Code)
			assert.Contains(t, rec4.Body.String(), "Missing token")
		})
	})

}
