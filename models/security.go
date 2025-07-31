package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Security struct {
	ID        string         `gorm:"type:uuid;primary_key;"`
	UserID    string         `gorm:"type:uuid;unique"`
	User      User           `gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (security *Security) BeforeCreate(tx *gorm.DB) (err error) {
	security.ID = utility.GenerateUUID()
	return
}
