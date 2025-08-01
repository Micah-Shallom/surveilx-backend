package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type GuestWatchList struct {
	ID           string         `json:"id" gorm:"type:uuid;primary_key;"`
	PlateNumber  string         `json:"plate_number"`
	IsEntry      bool           `json:"is_entry"`
	Type         string         `json:"type" validate:"oneof=bus car bike"`
	RegisteredBy string         `json:"registered_by" gorm:"type:uuid;"`
	Timestamp    time.Time      `json:"timestamp"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (guestLog *GuestWatchList) BeforeCreate(tx *gorm.DB) (err error) {
	guestLog.ID = utility.GenerateUUID()
	return
}

type LogGuestInput struct {
	PlateNumber  string `json:"plate_number" binding:"required"`
	IsEntry      bool   `json:"is_entry"`
	Type         string `json:"type" validate:"required,oneof=bus car bike"`
	RegisteredBy string `json:"registered_by" binding:"required"`
}
