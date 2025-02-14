package security_test

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecurityJWT(t *testing.T) {
	conf := config.GetConfig("../../config/config.json")
	jwt := security.NewJWT(conf)

	userJwt := &entity.User{
		ID:    "b79f40aa-2d78-4767-be8d-affb15537d82",
		Email: "test@gmail.com",
	}

	genJwt, errGenJwt := jwt.GenerateJWT(userJwt.ID, userJwt.Email)

	t.Run("positive case to test jwt, expected no error", func(t *testing.T) {
		t.Run("positive case while use func generate jwt, expected no error", func(t *testing.T) {

			assert.NoError(t, errGenJwt)
			assert.NotEmpty(t, len(genJwt))
		})

		t.Run("positive case while use func validate jwt, expected no error", func(t *testing.T) {
			isJwtValid := jwt.ValidateJWT(genJwt)

			assert.Equal(t, true, isJwtValid)
		})

		t.Run("positive case while use func get jwt, expected no error", func(t *testing.T) {
			getJwt := jwt.GetJWTClaims()

			assert.Equal(t, userJwt.ID, getJwt.UserID)
			assert.Equal(t, userJwt.Email, getJwt.Email)
		})

		t.Run("positive case while use func validate jwt, expected no error", func(t *testing.T) {
			decodeJwt, errDecodeJwt := jwt.DecodeJWT(genJwt)

			assert.NoError(t, errDecodeJwt)
			assert.Equal(t, userJwt.ID, decodeJwt["user_id"])
			assert.Equal(t, userJwt.Email, decodeJwt["email"])
		})
	})

}
