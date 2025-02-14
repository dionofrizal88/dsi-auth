package routes_test

import (
	"github.com/dionofrizal88/dsi/auth/pkg/tests"
	"github.com/dionofrizal88/dsi/auth/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoutesInitiate(t *testing.T) {
	testSuite := tests.InitTestSuite()

	t.Cleanup(func() {
		testSuite.Clean()
	})

	// Init Router
	e := routes.
		NewRouter(
			routes.WithConfig(testSuite.Config),
			routes.WithDBService(testSuite.DBService),
			routes.WithRedisDB(testSuite.RedisClient),
		).
		Init()

	t.Run("positive case to test routes, expected no error", func(t *testing.T) {
		t.Run("positive case while use func routes initiate, expected no error", func(t *testing.T) {

			require.NotNil(t, e)
			assert.NotEmpty(t, e.Logger.Errorf)
		})
	})

}
