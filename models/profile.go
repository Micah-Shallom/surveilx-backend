package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID            string         `gorm:"type:uuid;primary_key;" json:"id"`
	UserID        string         `gorm:"type:uuid;not null" json:"user_id"` // Link to users table
	FirstName     string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName      string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Email         string         `gorm:"unique;not null" json:"email"`
	Phone         string         `gorm:"type:varchar(20)" json:"phone,omitempty"` // Useful for alerts
	AvatarURL     string         `gorm:"type:varchar(255)" json:"avatar_url,omitempty"`
	DateOfBirth   *time.Time     `gorm:"type:date" json:"date_of_birth,omitempty"`      
	Department    string         `gorm:"type:varchar(100)" json:"department,omitempty"` // Optional, depending on school
	IsBlacklisted bool           `gorm:"default:false" json:"is_blacklisted"`           // For flagged users
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (profile *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	profile.ID = utility.GenerateUUID()
	return
}

type UpdateProfileInput struct {
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Phone       string     `json:"phone,omitempty"`
	AvatarURL   string     `json:"avatar_url,omitempty"`
	Department  string     `json:"department,omitempty"`
}
