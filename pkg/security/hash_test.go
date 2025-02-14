package security_test

import (
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecurityHash(t *testing.T) {
	salt := security.GenerateSalt(64)

	t.Run("positive case to test hash security, expected no error", func(t *testing.T) {
		t.Run("positive case while use func generate salt, expected no error", func(t *testing.T) {

			assert.NotEmpty(t, len(salt))
			assert.Equal(t, 64, len(salt))
		})

		t.Run("positive case while use func hash password with salt, expected no error", func(t *testing.T) {
			hashPassword := security.HashPasswordWithSalt("examplePassword12344", salt)

			assert.NotEmpty(t, hashPassword)
		})

		t.Run("positive case while use func hash password with salt, expected no error", func(t *testing.T) {
			hashPassword := security.GenerateRecoveryToken("exampleString12344", salt)

			assert.NotEmpty(t, hashPassword)
		})

	})

}
