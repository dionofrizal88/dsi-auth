package dao

import (
	domainRegistry "github.com/dionofrizal88/dsi/auth/domain/registry"
	"github.com/dionofrizal88/dsi/auth/pkg/registry"
)

// NewRegistry will initialize registry.Registry.
// Return []registry.Entity and []registry.Table.
// The registry is uses to auto migrate and reset database when running test mode.
func NewRegistry() *registry.Registry {
	var entityRegistry []registry.Entity
	var tableRegistry []registry.Table
	entityRegistry = append(entityRegistry, domainRegistry.CollectEntities()...)
	tableRegistry = append(tableRegistry, domainRegistry.CollectTableNames()...)

	return &registry.Registry{
		Entities: entityRegistry,
		Table:    tableRegistry,
	}
}
