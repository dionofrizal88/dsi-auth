package entity_test

import (
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEntityUser(t *testing.T) {
	userData1 := &entity.User{
		Name:     "User Testing 1",
		Email:    "userTesting1@gmail.com",
		Password: "userTesting1234",
	}

	t.Run("positive case to test user entity, expected no error", func(t *testing.T) {
		t.Run("positive case while func on user entity, expected no error", func(t *testing.T) {

			assert.NotEmpty(t, userData1.TableName())
			assert.NotEmpty(t, userData1.FilterableFields())
			assert.NotEmpty(t, userData1.TimeFields())
		})

		t.Run("positive case while func before create on user entity, expected no error", func(t *testing.T) {
			err := userData1.BeforeCreate(nil)

			require.NoError(t, err)
			assert.NotEmpty(t, userData1.ID)
		})

	})

}
