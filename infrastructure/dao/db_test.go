package dao_test

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/infrastructure/dao"
	"github.com/dionofrizal88/dsi/auth/pkg/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDaoDB(t *testing.T) {
	conf := config.GetConfig("../../config/config.json")

	dbConnection, _ := db.NewPostgresDBTestConnection(conf)

	t.Run("positive case to test dao db repository, expected no error", func(t *testing.T) {
		t.Run("positive case while use func initiate repository, expected no error", func(t *testing.T) {

			repo := dao.NewDBService(dbConnection)
			require.NotNil(t, repo)
		})
	})

}
