package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represent schema of table users.
type User struct {
	ID         string         `gorm:"size:36;not null;unique;index;primary_key;" json:"id" uri:"id"`
	Name       string         `gorm:"size:100;not null;index;" json:"name" form:"name"`
	Email      string         `gorm:"size:100;" json:"email" form:"email"`
	Password   string         `gorm:"size:255;" json:"password" form:"password"`
	IsRecovery bool           `json:"is_recovery" form:"is_recovery"`
	IsAdmin    bool           `json:"is_admin" form:"is_admin"`
	CreatedAt  *time.Time     `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
}

var _ Interface = &User{}

// Users represent multiple Users.
type Users []*User

// TableName return name of table.
func (u *User) TableName() string {
	return "users"
}

// FilterableFields return fields.
func (u *User) FilterableFields() []interface{} {
	return []interface{}{"id"}
}

// TimeFields return fields.
func (u *User) TimeFields() []interface{} {
	return []interface{}{"created_at", "updated_at", "deleted_at"}
}

// BeforeCreate handle hook before create row in db.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		generateUUID := uuid.New()
		u.ID = generateUUID.String()
	}

	return nil
}
