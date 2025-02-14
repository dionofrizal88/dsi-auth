package dao_test

import (
	"github.com/dionofrizal88/dsi/auth/infrastructure/dao"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDaoRegistry(t *testing.T) {
	t.Run("positive case to test dao db repository, expected no error", func(t *testing.T) {
		t.Run("positive case while use func initiate repository, expected no error", func(t *testing.T) {
			registry := dao.NewRegistry()

			require.NotNil(t, registry)
		})
	})

}
