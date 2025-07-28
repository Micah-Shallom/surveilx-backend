package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        string         `gorm:"type:uuid;primary_key;"`
	UserID    string         `gorm:"type:uuid;unique"`
	User      User           `gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (admin *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	admin.ID = utility.GenerateUUID()
	return
}
