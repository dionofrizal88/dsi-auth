package dao

import (
	"github.com/dionofrizal88/dsi/auth/domain/repository"
	"github.com/dionofrizal88/dsi/auth/infrastructure/persistence"
	"gorm.io/gorm"
)

// Repositories represent it self.
type Repositories struct {
	User repository.UserRepositoryInterface
	DB   *gorm.DB
}

// NewDBService will initialize db connection and return repositories.
func NewDBService(db *gorm.DB) *Repositories {
	return &Repositories{
		User: persistence.NewUserRepository(db),
		DB:   db,
	}
}
