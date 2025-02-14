package persistence_test

import (
	"context"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/infrastructure/persistence"
	"github.com/dionofrizal88/dsi/auth/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPersistenceUser(t *testing.T) {
	testSuite := tests.InitTestSuite()

	t.Cleanup(func() {
		testSuite.Clean()
	})

	ctx := context.Background()
	userData1 := &entity.User{
		ID:       "780d709e-efb9-4916-930d-0e81f0cf1e0e",
		Name:     "User Testing 1",
		Email:    "userTesting1@gmail.com",
		Password: "userTesting1234",
	}

	t.Run("positive case to test user repository, expected no error", func(t *testing.T) {
		repo := persistence.NewUserRepository(testSuite.DBService.DB)

		t.Run("positive case while create user, expected no error", func(t *testing.T) {

			data, err := repo.CreateUser(ctx, userData1, testSuite.Config.AppSecret+userData1.Email)
			require.NoError(t, err)
			assert.NotNil(t, data)
			assert.Equal(t, userData1.Name, data.Name)
			assert.Equal(t, userData1.Email, data.Email)
		})

		t.Run("positive case while find user, expected no error", func(t *testing.T) {

			findUserData, err := repo.FindUser(ctx, userData1.ID)
			require.NoError(t, err)
			assert.NotNil(t, findUserData)
			assert.Equal(t, userData1.Name, findUserData.Name)
			assert.Equal(t, userData1.Email, findUserData.Email)
		})

		t.Run("positive case while find user by email, expected no error", func(t *testing.T) {

			findUserByEmailData, err := repo.FindUserByEmail(ctx, userData1.Email)
			require.NoError(t, err)
			assert.NotNil(t, findUserByEmailData)
			assert.Equal(t, userData1.Name, findUserByEmailData.Name)
			assert.Equal(t, userData1.Email, findUserByEmailData.Email)
		})

		t.Run("positive case while update user data, expected no error", func(t *testing.T) {

			target := &entity.User{
				ID: userData1.ID,
			}

			updateRequest := map[string]interface{}{
				"name":  "User Testing 1 updated",
				"email": "updateUserTesting1@gmail.com",
			}

			err := repo.UpdateUser(ctx, target, updateRequest)
			require.NoError(t, err)
			assert.Equal(t, updateRequest["name"], target.Name)
			assert.Equal(t, updateRequest["email"], target.Email)
		})

		t.Run("positive case while delete user, expected no error", func(t *testing.T) {

			target := &entity.User{
				ID: userData1.ID,
			}
			err := repo.DeleteUser(ctx, target)
			require.NoError(t, err)
		})
	})

}
