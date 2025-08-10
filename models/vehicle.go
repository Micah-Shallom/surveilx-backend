package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID          string         `json:"id" gorm:"type:uuid;primary_key;"`
	UserID      string         `json:"user_id" gorm:"type:uuid;"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
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

type VisitorType string

const (
	VisitorTypeRegistered VisitorType = "registered"
	VisitorTypeGuest      VisitorType = "guest"
)

type VehicleActivity struct {
	ID          string      `json:"id" gorm:"type:uuid;primary_key;"`
	PlateNumber string      `json:"plate_number" gorm:"not null;index"`
	VisitorType VisitorType `json:"visitor_type" gorm:"type:varchar(20);not null;index"`

	// For registered vehicles
	VehicleID *string  `json:"vehicle_id,omitempty" gorm:"type:uuid;index"`
	UserID    *string  `json:"user_id,omitempty" gorm:"type:uuid;index"`
	Vehicle   *Vehicle `json:"vehicle,omitempty" gorm:"foreignKey:VehicleID"`
	User      *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Common fields
	IsEntry     bool   `json:"is_entry"`
	VehicleType string `json:"vehicle_type" gorm:"type:varchar(20)" validate:"oneof=bus car bike"`

	// Entry/Exit points
	EntryPointID *string          `json:"entry_point_id,omitempty" gorm:"type:uuid"`
	ExitPointID  *string          `json:"exit_point_id,omitempty" gorm:"type:uuid"`
	EntryPoint   *AccessExitPoint `json:"entry_point,omitempty" gorm:"foreignKey:EntryPointID"`
	ExitPoint    *AccessExitPoint `json:"exit_point,omitempty" gorm:"foreignKey:ExitPointID"`

	// Audit fields
	RegisteredBy     *string `json:"registered_by,omitempty" gorm:"type:uuid"` // Security personnel who logged guest
	RegisteredByUser *User   `json:"registered_by_user,omitempty" gorm:"foreignKey:RegisteredBy"`

	// Timestamps
	Timestamp time.Time      `json:"timestamp" gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (va *VehicleActivity) BeforeCreate(tx *gorm.DB) (err error) {
	va.ID = utility.GenerateUUID()
	if va.Timestamp.IsZero() {
		va.Timestamp = time.Now()
	}
	return
}

// Input models
type LogVehicleActivityInput struct {
	PlateNumber  string      `json:"plate_number" binding:"required"`
	VisitorType  VisitorType `json:"visitor_type" binding:"required" validate:"oneof=registered guest"`
	IsEntry      bool        `json:"is_entry"`
	EntryPointID string      `json:"entry_point_id,omitempty"`
	ExitPointID  string      `json:"exit_point_id,omitempty"`
}

// Response models for different contexts
type VehicleActivityResponse struct {
	ID          string      `json:"id"`
	PlateNumber string      `json:"plate_number"`
	VisitorType VisitorType `json:"visitor_type"`
	IsEntry     bool        `json:"is_entry"`
	VehicleType string      `json:"vehicle_type"`
	Timestamp   time.Time   `json:"timestamp"`
}
