package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

// UpdateUserProfileInput represents the input for updating user profile
type UpdateUserProfileInput struct {
	Phone    string `json:"phone"`
	UserName string `json:"username"`
	FullName string `json:"full_name"`
}

type Profile struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"profile_id"`
	FullName  string         `gorm:"column:full_name; type:text;" json:"full_name"`
	UserName  string         `gorm:"column:user_name; type:text;" json:"username"`
	Phone     string         `gorm:"type:varchar(255)" json:"phone"`
	UserID    string         `gorm:"type:uuid;" json:"user_id"`
	CreatedAt time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Profile) GetUserProfile(db *gorm.DB, user_id string) error {
	return db.Where("user_id = ?", user_id).First(&p).Error
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = utility.GenerateUUID()
	return
}

func (p *Profile) CreateProfile(db *gorm.DB) error {
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}
