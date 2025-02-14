package registry

import (
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/pkg/registry"
)

// CollectEntities will return collections of replication entity.
func CollectEntities() []registry.Entity {
	return []registry.Entity{
		{Entity: entity.User{}},
	}
}

// CollectTableNames will return collections of replication table name.
func CollectTableNames() []registry.Table {
	var (
		user entity.User
	)

	return []registry.Table{
		{Name: user.TableName()},
	}
}
