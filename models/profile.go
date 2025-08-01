package models

import (
	"time"

	"gorm.io/gorm"
)

// UpdateUserProfileInput represents the input for updating user profile
type UpdateUserProfileInput struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
	UserName    string `json:"username" validate:"required"`
	Phone       string `json:"phone" validate:"omitempty,phone"`
	AvatarURL   string `json:"avatar_url" validate:"omitempty,url"`
	DisplayName string `json:"display_name" validate:"omitempty"`
}

type Profile struct {
	ID          string         `gorm:"type:uuid;primary_key" json:"profile_id"`
	FirstName   string         `gorm:"column:first_name; type:text; not null" json:"first_name"`
	LastName    string         `gorm:"column:last_name; type:text;not null" json:"last_name"`
	FullName    string         `gorm:"column:full_name; type:text;" json:"full_name"`
	UserName    string         `gorm:"column:user_name; type:text;" json:"username"`
	Phone       string         `gorm:"type:varchar(255)" json:"phone"`
	AvatarURL   string         `gorm:"type:varchar(255)" json:"avatar_url"`
	UserID      string         `gorm:"type:uuid;" json:"user_id"`
	CreatedAt   time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	DisplayName string         `gorm:"type:varchar(255)" json:"display_name"`
}
