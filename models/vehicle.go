package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID          string         `json:"id" gorm:"type:uuid;primary_key;"`
	UserID      string         `json:"user_id" gorm:"type:uuid;"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	PlateNumber string         `json:"plate_number" gorm:"unique"`
	Type        string         `json:"type" validate:"oneof=bus car bike"`
	Model       string         `json:"model"`
	Color       string         `json:"color"`
	CreatedAt   time.Time      `json:"createdAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (vehicle *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	vehicle.ID = utility.GenerateUUID()
	return
}

type VehicleLog struct {
	ID            string          `json:"id" gorm:"type:uuid;primary_key;"`
	VehicleID     string          `json:"uuid" gorm:"type:uuid;"`
	Vehicle       Vehicle         `json:"vehicle" gorm:"foreignKey:VehicleID"`
	UserID        string          `json:"user_id" gorm:"type:uuid;"`
	User          User            `json:"user" gorm:"foreignKey:UserID"`
	Timestamp     time.Time       `json:"timestamp"`
	IsEntry       bool            `json:"is_entry"`
	EntryPointID  *string         `json:"entry_point_id" gorm:"type:uuid"`
	EntryPoint    AccessExitPoint `json:"entry_point" gorm:"foreignKey:EntryPointID"`
	ExitPointID   *string         `json:"exit_point_id" gorm:"type:uuid"`
	ExitPoint     AccessExitPoint `json:"exit_point" gorm:"foreignKey:ExitPointID"`
	CreatedAt     time.Time       `json:"createdAt"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"deletedAt"`
}

func (vehicleLog *VehicleLog) BeforeCreate(tx *gorm.DB) (err error) {
	vehicleLog.ID = utility.GenerateUUID()
	return
}
