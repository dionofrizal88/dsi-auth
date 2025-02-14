package registry_test

import (
	"github.com/dionofrizal88/dsi/auth/domain/registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistryDomainRegistry(t *testing.T) {

	t.Run("positive case to test domain registry, expected no error", func(t *testing.T) {
		t.Run("positive case while use func collect entities, expected no error", func(t *testing.T) {
			assert.NotEmpty(t, len(registry.CollectEntities()))
		})

		t.Run("positive case while use func collect table names, expected no error", func(t *testing.T) {
			assert.NotEmpty(t, len(registry.CollectTableNames()))
		})
	})

}
