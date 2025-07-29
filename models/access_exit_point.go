package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type AccessExitPoint struct {
	ID        string         `gorm:"type:uuid;primary_key;"`
	Name      string         `json:"name" gorm:"unique"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (point *AccessExitPoint) BeforeCreate(tx *gorm.DB) (err error) {
	point.ID = utility.GenerateUUID()
	return
}
