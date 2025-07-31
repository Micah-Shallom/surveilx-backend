package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type GuestLog struct {
	ID          string         `json:"id" gorm:"type:uuid;primary_key;"`
	PlateNumber string         `json:"plate_number"`
	IsEntry     bool           `json:"is_entry"`
	Timestamp   time.Time      `json:"timestamp"`
	CreatedAt   time.Time      `json:"createdAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (guestLog *GuestLog) BeforeCreate(tx *gorm.DB) (err error) {
	guestLog.ID = utility.GenerateUUID()
	return
}

type LogGuestInput struct {
	PlateNumber string `json:"plate_number" binding:"required"`
	IsEntry     bool   `json:"is_entry"`
}
