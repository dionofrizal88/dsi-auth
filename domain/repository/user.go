package repository

import (
	"context"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
)

// UserRepositoryInterface need to be implemented in persistence repository.
type UserRepositoryInterface interface {
	FindUser(ctx context.Context, userID string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User, salt string) (*entity.User, error)
	UpdateUser(ctx context.Context, target *entity.User, value map[string]interface{}) error
	DeleteUser(ctx context.Context, target *entity.User) error
}
