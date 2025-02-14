package persistence

import (
	"context"
	"github.com/dionofrizal88/dsi/auth/domain/entity"
	"github.com/dionofrizal88/dsi/auth/domain/repository"
	"github.com/dionofrizal88/dsi/auth/pkg/security"
	"gorm.io/gorm"
)

// UserRepo is a struct to store db connection.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepository will initialize UserRepo repository.
func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

// UserRepo implements the repository.UserRepositoryInterface interface.
var _ repository.UserRepositoryInterface = &UserRepo{}

// FindUser will find user by User id from database storage.
func (m UserRepo) FindUser(ctx context.Context, userID string) (*entity.User, error) {
	var dataEntity entity.User

	err := m.db.WithContext(ctx).Where("id = ?", userID).Take(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}

// CreateUser will create user into database storage.
func (m UserRepo) CreateUser(ctx context.Context, user *entity.User, salt string) (*entity.User, error) {
	// Hash the password with the salt
	hashedPassword := security.HashPasswordWithSalt(user.Password, salt)

	dataEntity := entity.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		IsAdmin:  user.IsAdmin,
	}

	err := m.db.WithContext(ctx).Create(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}

// UpdateUser is to update single row of data.
func (m UserRepo) UpdateUser(ctx context.Context, target *entity.User, value map[string]interface{}) error {
	err := m.db.WithContext(ctx).Model(target).Updates(value).Find(&target).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser is to delete single row of data.
func (m UserRepo) DeleteUser(ctx context.Context, target *entity.User) error {
	return m.db.WithContext(ctx).Delete(target).Error
}

// FindUserByEmail will find user by User email from database storage.
func (m UserRepo) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var dataEntity entity.User

	err := m.db.WithContext(ctx).Where("email = ?", email).Take(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}
