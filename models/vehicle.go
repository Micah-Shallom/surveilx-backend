package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID          string         `json:"id" gorm:"column:id;type:uuid;primaryKey;"`
	UserID      string         `json:"user_id" gorm:"column:user_id;type:uuid;"`
	PlateNumber string         `json:"plate_number" gorm:"column:plate_number;unique"`
	Type        string         `json:"type" validate:"oneof=bus car bike" gorm:"column:type"`
	Model       string         `json:"model" gorm:"column:model"`
	Color       string         `json:"color" gorm:"column:color"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"column:created_at"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at"`
}

type VehicleInfo struct {
	PlateNumber string `json:"plate_number,omitempty" gorm:"column:plate_number;uniqueIndex"`
	Type        string `json:"type,omitempty" validate:"oneof=bus car bike" gorm:"column:type"`
	Model       string `json:"model,omitempty" gorm:"column:model"`
	Color       string `json:"color,omitempty" gorm:"column:color"`
}

func (vehicle *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	vehicle.ID = utility.GenerateUUID()
	return
}

type VisitorType string

const (
	VisitorTypeRegistered VisitorType = "registered"
	VisitorTypeGuest      VisitorType = "guest"
)

type VehicleActivity struct {
	ID          string      `json:"id" gorm:"column:id;type:uuid;primary_key;"`
	PlateNumber string      `json:"plate_number" gorm:"column:plate_number;not null;index"`
	Model       string      `json:"model" gorm:"column:model;not null;index"`
	VisitorType VisitorType `json:"visitor_type" gorm:"column:visitor_type;type:varchar(20);not null;index"`

	Vehicle   *Vehicle `json:"vehicle,omitempty" gorm:"foreignKey:VehicleID"`
	VehicleID *string  `json:"vehicle_id,omitempty" gorm:"column:vehicle_id;type:uuid;index"`

	// Common fields
	IsEntry     bool   `json:"is_entry" gorm:"column:is_entry"`
	VehicleType string `json:"vehicle_type" gorm:"column:vehicle_type;type:varchar(20)" validate:"oneof=bus car bike"`

	// Entry/Exit points
	EntryPointID *string          `json:"entry_point_id,omitempty" gorm:"column:entry_point_id;type:uuid"`
	ExitPointID  *string          `json:"exit_point_id,omitempty" gorm:"column:exit_point_id;type:uuid"`
	EntryPoint   *AccessExitPoint `json:"entry_point,omitempty" gorm:"foreignKey:EntryPointID"`
	ExitPoint    *AccessExitPoint `json:"exit_point,omitempty" gorm:"foreignKey:ExitPointID"`
	GateName     string           `json:"gate_name,omitempty" gorm:"column:gate_name"`

	Timestamp time.Time      `json:"timestamp" gorm:"column:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}

type GuestVehicleActivity struct {
	ID          string `json:"id" gorm:"column:id;type:uuid;primaryKey;"`
	PlateNumber string `json:"plate_number" gorm:"column:plate_number;"`

	IsEntry      bool             `json:"is_entry" gorm:"column:is_entry"`
	EntryPointID *string          `json:"entry_point_id,omitempty" gorm:"column:entry_point_id;type:uuid"`
	ExitPointID  *string          `json:"exit_point_id,omitempty" gorm:"column:exit_point_id;type:uuid"`
	EntryPoint   *AccessExitPoint `json:"entry_point,omitempty" gorm:"foreignKey:EntryPointID"`
	ExitPoint    *AccessExitPoint `json:"exit_point,omitempty" gorm:"foreignKey:ExitPointID"`

	Timestamp time.Time      `json:"timestamp" gorm:"column:timestamp;not null;default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at"`
}

func (va *VehicleActivity) BeforeCreate(tx *gorm.DB) (err error) {
	va.ID = utility.GenerateUUID()
	if va.Timestamp.IsZero() {
		va.Timestamp = time.Now()
	}
	return
}

func (va *GuestVehicleActivity) BeforeCreate(tx *gorm.DB) (err error) {
	va.ID = utility.GenerateUUID()
	if va.Timestamp.IsZero() {
		va.Timestamp = time.Now()
	}
	return
}

type LogVehicleActivityInput struct {
	PlateNumber  string      `json:"plate_number" binding:"required"`
	VisitorType  VisitorType `json:"visitor_type"`
	IsEntry      bool        `json:"is_entry"`
	EntryPointID string      `json:"entry_point_id,omitempty"`
	ExitPointID  string      `json:"exit_point_id,omitempty"`
}

type VehicleActivityResponse struct {
	ID          string      `json:"id"`
	PlateNumber string      `json:"plate_number"`
	VisitorType VisitorType `json:"visitor_type"`
	IsEntry     bool        `json:"is_entry"`
	VehicleType string      `json:"vehicle_type,omitempty"`
	Model       string      `json:"model,omitempty"`
	GateName    string      `json:"gate_name,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
}

type VehicleIdentity struct {
	PlateNumber  string `json:"plate_number"`
	Status       string `json:"status"`
	IsRegistered bool   `json:"is_registered"`
}

func (v *Vehicle) DeRegister(db *gorm.DB) error {
	if err := db.Unscoped().Delete(&v).Error; err != nil {
		return err
	}

	return nil
}
