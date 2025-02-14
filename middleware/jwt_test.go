package middleware_test

import (
	"context"
	"fmt"
	"github.com/dionofrizal88/dsi/auth/middleware"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/dionofrizal88/dsi/auth/pkg/tests"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddlewareJWT(t *testing.T) {
	testSuite := tests.InitTestSuite()

	t.Cleanup(func() {
		testSuite.Clean()
	})

	ctx := context.Background()
	jwtMiddleware := middleware.JWTMiddleware(testSuite.Config, testSuite.RedisClient)

	t.Run("positive case to test middleware, expected no error", func(t *testing.T) {
		t.Run("positive case while send valid authorization token, expected no error", func(t *testing.T) {
			id := uuid.NewString() + "-test"
			jwtGenerator := security.NewJWT(testSuite.Config)
			token, err := jwtGenerator.GenerateJWT(id, "example1234@gmail.com")

			assert.NoError(t, err)

			err = testSuite.RedisClient.Set(ctx, fmt.Sprintf("user:%s:token", id), token, 3*time.Minute).Err()

			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = jwtMiddleware(func(c echo.Context) error {
				return c.JSON(http.StatusOK, "Success")
			})(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Success")
		})
	})

	t.Run("negative case to test middleware, expected no error", func(t *testing.T) {
		t.Run("negative case while not send authorization token, expected no error", func(t *testing.T) {

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := jwtMiddleware(func(c echo.Context) error {
				return c.JSON(http.StatusOK, "Success")
			})(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Contains(t, rec.Body.String(), "Missing token")
		})

		t.Run("negative case while send invalid authorization token, expected no error", func(t *testing.T) {

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "invalid_token")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := jwtMiddleware(func(c echo.Context) error {
				return c.JSON(http.StatusOK, "Success")
			})(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Contains(t, rec.Body.String(), "Invalid token")
		})
	})

}
