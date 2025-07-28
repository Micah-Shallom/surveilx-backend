package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primary_key;"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"`
	Role      string         `json:"role" gorm:"default:'user'"`
	Token     string         `json:"token,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = utility.GenerateUUID()
	return
}
