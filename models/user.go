package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID        string         `gorm:"column:id;type:uuid;primaryKey;"`
    Name      string         `json:"name" gorm:"column:name"`
    Email     string         `json:"email" gorm:"column:email;unique"`
    Password  string         `json:"-" gorm:"column:password"`
    Role      string         `json:"role" gorm:"column:role;default:'user'"`
    Token     string         `json:"token,omitempty" gorm:"column:token"`
    CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" gorm:"column:deleted_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = utility.GenerateUUID()
	return
}

func (user *User) CreateUser(db *gorm.DB) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}