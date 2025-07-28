package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID        string `gorm:"type:uuid;primary_key;"`
	UserID    string `gorm:"type:uuid;"`
	User      User
	PlateNumber string `json:"plate_number" gorm:"unique"`
	Model     string `json:"model"`
	Color     string `json:"color"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (vehicle *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	vehicle.ID = utility.GenerateUUID()
	return
}

type VehicleLog struct {
	ID        string `gorm:"type:uuid;primary_key;"`
	VehicleID string `gorm:"type:uuid;"`
	Vehicle   Vehicle
	UserID    string `gorm:"type:uuid;"`
	User      User
	Timestamp time.Time `json:"timestamp"`
	IsEntry   bool      `json:"is_entry"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (vehicleLog *VehicleLog) BeforeCreate(tx *gorm.DB) (err error) {
	vehicleLog.ID = utility.GenerateUUID()
	return
}
