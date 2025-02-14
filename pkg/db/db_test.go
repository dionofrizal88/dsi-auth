package db_test

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDBSQL(t *testing.T) {
	conf := config.GetConfig("../../config/config.json")

	t.Run("positive case to test db sql, expected no error", func(t *testing.T) {
		t.Run("positive case while use func to create postgres connection, expected no error", func(t *testing.T) {
			dbConnection, errDBConnection := db.NewPostgresDBConnection(conf)

			require.NoError(t, errDBConnection)
			assert.NotNil(t, dbConnection)
		})

		t.Run("positive case while use func to create test postgres connection, expected no error", func(t *testing.T) {
			dbConnection, errDBConnection := db.NewPostgresDBTestConnection(conf)

			require.NoError(t, errDBConnection)
			assert.NotNil(t, dbConnection)
		})
	})

}
